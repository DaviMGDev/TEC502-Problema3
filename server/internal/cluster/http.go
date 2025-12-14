package cluster

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/raft"
)

// GinHttpTransport implementa a ClusterTransportInterface
type GinHttpTransport struct {
	bindAddress string
	nodeID      string
	router      *gin.Engine
	client      *resty.Client
	raftNode    *raft.Raft
	timeout     time.Duration
}

// NewGinHttpTransport cria a instância do transporte
func NewGinHttpTransport(bindAddress, nodeID string, raftNode *raft.Raft) ClusterTransportInterface {
	router := gin.Default()
	transport := &GinHttpTransport{
		bindAddress: bindAddress,
		nodeID:      nodeID,
		router:      router,
		client:      resty.New(),
		raftNode:    raftNode,
		timeout:     10 * time.Second,
	}
	transport.setupRoutes()
	return transport
}

func (t *GinHttpTransport) setupRoutes() {
	group := t.router.Group("/raft")
	group.POST("/join", t.handleJoin)
	group.POST("/command", t.handleCommand)
}

// Start inicia o servidor HTTP (Gin) em background
func (t *GinHttpTransport) Start() error {
	go func() {
		if err := t.router.Run(t.bindAddress); err != nil {
			// Log fatal em caso de falha ao iniciar o servidor HTTP
		}
	}()
	return nil
}

// JoinCluster é usado por um nó novo para pedir entrada no cluster
func (t *GinHttpTransport) JoinCluster(targetAddress string, myRaftID string, myRaftAddress string) error {
	req := JoinRequest{
		NodeID:      myRaftID,
		NodeAddress: myRaftAddress,
	}

	resp, err := t.client.R().
		SetBody(req).
		Post(fmt.Sprintf("http://%s/raft/join", targetAddress))

	if err != nil {
		return fmt.Errorf("falha ao enviar requisição de join para %s: %w", targetAddress, err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("erro ao entrar no cluster. status: %s, body: %s", resp.Status(), resp.String())
	}

	return nil
}

// ForwardCommand é usado por um nó seguidor para repassar um evento ao líder
func (t *GinHttpTransport) ForwardCommand(leaderAddress string, eventBytes []byte) error {
	resp, err := t.client.R().
		SetBody(bytes.NewReader(eventBytes)).
		SetHeader("Content-Type", "application/json").
		Post(fmt.Sprintf("http://%s/raft/command", leaderAddress))

	if err != nil {
		return fmt.Errorf("falha ao encaminhar comando para o líder %s: %w", leaderAddress, err)
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("erro do líder ao processar comando. status: %s, body: %s", resp.Status(), resp.String())
	}
	return nil
}

// Handlers do Gin (privados)
func (t *GinHttpTransport) handleJoin(c *gin.Context) {
	var req JoinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "corpo da requisição inválido: " + err.Error()})
		return
	}

	if t.raftNode.State() != raft.Leader {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "não sou o líder"})
		return
	}

	addVoterFuture := t.raftNode.AddVoter(raft.ServerID(req.NodeID), raft.ServerAddress(req.NodeAddress), 0, 0)
	if err := addVoterFuture.Error(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "falha ao adicionar nó ao cluster: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (t *GinHttpTransport) handleCommand(c *gin.Context) {
	if t.raftNode.State() != raft.Leader {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "não sou o líder"})
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "falha ao ler corpo da requisição: " + err.Error()})
		return
	}

	applyFuture := t.raftNode.Apply(body, t.timeout)
	if err := applyFuture.Error(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "falha ao aplicar comando no raft: " + err.Error()})
		return
	}

	// Opcional: retornar a resposta da FSM para o nó seguidor
	res := applyFuture.Response()
	if resErr, ok := res.(error); ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": resErr.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
