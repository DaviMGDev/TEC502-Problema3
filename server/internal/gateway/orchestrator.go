package gateway

import (
	"cod-server/internal/api/protocol"
	"cod-server/internal/handlers"
	"cod-server/internal/state"
)

type Orchestrator struct {
	handlers.Handlers
	state	*state.State
}

func (orch *Orchestrator) OnRegisterEvent(event protocol.Event) protocol.Event {
	return orch.Handlers.OnRegisterEvent(event)
}

func (orch *Orchestrator) OnLoginEvent(event protocol.Event) protocol.Event {
	return orch.Handlers.OnLoginEvent(event)
}

func (orch *Orchestrator) OnBuyCardPackEvent(event protocol.Event) protocol.Event {
	return orch.Handlers.OnBuyCardPackEvent(event)
}

func (orch *Orchestrator) OnTradeEvent(event protocol.Event) protocol.Event {
	return orch.Handlers.OnTradeEvent(event)
}

func (orch *Orchestrator) OnListUserCardsEvent(event protocol.Event) protocol.Event {
	return orch.Handlers.OnListUserCardsEvent(event)
}

func (orch *Orchestrator) OnStartMatchEvent(event protocol.Event) protocol.Event {
	return orch.Handlers.OnStartMatchEvent(event)
}

func (orch *Orchestrator) OnPlayCardEvent(event protocol.Event) protocol.Event {
	return orch.Handlers.OnPlayCardEvent(event)
}

func (orch *Orchestrator) OnGetMatchResultEvent(event protocol.Event) protocol.Event {
	return orch.Handlers.OnGetMatchResultEvent(event)
}
