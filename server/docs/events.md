# Requests on mqtt topic "{feature}/{method}/requests"
```json
{
    "method": "<method>",
    "timestamp": "<timestamp>",
    "payload": {
        "<key>": "<value>"
    }
}

{
    "method": "register",
    "timestamp": "2024-06-15T12:00:00Z",
    "payload": {
        "username": "new_user",
        "password": "secure_password",
        "client_id": "unique_client_identifier"
    }
}

{
    "method": "login",
    "timestamp": "2024-06-15T12:05:00Z",
    "payload": {
        "username": "existing_user",
        "password": "user_password",
        "client_id": "unique_client_identifier"
    }
}

{
    "method": "start_game",
    "timestamp": "2024-06-15T12:10:00Z",
    "payload": {
        "user_id": "12345",
    }
}

{
    "method": "play",
    "timestamp": "2024-06-15T12:15:00Z",
    "payload": {
        "user_id": "12345",
        "move": "rock",
        "match_id": "67890",
    }
}

{
    "method": "surrender",
    "timestamp": "2024-06-15T12:20:00Z",
    "payload": {
        "user_id": "12345",
        "match_id": "67890"
    }
}

{
    "method": "list_cards",
    "timestamp": "2024-06-15T12:25:00Z",
    "payload": {
        "user_id": "12345"
    }
}

{
    "method": "buy_pack",
    "timestamp": "2024-06-15T12:30:00Z",
    "payload": {
        "user_id": "12345",
    }
}

{
    "method": "trade",
    "timestamp": "2024-06-15T12:35:00Z",
    "payload": {
        "trader_id": "12345",
        "card_type": "paper",
        "username": "trading_user"
    }
}

```

# Responses on mqtt (each client subscribes to its own response topics)

As respostas são publicadas em tópicos específicos, construídos para direcionar a mensagem de volta ao cliente ou usuário correto.

**Estrutura Geral do Tópico de Resposta:** `cod/response/{feature}/{method}/{identifier}`

*   `cod/response/users/{method}/{client_id}`: Para eventos de `register` e `login`.
*   `cod/response/game/{method}/{user_id}`: Para eventos `start_game`.
*   `cod/response/game/{method}/{match_id}/{user_id}`: Para eventos `play` (ou `move`) e `surrender`.
*   `cod/response/cards/{method}/{user_id}`: Para eventos `list_cards`, `buy_pack`, `trade`.

**Formato do Payload de Resposta:**
```json
{
    "method": "<method>",
    "timestamp": "<timestamp>",
    "payload": {
        "status": "success || error",
        "status_message": "<message>",
        "<key>": "<value>"
    }
}

{
    "method": "register",
    "timestamp": "2024-06-15T12:00:01Z",
    "payload": {
        "status": "success || error",
        "status_message": "<message>",
        "user_id": "12345 || empty"
    }
}

{
    "method": "login",
    "timestamp": "2024-06-15T12:05:01Z",
    "payload": {
        "status": "success || error",
        "status_message": "<message>",
        "user_id": "12345 || empty"
    }
}

{
    "method": "start_game",
    "timestamp": "2024-06-15T12:10:01Z",
    "payload": {
        "status": "success || error",
        "status_message": "<message>",
        "match_id": "67890 || empty"
    }
}

{
    "method": "play",
    "timestamp": "2024-06-15T12:15:01Z",
    "payload": {
        "status": "success || error",
        "status_message": "<message>",
        "round_result": "win || lose || draw || empty"
    }
}

{
    "method": "surrender",
    "timestamp": "2024-06-15T12:20:01Z",
    "payload": {
        "status": "success || error",
        "status_message": "<message>"
    }
}

{
    "method": "list_cards",
    "timestamp": "2024-06-15T12:25:01Z",
    "payload": {
        "status": "success || error",
        "status_message": "<message>",
        "cards": ["card1", "card2", "..."] || []
    }
}

{
    "method": "buy_pack",
    "timestamp": "2024-06-15T12:30:01Z",
    "payload": {
        "status": "success || error",
        "status_message": "<message>",
        "new_cards": ["cardA", "cardB", "..."] || []
    }
}

{
    "method": "trade",
    "timestamp": "2024-06-15T12:35:01Z",
    "payload": {
        "status": "success || error",
        "status_message": "<message>"
    }
}
```
