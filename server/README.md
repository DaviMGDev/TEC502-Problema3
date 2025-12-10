# Servidor do Jogo "Cards of Destiny"

Este projeto é o backend para "Cards of Destiny", um jogo de cartas colecionáveis multiplayer. O servidor é construído em Go e utiliza uma arquitetura em camadas orientada a eventos. O estado atual do projeto é de um esqueleto funcional: a lógica de negócio principal está implementada e testada, mas a camada de comunicação (handlers) e o ponto de entrada da aplicação ainda estão pendentes.

## 1. Estrutura do Projeto

A seguinte estrutura de diretórios organiza os componentes do servidor:

```
/
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
        └───map.go
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

## 3. Estado do Projeto e Próximos Passos

### Itens Concluídos (Recentemente)
*   `[x]` Implementada toda a lógica de negócio na camada de Serviços (`UserService`, `CardService`, `GameService`).
*   `[x]` Implementado sistema de fila (matchmaking) no `GameService`.
*   `[x]` Adicionados testes unitários e de integração para toda a camada de `Services`, que estão passando.
*   `[x]` Refatorados os testes de serviço para usar `InMemoryRepository` em vez de mocks manuais.
*   `[x]` Criados arquivos de teste placeholder para todos os pacotes do projeto.
*   `[x]` Refatorados os testes dos `Handlers` para seguir a filosofia TDD corretamente e falhar de forma explícita.
*   `[x]` Corrigido bug no repositório que permitia a criação de usuários duplicados.

### Próximos Passos para o MVP (Prioridade Alta)
*   `[ ]` **Implementar a Lógica dos Handlers:** Implementar a lógica em `internal/handlers/generic_handlers.go` para fazer os testes TDD passarem.
*   `[ ]` **Corrigir Bug Crítico da API:** Implementar a função `InferEventTopic` em `internal/api/codmqtt/handlers.go` para permitir o envio de respostas MQTT.
*   `[ ]` **Implementar Ponto de Entrada:** Escrever o código em `cmd/main.go` para instanciar e conectar todos os componentes e iniciar o servidor.

### Melhorias Futuras (Pós-MVP)
*   `[ ]` Adicionar persistência de dados com um banco de dados real (ex: PostgreSQL).
*   `[ ]` Desenvolver a API REST planejada em `internal/api/rest`.
*   `[ ]` Escrever testes funcionais para os pacotes com placeholders (`gateway`, `codmqtt`, etc.).
*   `[ ]` Melhorar a geração de IDs (ex: usando UUIDs).

## 4. Documentação de Arquitetura e API

*   **`details.md`**: Contém uma análise detalhada e incremental de cada pacote do projeto.
*   **`events.md`**: Descreve o contrato da API, com a especificação JSON para cada evento de requisição e resposta, **incluindo a estratégia de tópicos de resposta que utiliza `client_id` e `user_id` para direcionamento e a estrutura de respostas com `status` dentro do `payload`.**
*   **`game.md`**: Explica as regras e o fluxo do jogo "Cards of Destiny".
*   **`tests.md`**: Descreve a metodologia e a filosofia de testes adotadas no projeto.

## 5. Como Executar e Testar

### 5.1. Pré-requisitos

*   Go (versão definida em `go.mod`).
*   Um broker MQTT acessível (ex: Mosquitto) para comunicação, quando o `main.go` for implementado.

### 5.2. Execução de Testes

A ferramenta de teste padrão do Go é o `go test`. Para ter uma visão completa do estado do projeto, execute o seguinte comando na raiz:

```bash
# Executa todos os testes em todos os pacotes
go test -v ./...
```

Isso irá mostrar:
*   **PASS** para os pacotes `domain`, `utils`, `data` e `services`.
*   **FAIL** para `handlers` (pois os testes TDD esperam uma implementação que ainda não existe).
*   **FAIL** para `gateway`, `api/*` e `state` (pois os testes são placeholders que precisam ser implementados).

### 5.3. Compilação e Execução do Servidor

**Nota:** O servidor não é executável até que o **Passo 3.1** do plano do MVP seja concluído.

```bash
# Após implementar cmd/main.go
go build -o server ./cmd
./server
```

## 6. Dependências

As dependências do projeto são geridas pelo `go modules` e podem ser encontradas no arquivo `go.mod`.