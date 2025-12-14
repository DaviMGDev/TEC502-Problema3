package cluster

import (
	"fmt"
	"log"
	"net"
	"time"
)

const (
	DiscoveryPort     = 9999
	DiscoveryInterval = 5 * time.Second
	DiscoveryMessage  = "COD_SERVER_DISCOVERY"
)

// DiscoveryServiceInterface define o contrato para o serviço de descoberta
type DiscoveryServiceInterface interface {
	Start()
}

// DiscoveryService gerencia a descoberta automática de nós via UDP
type DiscoveryService struct {
	raftAddress      string
	Port             int
	Interval         time.Duration
	Message          string
	OnPeerDiscovered func(peerRaftAddress string)
}

// NewDiscoveryService cria uma nova instância
func NewDiscoveryService(raftAddress string) *DiscoveryService {
	return &DiscoveryService{
		raftAddress:      raftAddress,
		Port:             DiscoveryPort,
		Interval:         DiscoveryInterval,
		Message:          DiscoveryMessage,
	}
}

// Start inicia as rotinas de escuta e broadcast
func (ds *DiscoveryService) Start() {
	go ds.listen()
	go ds.broadcast()
}

func (ds *DiscoveryService) listen() {
	addr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf(":%d", ds.Port))
	if err != nil {
		log.Fatalf("[DISCOVERY] Falha ao resolver endereço UDP para escuta: %v", err)
	}

	conn, err := net.ListenUDP("udp4", addr)
	if err != nil {
		log.Fatalf("[DISCOVERY] Falha ao escutar na porta UDP: %v", err)
	}
	defer conn.Close()

	buf := make([]byte, 1024)
	for {
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Printf("[DISCOVERY] Erro ao ler do UDP: %v", err)
			continue
		}

		message := string(buf[:n])
		// Extrai o endereço do Raft da mensagem, que é o payload
		if len(message) > len(ds.Message) && message[:len(ds.Message)] == ds.Message {
			peerRaftAddr := message[len(ds.Message):]

			// Não reagir às próprias mensagens
			if peerRaftAddr == ds.raftAddress {
				continue
			}

			log.Printf("[DISCOVERY] Nó par descoberto com endereço Raft: %s", peerRaftAddr)

			if ds.OnPeerDiscovered != nil {
				go ds.OnPeerDiscovered(peerRaftAddr)
			}
		}
	}
}

func (ds *DiscoveryService) broadcast() {
	broadcastAddr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("255.255.255.255:%d", ds.Port))
	if err != nil {
		log.Fatalf("[DISCOVERY] Falha ao resolver endereço de broadcast: %v", err)
	}

	conn, err := net.DialUDP("udp4", nil, broadcastAddr)
	if err != nil {
		log.Fatalf("[DISCOVERY] Falha ao conectar para broadcast: %v", err)
	}
	defer conn.Close()

	ticker := time.NewTicker(ds.Interval)
	defer ticker.Stop()

	// Mensagem a ser enviada = MagicString + nosso endereço Raft
	message := ds.Message + ds.raftAddress

	for {
		<-ticker.C
		_, err := conn.Write([]byte(message))
		if err != nil {
			log.Printf("[DISCOVERY] Falha ao enviar broadcast: %v", err)
		} else {
			// log.Printf("[DISCOVERY] Broadcast enviado.")
		}
	}
}
