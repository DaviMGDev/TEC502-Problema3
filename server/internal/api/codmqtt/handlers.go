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
        "cod/request/register":      mh.wrap(mh.handlers.OnRegisterEvent),
        "cod/request/login":         mh.wrap(mh.handlers.OnLoginEvent),
		"cod/request/buy_card_pack": mh.wrap(mh.handlers.OnBuyCardPackEvent),
		"cod/request/trade":         mh.wrap(mh.handlers.OnTradeEvent), // Renamed from swap_card
		"cod/request/list_user_cards": mh.wrap(mh.handlers.OnListUserCardsEvent),
		"cod/request/start_match":   mh.wrap(mh.handlers.OnStartMatchEvent),
		"cod/request/play":          mh.wrap(mh.handlers.OnPlayCardEvent), // Renamed from play_card
		"cod/request/surrender":     mh.wrap(mh.handlers.OnPlayerSurrenderEvent), // Renamed from PLayerSurrenderEvent
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
	var feature string
	var identifier string

	// Determine the feature based on the method
	switch event.Method {
	case "register", "login":
		feature = "users"
	case "start_game", "play", "surrender":
		feature = "game"
	case "list_cards", "buy_pack", "trade":
		feature = "cards"
	default:
		// Fallback for unknown methods, should ideally not happen if all methods are covered
		return "cod/response/unknown_feature/unknown_method/unknown_identifier"
	}

	// Determine the identifier based on the method and available payload fields
	switch event.Method {
	case "register", "login":
		if clientID, ok := event.Payload["client_id"].(string); ok && clientID != "" {
			identifier = clientID
		} else {
			identifier = "fallback_client_id" // Use a generic ID if not found, though client_id should be mandatory
		}
	case "start_game", "list_cards", "buy_pack", "trade": // These use user_id as the primary identifier
		if userID, ok := event.Payload["user_id"].(string); ok && userID != "" {
			identifier = userID
		} else {
			identifier = "fallback_user_id" // Use a generic ID if not found
		}
	case "play", "surrender": // These require both match_id and user_id
		matchID, matchOK := event.Payload["match_id"].(string)
		userID, userOK := event.Payload["user_id"].(string) // For 'play' event, user_id is the player making the move
		if matchOK && matchID != "" && userOK && userID != "" {
			identifier = fmt.Sprintf("%s/%s", matchID, userID)
		} else {
			// If matchID or userID is missing, use fallbacks
			if !matchOK || matchID == "" { matchID = "fallback_match" }
			if !userOK || userID == "" { userID = "fallback_user" }
			identifier = fmt.Sprintf("%s/%s", matchID, userID)
		}
	default:
		identifier = "unknown_identifier"
	}
	
	// Construct the final topic string
	// The problem states the topic should be `cod/response/{feature}/{method}/{identifier}`
	return fmt.Sprintf("cod/response/%s/%s/%s", feature, event.Method, identifier)
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

