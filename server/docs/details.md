# Análise Detalhada do Projeto

Este ficheiro documenta a análise incremental do código-fonte do servidor do jogo.

## 1. Pacote `internal/domain`

**Data da Análise:** 2025-12-08

**Ficheiros Analisados:** `card.go`, `user.go`, `match.go`, `package.go`

### Resumo

O pacote `domain` é o núcleo do sistema, definindo as entidades de negócio e as regras fundamentais do jogo. Todas as estruturas de dados principais residem aqui.

### Estruturas Principais:

*   **`Card`**: Representa uma carta de jogo.
    *   `ID`: Identificador único.
    *   `Type`: O tipo da carta (`Rock`, `Paper`, ou `Scissors`).
    *   `Level`: Um nível associado à carta, sugerindo uma mecânica de progressão.
    *   **Método `Against(*Card)`**: Contém a lógica central do jogo (Pedra-Papel-Tesoura), retornando `1`, `0`, ou `-1` para vitória, empate ou derrota.

*   **`User`**: Representa um jogador.
    *   `ID`, `Username`, `Password`: Dados de autenticação e identificação.
    *   `Cards`: Uma coleção de cartas (`utils.Map[string, *Card]`) que o jogador possui.

*   **`Match`**: Representa uma partida ou uma rodada.
    *   `ID`: Identificador da partida.
    *   `Players`: As identificações dos dois jogadores.
    *   `Moves`: As duas cartas jogadas na partida.
    *   `Winner`: A identificação do jogador vencedor.

*   **`Package`**: Representa um pacote de cartas.
    *   `ID`: Identificador do pacote.
    *   `Cards`: As cartas contidas no pacote.
    *   **Método `Unpack()`**: Retorna três cartas (pedra, papel, tesoura), o que reforça a ideia de ser um "pacote inicial" para novos jogadores.

### Dependências:

*   **`cod-server/internal/utils`**: Utiliza o tipo `utils.Map` para as coleções de cartas, o que sugere que estas coleções precisam de ser seguras para acesso concorrente (thread-safe).

### Conclusão Inicial:

Este pacote estabelece a base de um jogo de cartas com uma mecânica de Pedra-Papel-Tesoura e um sistema de progressão (`Level`). A estrutura dos dados está claramente definida e a utilização de um mapa customizado aponta para requisitos de concorrência.

## 2. Pacote `internal/utils`

**Data da Análise:** 2025-12-08

**Ficheiros Analisados:** `dict.go`, `list.go`, `map.go`, `mux.go`, `parser.go`

### Resumo

O pacote `utils` fornece um conjunto de componentes genéricos e reutilizáveis, com um foco claro em concorrência e manipulação de comandos. É a caixa de ferramentas fundamental para o resto da aplicação.

### Estruturas e Funções Principais:

*   **`SafeMap[K, V]`**: Uma implementação de mapa genérico que é seguro para acesso concorrente (thread-safe) através do uso de `sync.RWMutex`. É usado pelo pacote `domain` e é central para a arquitetura do sistema.
*   **`SafeList[T]`**: Uma implementação de lista/slice genérica, também thread-safe, que oferece operações comuns como `Append`, `Remove`, etc.
*   **`Mux[fn]`**: Um multiplexer genérico que mapeia chaves `string` (comandos) a funções de qualquer tipo. Inclui suporte para uma função `default`, sendo ideal para encaminhar comandos para os seus manipuladores.
*   **`ParseCommand(string)`**: Uma função que processa uma `string` de entrada. Se a `string` começar com `/`, a função extrai um comando e os seus argumentos. Caso contrário, devolve o comando `chat`. Isto indica a existência de uma interface de utilizador baseada em comandos de texto.
*   **`Dict`**: Um alias para `map[string]any` com métodos utilitários para converter o mapa para uma `string` ou `[]byte` em formato JSON. Útil para criar respostas de API dinâmicas.

### Conclusão Inicial:

Este pacote é a fundação técnica da aplicação. A forte aposta em estruturas thread-safe (`SafeMap`, `SafeList`) e em mecanismos de despacho de comandos (`Mux`, `ParseCommand`) sugere uma arquitetura orientada a eventos e preparada para alta concorrência, como um servidor de jogo em tempo real.

## 3. Pacote `internal/data`

**Data da Análise:** 2025-12-08

**Ficheiros Analisados:** `repository.go`

### Resumo

O pacote `data` abstrai o acesso aos dados da aplicação, separando a lógica de negócio dos detalhes de armazenamento.

### Padrão de Projeto:

*   **Repository Pattern**: O pacote implementa este padrão de forma genérica.

### Estruturas Principais:

