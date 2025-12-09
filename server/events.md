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
        "password": "secure_password"
    }
}

{
    "method": "login",
    "timestamp": "2024-06-15T12:05:00Z",
    "payload": {
        "username": "existing_user",
        "password": "user_password"
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

# Responses on mqtt (topics can vary based on method)
```json
{
    "method": "<method>",
    "timestamp": "<timestamp>",
    "status": "success || error",
    "payload": {
        "status_message": "<message>",
        "<key>": "<value>"
    }
}

{
    "method": "register",
    "timestamp": "2024-06-15T12:00:01Z",
    "status": "success || error",
    "payload": {
        "status_message": "<message>",
        "user_id": "12345 || empty"
    }
}

{
    "method": "login",
    "timestamp": "2024-06-15T12:05:01Z",
    "status": "success || error",
    "payload": {
        "status_message": "<message>",
        "user_id": "12345 || empty"
    }
}

{
    "method": "start_game",
    "timestamp": "2024-06-15T12:10:01Z",
    "status": "success || error",
    "payload": {
        "status_message": "<message>",
        "match_id": "67890 || empty"
    }
}

{
    "method": "play",
    "timestamp": "2024-06-15T12:15:01Z",
    "status": "success || error",
    "payload": {
        "status_message": "<message>",
        "round_result": "win || lose || draw || empty"
    }
}

{
    "method": "surrender",
    "timestamp": "2024-06-15T12:20:01Z",
    "status": "success || error",
    "payload": {
        "status_message": "<message>"
    }
}

{
    "method": "list_cards",
    "timestamp": "2024-06-15T12:25:01Z",
    "status": "success || error",
    "payload": {
        "status_message": "<message>",
        "cards": ["card1", "card2", "..."] || []
    }
}

{
    "method": "buy_pack",
    "timestamp": "2024-06-15T12:30:01Z",
    "status": "success || error",
    "payload": {
        "status_message": "<message>",
        "new_cards": ["cardA", "cardB", "..."] || []
    }
}

{
    "method": "trade",
    "timestamp": "2024-06-15T12:35:01Z",
    "status": "success || error",
    "payload": {
        "status_message": "<message>"
    }
}
```
