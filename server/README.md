# Projeto do Servidor do Jogo de Cartas

## 1. Estrutura do Projeto

A seguinte estrutura de diretórios organiza os componentes do servidor:

```
/home/davi/Distrobox/arch-box/davi/tmp/TEC502-Problema3/server/
├───go.mod
├───go.sum
├───cmd/
│   └───main.go
└───internal/
    ├───api/
    │   ├───codmqtt/
    │   │   └───handlers.go
    │   ├───protocol/
    │   │   └───event.go
    │   └───rest/
    ├───data/
    │   └───repository.go
    ├───domain/
    │   ├───card.go
    │   ├───match.go
    │   ├───package.go
    │   └───user.go
    ├───gateway/
    │   └───orchestrator.go
    ├───handlers/
    │   └───generic_handlers.go
    ├───services/
    │   ├───cards.go
    │   ├───game.go
    │   └───users.go
    ├───state/
    │   └───state.go
    └───utils/
        ├───dict.go
        ├───list.go
        ├───map.go
        └───mux.go
        └───parser.go
```

## 2. Diagrama de Arquitetura

O sistema segue uma arquitetura baseada em camadas e orientada a eventos. Um evento (mensagem) flui através das camadas, sendo processado em cada uma delas.

```
+----------------+
|  Cliente MQTT  |
+----------------+
        | (protocol.Event)
        v
+-----------------------------+
|    API Layer (internal/api) |
|                             |
|  +-----------------------+  |
|  | codmqtt.MQTTHandler   |  | (Deserializa Evento)
|  | (Adapta MQTT p/ Event) |  | (Publica Resposta)
|  +-----------------------+  |
+-----------------------------+
        | (chama handlers.Handlers)
        v
+-----------------------------+
|    Gateway (internal/gateway)|
|                             |
|  +-----------------------+  |
|  |  Orchestrator         |  |
|  |  (usa state.State)      |  | (Pass-through para Handlers)
|  +-----------------------+  |
+-----------------------------+
        | (chama services.*Service)
        v
+-----------------------------+
|   Handlers (internal/handlers) |
|                             |
|  +-----------------------+  |
|  |  HandlersImplementation |  | (Traduz Evento p/ Chamada Service)
|  +-----------------------+  |
+-----------------------------+
        | (chama data.Repository)
        v
+-----------------------------+
|   Services (internal/services) |
|                             |
|  +-----------------------+  |
|  |  UserService          |  | (Lógica de Negócio p/ User)
|  |  CardService          |  | (Lógica de Negócio p/ Cards)
|  |  GameService          |  | (Lógica de Negócio p/ Jogo)
|  +-----------------------+  |
+-----------------------------+
        | (opera no armazenamento)
        v
+-----------------------------+
|    Data Layer (internal/data) |
|                             |
|  +-----------------------+  |
|  |  InMemoryRepository   |  | (Armazenamento Thread-safe em Memória)
|  +-----------------------+  |
+-----------------------------+
```

Nota:
- `protocol.Event` é o formato de mensagem padrão.
- `utils.SafeMap`, `utils.SafeList`, `utils.Mux` são usados por várias camadas para concorrência e despacho.
- `state.State` é usado pelo `Orchestrator` para aceder a estado e configuração globais (ex: endereço do broker MQTT).

## 3. Checklist de TODOs

Esta secção lista as funcionalidades incompletas, bugs e lacunas arquitetónicas identificadas durante a análise do projeto. Estas são tarefas cruciais para tornar a aplicação funcional e robusta.

-   [ ] **CRÍTICO: Criar o ponto de entrada da aplicação (`func main`) em `cmd/main.go`**: Atualmente, o ficheiro `main.go` está vazio, tornando a aplicação não executável. Esta é a tarefa de maior prioridade.
-   [ ] **Implementar a lógica dos manipuladores de eventos em `internal/handlers/generic_handlers.go`**: Os métodos da `HandlersImplementation` estão definidos na interface, mas a sua lógica para extrair dados do `protocol.Event` recebido e chamar os serviços correspondentes está ausente.
-   [ ] **Implementar a lógica de jogo em `internal/services/game.go`**: As funções principais do `GameService` (`StartGame`, `MakeMove`, `GetGameState`, `PLayerSurrender`) estão vazias ou comentadas. Esta é uma funcionalidade central do jogo que precisa ser desenvolvida para que o jogo possa ser jogado.
-   [ ] **Implementar o Nível da Carta como critério de desempate**: Na lógica de combate (`domain.Card.Against`), em caso de empate (ex: Pedra vs. Pedra), o `Level` da carta deve ser usado para determinar o vencedor. A carta de maior nível vence. Esta lógica precisa ser adicionada.
-   [ ] **Refinar a lógica de `MakeMove` para resolver a partida**: O método `MakeMove` no `GameService` deve ser implementado para: 1. Guardar a jogada do jogador. 2. Verificar se o oponente também já jogou. 3. Se sim, invocar um método interno para comparar as duas jogadas, determinar o vencedor e finalizar a partida.
-   [ ] **Corrigir o bug em `internal/api/codmqtt/handlers.go` na função `InferEventTopic`**: Atualmente, esta função retorna sempre `"unknown"`, o que impede que as respostas do servidor sejam publicadas nos tópicos MQTT corretos para os clientes. É necessário implementar a lógica para inferir o tópico de resposta adequado com base no método do evento original.
-   [ ] **Adicionar persistência de dados real**: A implementação atual (`data.InMemoryRepository`) armazena os dados apenas em memória RAM. Todos os dados são perdidos ao reiniciar o servidor. É necessário integrar uma solução de persistência (ex: banco de dados SQL como PostgreSQL, ou NoSQL como MongoDB/Redis) para garantir que os dados dos utilizadores e do jogo sejam guardados permanentemente.
-   [ ] **Desenvolver a API REST em `internal/api/rest`**: O diretório `internal/api/rest` está vazio, indicando que uma API REST foi planeada, mas nunca implementada. Se desejado, esta API deve ser desenvolvida para oferecer uma alternativa ou complemento à interface MQTT.
-   [ ] **Melhorar a geração de IDs**: Atualmente, a geração de IDs utiliza contadores simples que podem não ser robustos o suficiente em um ambiente distribuído ou para evitar colisões. Considerar o uso de UUIDs (Universally Unique Identifiers) para garantir unicidade e escalabilidade.

## 4. Como Executar

Para compilar e executar o servidor, os seguintes passos seriam necessários APÓS a implementação das partes em falta (nomeadamente o `main.go` e a lógica dos `handlers`):

1.  **Pré-requisitos**:
    *   Ter o Go instalado (versão 1.25.4 ou superior, conforme `go.mod`).
    *   Um broker MQTT acessível (ex: Mosquitto) para comunicação.

2.  **Compilação**:
    ```bash
    go build -o server ./cmd
    ```

3.  **Execução**:
    ```bash
    ./server
    ```
    (Ou `go run ./cmd` para executar diretamente sem compilação prévia)

## 5. Dependências

As dependências do projeto são geridas pelo `go modules`.

### Dependências Diretas:

*   `github.com/eclipse/paho.mqtt.golang v1.5.1` (Cliente MQTT para comunicação)

### Dependências Indiretas:

*   `github.com/gorilla/websocket v1.5.3` (Provavelmente para uma futura implementação REST/WebSocket)
*   `golang.org/x/net v0.44.0`
*   `golang.org/x/sync v0.17.0`

