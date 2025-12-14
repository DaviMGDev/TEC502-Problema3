package cluster

import (
	"cod-server/internal/api"
	"encoding/json"
	"fmt"
	"io"

	raft "github.com/hashicorp/raft"
)

// ClusterFSM é a máquina de estados que traduz logs do Raft em ações do sistema
type ClusterFSM struct {
	// Dependência: A interface que sabe lidar com a lógica do jogo
	eventHandler api.EventHandlerInterface
}

// NewClusterFSM cria a FSM injetando o handler de eventos
func NewClusterFSM(handler api.EventHandlerInterface) *ClusterFSM {
	return &ClusterFSM{
		eventHandler: handler,
	}
}

// Apply é chamado pelo Raft quando um log é commitado
func (fsm *ClusterFSM) Apply(log *raft.Log) interface{} {
	var event api.Event
	if err := json.Unmarshal(log.Data, &event); err != nil {
		return fmt.Errorf("failed to unmarshal log data: %w", err)
	}

	switch event.Method {
	case "register":
		return fsm.eventHandler.OnRegister(event)
	case "login":
		return fsm.eventHandler.OnLogin(event)
	case "get_cards":
		return fsm.eventHandler.OnGetCards(event)
	case "buy_pack":
		return fsm.eventHandler.OnBuyPack(event)
	case "offer_trade":
		return fsm.eventHandler.OnOfferTrade(event)
	case "start_match":
		return fsm.eventHandler.OnStartMatch(event)
	case "join_match":
		return fsm.eventHandler.OnJoinMatch(event)
	case "surrender_match":
		return fsm.eventHandler.OnSurrenderMatch(event)
	case "make_move":
		return fsm.eventHandler.OnMakeMove(event)
	default:
		return fmt.Errorf("unhandled fsm event method: %s", event.Method)
	}
}

// Snapshot retorna um "retrato" do estado atual do sistema
func (fsm *ClusterFSM) Snapshot() (raft.FSMSnapshot, error) {
	// Lógica em Pseudocódigo:

	// IDEALMENTE: Você chamaria algo como fsm.services.GetAllState()
	// PARA AGORA (Simplificação):
	// Retorne uma struct vazia que implemente FSMSnapshot, apenas para cumprir contrato.
	// Ex: return &NoOpSnapshot{}, nil
	return nil, nil
}

// Restore restaura o estado a partir de um backup
func (fsm *ClusterFSM) Restore(rc io.ReadCloser) error {
	// Lógica em Pseudocódigo:

	// Ler os dados de rc (JSON decoder)
	// Repopular os repositórios/services com esses dados
	return nil
}