*   **`Repository[T]` (Interface)**: Define um contrato comum para a persistência de dados, com métodos padrão de CRUD (`Create`, `Read`, `Update`, `Delete`, `List`). A natureza genérica (`[T any]`) permite que seja usado para qualquer entidade do `domain` (Users, Cards, etc.).
*   **`InMemoryRepository[T]` (Implementação)**: Uma implementação concreta da interface `Repository` que armazena os dados em memória.
    *   **Armazenamento**: Utiliza internamente o `utils.SafeMap` para guardar as entidades, garantindo que o repositório é thread-safe.
    *   **Persistência**: Os dados **não são persistentes**. Toda a informação será perdida quando o servidor for desligado.

### Conclusão Inicial:

A camada de dados está bem arquitetada. A separação da interface (`Repository`) da implementação (`InMemoryRepository`) é um ponto forte, pois permite que a aplicação seja desenvolvida e testada sem uma base de dados real. No futuro, a implementação em memória pode ser trocada por uma implementação persistente (ex: usando SQL ou uma base de dados NoSQL) com impacto mínimo no resto do código, desde que a nova implementação satisfaça a mesma interface.

## 4. Pacote `internal/state`

**Data da Análise:** 2025-12-08

**Ficheiros Analisados:** `state.go`

### Resumo

O pacote `state` serve como um contentor para configurações e variáveis de estado globais da aplicação. A sua função é centralizar informações que possam ser necessárias em diferentes partes do sistema.

### Estruturas Principais:

*   **`State`**: Contém informações de configuração e estado para a instância do servidor.
    *   `Address`: O endereço de rede da instância do servidor.
    *   `BrokerAddress`: O endereço do message broker (MQTT) ao qual o servidor se conecta.

### Conclusão Inicial:

Este pacote fornece um objeto de estado centralizado. Atualmente, armazena principalmente configurações de rede, mas foi projetado para ser expandido no futuro para guardar outras variáveis globais ou instâncias de singletons, disponibilizando-as de forma controlada para as camadas que delas necessitem (como o `gateway`).

## 5. Pacote `internal/services`

**Data da Análise:** 2025-12-08

**Ficheiros Analisados:** `users.go`, `cards.go`, `game.go`

### Resumo

O pacote `services` implementa a lógica de negócio da aplicação. Ele atua como um orquestrador, utilizando os pacotes `domain`, `data` e `utils` para executar as funcionalidades principais do jogo. A separação da lógica em múltiplos serviços (`User`, `Card`, `Game`) demonstra uma boa organização do código.

### Padrão de Projeto:

*   **Service Layer**: Cada ficheiro define um serviço que agrupa operações de negócio relacionadas.
*   **Injeção de Dependências**: Os serviços recebem as suas dependências (ex: `data.Repository`) nos seus construtores, o que facilita os testes e a manutenção.

### Análise dos Serviços:

*   **`UserService` (`users.go`)**:
    *   **Funcionalidades**: `Register` e `Login`.
    *   **Lógica**: Cria novos utilizadores e valida as suas credenciais, interagindo com o repositório de utilizadores. A geração de IDs é feita com um contador simples.

*   **`CardService` (`cards.go`)**:
    *   **Funcionalidades**: `BuyCardPack`, `SwapCard`, `ListUserCards`.
    *   **Lógica**: Permite aos utilizadores obterem e trocarem cartas.
    *   **Detalhe Chave**: Contém um método `checkStock` que gera novos pacotes de cartas com níveis aleatórios quando o "stock" em memória fica baixo. Isto revela como funciona a economia de cartas do jogo.

*   **`GameService` (`game.go`)**:
    *   **Funcionalidades**: A interface define as ações de uma partida: `StartGame`, `MakeMove`, `GetGameState`, `PLayerSurrender`.
    *   **Estado Atual**: **FUNCIONALIDADE INCOMPLETA**. A implementação deste serviço está praticamente toda comentada ou vazia. A lógica para iniciar uma partida, fazer uma jogada e determinar um vencedor não existe.
    *   **Potencial**: A estrutura do serviço contém um `queue chan *domain.Match`, sugerindo que o processamento de partidas foi pensado para ser assíncrono, mas esta funcionalidade não está implementada.

### Conclusão Inicial:

A camada de serviços está bem desenhada e implementa com sucesso a gestão de utilizadores e de cartas. No entanto, a funcionalidade mais crítica para um jogo — a própria jogabilidade (`GameService`) — está por fazer. Este é o ponto mais importante para uma futura lista de TODOs. A arquitetura está no lugar, mas o código principal do jogo precisa de ser escrito.

## 6. Pacote `internal/handlers`

**Data da Análise:** 2025-12-08

**Ficheiros Analisados:** `generic_handlers.go`

