package codmqtt

import (
	"cod-server/internal/api/protocol"
	"cod-server/internal/handlers"
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTHandler struct {
	client   mqtt.Client
	handlers handlers.Handlers
} 

func NewMQTTHandler(client mqtt.Client, handlers handlers.Handlers) *MQTTHandler {
	return &MQTTHandler{
		client:   client,
		handlers: handlers,
	}
}
// Definimos o tipo da função que a interface Handlers usa
type HandlerFunc func(protocol.Event) protocol.Event

// wrap cria um mqtt.MessageHandler padrão para qualquer função de negócio
func (mh *MQTTHandler) wrap(businessLogic HandlerFunc) mqtt.MessageHandler {
    return func(client mqtt.Client, msg mqtt.Message) {
        // 1. Deserialização (Bytes -> Struct)
        var reqEvent protocol.Event
        if err := json.Unmarshal(msg.Payload(), &reqEvent); err != nil {
            fmt.Printf("Erro ao decodificar JSON: %v\n", err)
            return // Ou publicar um erro genérico de volta
        }

        // 2. Execução (Chama a interface cega)
        // Aqui a mágica acontece: vai pro Orquestrador -> Lógica
        resEvent := businessLogic(reqEvent)

        // 3. Resposta (Struct -> Bytes -> Publicar)
				mh.PublishEvent(resEvent)
    }
}

func (mh *MQTTHandler) SubscribeToEvents() error {
    routes := map[string]mqtt.MessageHandler{
        "cod/request/register": mh.wrap(mh.handlers.OnRegisterEvent),
        "cod/request/login":    mh.wrap(mh.handlers.OnLoginEvent),
				"cod/request/buy_card_pack":   mh.wrap(mh.handlers.OnBuyCardPackEvent),
				"cod/request/swap_card":       mh.wrap(mh.handlers.OnSwapCardEvent),
				"cod/request/list_user_cards": mh.wrap(mh.handlers.OnListUserCardsEvent),
				"cod/request/start_match":     mh.wrap(mh.handlers.OnStartMatchEvent),
				"cod/request/play_card":       mh.wrap(mh.handlers.OnPlayCardEvent),
				"cod/request/get_match_result": mh.wrap(mh.handlers.OnGetMatchResultEvent),
    }

    for topic, handler := range routes {
        token := mh.client.Subscribe(topic, 0, handler)
        token.Wait()
        if token.Error() != nil {
            return token.Error()
        }
    }
    return nil
}

func InferEventTopic(event protocol.Event) string {
	switch event.Method {
	default:
		return "unknown"
	}
}

func (mh *MQTTHandler) PublishEvent(event protocol.Event) error {
	topic := InferEventTopic(event)
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}
	token := mh.client.Publish(topic, 0, false, payload)
	token.Wait()
	return token.Error()
}

