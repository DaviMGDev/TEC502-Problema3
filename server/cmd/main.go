package main

import (
	"cod-server/internal/api"
	"cod-server/internal/api/mqtt"
	"cod-server/internal/cluster"
	"cod-server/internal/data"
	"cod-server/internal/domain"
	"cod-server/internal/services"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
	"github.com/joho/godotenv"
)

// getEnv carrega uma variável de ambiente ou retorna um valor padrão.
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	// Carrega o arquivo .env (se existir)
	err := godotenv.Load()
	if err != nil {
		log.Printf("Aviso: Não foi possível carregar o arquivo .env: %v. Usando variáveis de ambiente existentes ou padrões.", err)
	}

	// Carrega as configurações das variáveis de ambiente
	raftDataDir := getEnv("COD_RAFT_DATA_DIR", "./raft-data")
	raftBindAddr := getEnv("COD_RAFT_BIND_ADDR", "127.0.0.1:10000")
	httpBindAddr := getEnv("COD_HTTP_BIND_ADDR", "127.0.0.1:8080")
	nodeID := getEnv("COD_NODE_ID", "node-1")
	mqttBrokerAddr := getEnv("COD_MQTT_BROKER_ADDR", "tcp://localhost:1883")
	isFirstNodeStr := getEnv("COD_IS_FIRST_NODE", "false")
	isFirstNode, _ := strconv.ParseBool(isFirstNodeStr)

	log.Println("Iniciando servidor COD...")

	// 1. Inicialização de Repositórios e Serviços
	userRepo := data.NewMemoryRepository[domain.UserInterface]()
	cardRepo := data.NewMemoryRepository[domain.CardInterface]()
	matchRepo := data.NewMemoryRepository[domain.MatchInterface]()

	userService := services.NewUserService(userRepo)
	cardsService := services.NewCardsService(cardRepo)
	matchService := services.NewMatchService(matchRepo, cardRepo, userRepo)

	// 2. Inicialização dos Handlers da API
	eventHandler := api.NewEventHandler(userService, cardsService, matchService)

	// 3. Inicialização da FSM (Máquina de Estados do Raft)
	fsm := cluster.NewClusterFSM(eventHandler)

	// 4. Configuração do Raft
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(nodeID)

	addr, err := net.ResolveTCPAddr("tcp", raftBindAddr)
	if err != nil {
		log.Fatalf("falha ao resolver endereço TCP do Raft: %v", err)
	}
	transport, err := raft.NewTCPTransport(raftBindAddr, addr, 3, 10*time.Second, os.Stderr)
	if err != nil {
		log.Fatalf("falha ao criar transporte TCP do Raft: %v", err)
	}

	if err := os.MkdirAll(raftDataDir, 0755); err != nil {
		log.Fatalf("falha ao criar diretório de dados do Raft: %v", err)
	}

	logStore, err := raftboltdb.NewBoltStore(raftDataDir + "/logs.db")
	if err != nil {
		log.Fatalf("falha ao criar log store: %v", err)
	}
	stableStore, err := raftboltdb.NewBoltStore(raftDataDir + "/stable.db")
	if err != nil {
		log.Fatalf("falha ao criar stable store: %v", err)
	}
	snapshotStore, err := raft.NewFileSnapshotStore(raftDataDir, 1, os.Stderr)
	if err != nil {
		log.Fatalf("falha ao criar snapshot store: %v", err)
	}

	raftNode, err := raft.NewRaft(config, fsm, logStore, stableStore, snapshotStore, transport)
	if err != nil {
		log.Fatalf("falha ao criar nó Raft: %v", err)
	}

	// Bootstrap (apenas para o primeiro nó)
	if isFirstNode {
		log.Println("Realizando bootstrap do cluster...")
		configuration := raft.Configuration{
			Servers: []raft.Server{
				{
					ID:      config.LocalID,
					Address: transport.LocalAddr(),
				},
			},
		}
		if err := raftNode.BootstrapCluster(configuration).Error(); err != nil {
			log.Fatalf("falha ao realizar bootstrap do cluster: %v", err)
		}
	}

	// 5. Inicialização da API Interna (HTTP)
	httpTransport := cluster.NewGinHttpTransport(httpBindAddr, nodeID, raftNode)
	if err := httpTransport.Start(); err != nil {
		log.Fatalf("Falha ao iniciar transporte HTTP: %v", err)
	}

	// 6. Inicialização do Coordenador
	coordinator := cluster.NewRaftCoordinator(raftNode, httpTransport)

	// 7. Configuração do MQTT
	mqttAdapter, err := mqtt.NewMQTTAdapter(mqttBrokerAddr, nodeID)
	if err != nil {
		log.Fatalf("Falha ao criar adaptador MQTT: %v", err)
	}
	if err := mqttAdapter.Connect(); err != nil {
		log.Fatalf("Falha ao conectar ao broker MQTT: %v", err)
	}
	defer mqttAdapter.Disconnect()

	mqttAdapter.Subscribe("game/actions", func(client paho.Client, msg paho.Message) {
		event, err := api.FromJson(msg.Payload())
		if err != nil {
			log.Printf("Erro ao desserializar evento MQTT: %v", err)
			return
		}
		log.Printf("Evento MQTT recebido: %+v", event)
		if err := coordinator.Handle(*event); err != nil {
			log.Printf("Erro ao processar evento via coordenador: %v", err)
		}
	})

	// 8. Serviço de Descoberta
	discovery := cluster.NewDiscoveryService(string(transport.LocalAddr()))
	discovery.OnPeerDiscovered = func(peerIp string) {
		targetAddr := fmt.Sprintf("%s:%s", peerIp, "8080") // Assumindo porta 8080
		log.Printf("Nó par descoberto em %s. Tentando adicionar ao cluster...", targetAddr)

		// Lógica para evitar adicionar nós duplicados ou a si mesmo
		if raftNode.State() != raft.Leader {
			log.Println("Não sou o líder, não posso adicionar um novo nó.")
			return
		}

		future := raftNode.GetConfiguration()
		if err := future.Error(); err != nil {
			log.Printf("Falha ao obter configuração do cluster: %v", err)
			return
		}
		for _, srv := range future.Configuration().Servers {
			if srv.Address == raft.ServerAddress(targetAddr) {
				log.Printf("Nó %s já faz parte do cluster.", targetAddr)
				return
			}
		}

		// Adicionando o nó descoberto
		addVoterFuture := raftNode.AddVoter(raft.ServerID(peerIp), raft.ServerAddress(targetAddr), 0, 0)
		if err := addVoterFuture.Error(); err != nil {
			log.Printf("Falha ao adicionar novo nó ao cluster: %v", err)
		} else {
			log.Printf("Nó %s adicionado ao cluster.", targetAddr)
		}
	}
	discovery.Start()

	// 9. Bloqueio para manter o servidor rodando
	log.Println("Servidor COD rodando. Pressione CTRL-C para sair.")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	log.Println("Desligando o servidor...")
}