### Resumo

O pacote `handlers` define a camada que conecta os eventos de protocolo de baixo nível (como uma mensagem MQTT recebida) à lógica de negócio de alto nível (a camada de `services`). Atua como a camada "Controller" da aplicação.

### Estruturas Principais:

*   **`Handlers` (Interface)**: Define um contrato com um método para cada ação possível do utilizador (`OnRegisterEvent`, `OnLoginEvent`, `OnStartMatchEvent`, etc.). Este padrão torna claro que eventos a aplicação suporta.
*   **`HandlersImplementation` (Struct)**: A estrutura que deveria implementar a interface `Handlers`. Ela contém instâncias de todos os serviços (`UserService`, `CardService`, `GameService`) como dependências, que lhe são fornecidas através de injeção no construtor `NewHandlers`.

### Interação e Protocolo:

*   Cada método da interface aceita um `protocol.Event` e retorna um `protocol.Event`. Isto impõe um padrão de comunicação síncrona (pedido/resposta) para todas as interações e aponta para a importância do pacote `internal/api/protocol`.

### Estado Atual:

*   **FUNCIONALIDADE INCOMPLETA**: Assim como o `GameService`, este pacote é um esqueleto. A `struct` `HandlersImplementation` existe, mas **não implementa nenhum dos métodos** da interface `Handlers`. A lógica que extrai dados do evento de entrada e chama o serviço apropriado está completamente ausente.

### Conclusão Inicial:

A arquitetura para manipular eventos está bem definida, mas não implementada. É mais uma peça fundamental do puzzle que está em falta. Sem a implementação destes manipuladores, o servidor não pode responder a nenhum comando, mesmo àqueles cujos serviços (`UserService`, `CardService`) estão prontos. A próxima análise, do pacote `gateway`, deve revelar como estes `handlers` são chamados.

## 7. Pacote `internal/gateway`

**Data da Análise:** 2025-12-08

**Ficheiros Analisados:** `orchestrator.go`

### Resumo

O pacote `gateway` serve como um despachante de eventos central. É a camada que se situa entre a API de rede e os manipuladores de lógica de negócio, atuando como um ponto de entrada único para todos os eventos.

### Estruturas Principais:

*   **`Orchestrator`**:
    *   **Composição**: A estrutura é composta pela interface `handlers.Handlers` (embutida) e por uma referência ao `*state.State` global.
    *   **Implementação "Pass-Through"**: No seu estado atual, o `Orchestrator` implementa os métodos da interface `handlers.Handlers`, mas apenas para passar a chamada diretamente para a interface embutida.

### Implicações Arquiteturais:

Esta camada "extra", embora atualmente pareça redundante, serve a um propósito de design importante. É um local preparado para a implementação de "middleware" ou lógica transversal.

*   **Exemplo de Utilização Futura**: Se fosse necessário adicionar logging, métricas ou validações de segurança a todos os eventos antes de serem processados, o `Orchestrator` seria o local ideal para adicionar essa lógica, sem poluir a camada de `handlers` ou de `services`. Por ter acesso ao `state`, ele pode tomar decisões com base na configuração global da aplicação.

### Conclusão Inicial:

O `Orchestrator` atua como um ponto de entrada limpo e um potencial "intercetor" de eventos. O seu design, embora simples na implementação atual, permite a extensão futura com funcionalidades transversais (como logging ou validação) de forma elegante e centralizada.

## 8. Pacote `internal/api` e Subpacotes

### 8.1. Subpacote `internal/api/protocol`

**Data da Análise:** 2025-12-08

**Ficheiros Analisados:** `event.go`

#### Resumo

Este pacote define o formato de dados padrão para toda a comunicação orientada a eventos no sistema.

#### Estruturas Principais:

*   **`Event`**: A estrutura de mensagem unificada.
    *   `Method` (string): Identifica o comando a ser executado (ex: "login", "startMatch"). Esta é a chave de roteamento para o `utils.Mux`.
    *   `Timestamp` (time.Time): Data e hora da criação do evento.
    *   `Payload` (utils.Dict): Um mapa flexível (`map[string]any`) para transportar os dados específicos do evento, permitindo que qualquer tipo de dados JSON seja comunicado.

#### Conclusão Inicial:

O `protocol.Event` é a espinha dorsal da comunicação interna. Ao padronizar o formato da mensagem, o sistema garante que as camadas de negócio (`gateway`, `handlers`) são agnósticas em relação à fonte do evento (MQTT, REST, etc.), o que é um excelente design para a manutenibilidade e extensibilidade.

### 8.2. Subpacote `internal/api/codmqtt`

**Data da Análise:** 2025-12-08

**Ficheiros Analisados:** `handlers.go`

#### Resumo

Este pacote atua como um "Adaptador de Protocolo" para MQTT. Ele é responsável por traduzir o mundo do MQTT (tópicos, mensagens em bytes) para o mundo da aplicação (objetos `protocol.Event`).

#### Estruturas e Funções Principais:

*   **`MQTTHandler`**: Contém um cliente MQTT e uma instância da interface `handlers.Handlers` (que na prática será o `gateway.Orchestrator`).
*   **`wrap(HandlerFunc)`**: Uma função de ordem superior que adapta a interface `handlers.Handlers` para a interface `mqtt.MessageHandler`. Este wrapper:
    1.  **Desserializa** a mensagem MQTT recebida para um `protocol.Event`.
    2.  **Chama** a função de negócio correspondente (que desce pela cadeia `Orchestrator` -> `Handlers` -> `Services`).
    3.  **Publica** o `protocol.Event` de resposta de volta ao broker MQTT.
*   **`SubscribeToEvents()`**: O ponto de entrada da API. Mapeia tópicos MQTT (ex: `"cod/request/register"`) para as funções de negócio correspondentes, já adaptadas pela função `wrap`, e subscreve-os no broker.
*   **`InferEventTopic()`**: **FUNCIONALIDADE INCOMPLETA/COM BUG**. Esta função deveria determinar o tópico de resposta correto com base no evento, mas atualmente retorna sempre `"unknown"`. Isto impede que as respostas sejam enviadas para o destino certo.

#### Conclusão Inicial:

A camada de MQTT está bem desenhada, servindo como uma ponte limpa entre a rede e a aplicação. O bug na função `InferEventTopic` é crítico para a comunicação de resposta. O subpacote `rest`, que se encontra vazio, indica que uma API REST foi considerada mas não implementada.

## 9. Pacote `cmd`

**Data da Análise:** 2025-12-08

**Ficheiros Analisados:** `main.go`

### Resumo

A análise deste pacote resultou na descoberta mais crítica de todo o projeto.

### Detalhes:

*   **`main.go` está vazio**: O ficheiro que deveria conter o ponto de entrada da aplicação (`func main`) está completamente em branco.

### Conclusão Inicial:

A aplicação não tem um ponto de entrada e, portanto, **não pode ser compilada ou executada**. Apesar de ter uma arquitetura interna bem definida e em camadas, falta a "cola" final que instancia e conecta todos os componentes (`data`, `services`, `handlers`, `gateway`, `api`). O projeto, no seu estado atual, é um esqueleto não funcional. É como ter o motor, o chassi e as rodas de um carro, mas sem nenhuma peça que os ligue.

## 10. API de Eventos (`events.md`)

**Data da Análise:** 2025-12-08

**Ficheiros Analisados:** `events.md`

### Resumo

O arquivo `events.md` define o contrato da API para a comunicação baseada em eventos via MQTT. Ele especifica a estrutura exata dos payloads JSON para cada tipo de requisição e resposta, garantindo a interoperabilidade entre cliente e servidor.

### Estrutura Geral

*   **Requisições (Requests):** Enviadas para tópicos como `{feature}/{method}/requests`.
    *   `method`: `string` - O nome da ação a ser executada.
    *   `timestamp`: `string` - Data/hora do evento.
    *   `payload`: `object` - Um dicionário contendo os dados da requisição.

*   **Respostas (Responses):** Publicadas em tópicos de resposta.
    *   `method`: `string` - O método da requisição original.
    *   `timestamp`: `string` - Data/hora da resposta.
    *   `status`: `string` - `"success"` ou `"error"`.
    *   `payload`: `object` - Contém uma `status_message` e outros dados de resposta.

### Payloads de Requisição Detalhados

*   `register`: `{ "username": "<string>", "password": "<string>" }`
*   `login`: `{ "username": "<string>", "password": "<string>" }`
*   `start_game`: `{ "user_id": "<string>" }`
*   `play`: `{ "user_id": "<string>", "move": "<string>", "match_id": "<string>" }`
*   `surrender`: `{ "user_id": "<string>", "match_id": "<string>" }`
*   `list_cards`: `{ "user_id": "<string>" }`
*   `buy_pack`: `{ "user_id": "<string>" }`
*   `trade`: `{ "trader_id": "<string>", "card_type": "<string>", "username": "<string>" }`

### Conclusão

A especificação em `events.md` é a "fonte da verdade" para a camada de `handlers`. A lógica dos handlers deve validar os payloads de entrada contra esta especificação e formatar as respostas de acordo. Qualquer divergência entre o código e este arquivo deve ser corrigida, tratando o `events.md` como o contrato a ser seguido.
