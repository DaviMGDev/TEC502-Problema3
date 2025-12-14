# Análise Conceitual e Arquitetural do Projeto

Este documento detalha a arquitetura intencional, a estrutura e o estado de desenvolvimento atual do projeto, com base em uma análise de alto nível de seus arquivos-fonte. O objetivo é criar um mapa conceitual que sirva de base para futuras contribuições.

---

## `cmd/main.go`

### Função Prevista
Este arquivo é o ponto de entrada (`main`) da aplicação. Sua responsabilidade é orquestrar a inicialização e a interconexão de todos os módulos do sistema, incluindo a camada de dados, serviços de domínio, o cluster Raft, a API interna (HTTP), o cliente MQTT para entrada de eventos e o serviço de descoberta de nós. Ele funciona como um "maestro" que monta a aplicação.

### Componentes Notáveis (Existentes ou Implícitos)
*   **`main()`**: A função principal que, atualmente, contém um roteiro detalhado em forma de comentários e pseudocódigo.
*   **Constantes de Configuração**: `RaftDataDir`, `RaftBindAddr`, `HttpBindAddr`, `NodeID` definem configurações estáticas para um nó, indicando a necessidade de um sistema de configuração mais robusto no futuro (e.g., via variáveis de ambiente).
*   **Roteiro de Inicialização (em comentários)**: O arquivo descreve uma sequência clara de 9 passos para iniciar o servidor:
    1.  Inicializar Repositórios e Serviços.
    2.  Inicializar Handlers da API.
    3.  Inicializar a Máquina de Estados Finitos (FSM) do Raft.
    4.  Configurar e iniciar o nó Raft.
    5.  Iniciar a API HTTP interna.
    6.  Inicializar o `RaftCoordinator`.
    7.  Configurar o cliente MQTT.
    8.  Iniciar o serviço de Descoberta.
    9.  Bloquear a execução para manter o servidor ativo.

### Relações e Interdependências Esperadas
*   **Hub Central:** Este arquivo importa quase todos os outros pacotes principais do projeto (`api`, `cluster`, `data`, `domain`, `services`), demonstrando seu papel como o ponto central de integração.
*   **Cadeia de Dependências Clara:** A lógica pretendida segue o padrão: `Repositório` -> `Serviço` -> `Handler` -> `FSM`.
*   **Ponto Crítico de Integração:** A integração entre o `MQTT` (entrada de dados), o `Coordinator` e o `Raft` é o coração do sistema. O `Coordinator` deve receber eventos do MQTT e usar o Raft para replicá-los de forma consistente no cluster.
*   **Cluster Auto-Join:** O `DiscoveryService` deve interagir com o transporte HTTP (`GinHttpTransport`) para permitir que novos nós se juntem ao cluster dinamicamente.

### Observações sobre Incompletude
*   **Totalmente Incompleto:** O arquivo consiste inteiramente em pseudocódigo e comentários. Nenhuma lógica de inicialização é de fato executada.
*   **Blueprint Arquitetural:** Ele funciona como um "TODO" ou um plano de implementação detalhado para o startup do servidor. Os nomes de funções e variáveis no pseudocódigo (ex: `services.NewUserService`, `api.NewEventHandler`, `cluster.NewClusterFSM`) são pistas fortes sobre as assinaturas de construtores e interfaces esperadas nos outros pacotes.

---

## `internal/services/cards.go`

### Função Prevista
Este arquivo implementa o `CardsService`, responsável por toda a lógica de negócios relacionada às cartas dos jogadores. Isso inclui consultar o "álbum" de um usuário, gerenciar a compra de novos pacotes de cartas e orquestrar as operações de troca de cartas entre usuários. Ele serve como a camada de serviço que isola a lógica de domínio das cartas.

### Componentes Notáveis (Existentes ou Implícitos)
*   **`CardsService`**: A struct que contém uma dependência do repositório de cartas (`data.Repository[domain.CardInterface]`), indicando que ele delega as operações de persistência.
*   **`NewCardsService(...)`**: O construtor do serviço, que injeta o repositório como dependência. Retorna uma `CardsServiceInterface`, sugerindo que o contrato da interface é definido em outro lugar (provavelmente `services.go`).
*   **`GetCards(userID string)`**: Um método que busca no repositório todas as cartas associadas a um ID de usuário.
*   **`BuyPack(userID string)`**: Stub para a funcionalidade de comprar um pacote de cartas.
*   **`OfferTrade(...)`**: Stub para a funcionalidade de um usuário oferecer uma carta para troca.
*   **`AcceptTrade(...)`**: Stub para a funcionalidade de um usuário aceitar uma oferta de troca.

### Relações e Interdependências Esperadas
*   **Camada de Dados (`data`):** Depende diretamente de `data.Repository` para abstrair o acesso ao banco de dados, especificamente um repositório que trabalha com a interface `domain.CardInterface`.
*   **Camada de Domínio (`domain`):** Opera sobre a entidade `domain.CardInterface`, mantendo a lógica de negócio desacoplada da implementação concreta do modelo.
*   **Camada de API (`api`/`handlers`):** Espera-se que este serviço seja invocado por handlers de eventos ou de rotas HTTP que processam as ações do usuário, como um `POST /cards/buy` ou um evento MQTT correspondente.
*   **Arquivo `services.go`:** A existência da `CardsServiceInterface` como tipo de retorno do construtor implica que `internal/services/services.go` define as interfaces para todos os serviços, promovendo um contrato unificado.

### Observações sobre Incompletude
*   **Parcialmente Implementado:** O método `GetCards` parece funcional, assumindo que o repositório genérico está implementado.
*   **Lógica de Negócio Ausente:** As funções mais importantes (`BuyPack`, `OfferTrade`, `AcceptTrade`) são apenas stubs. Elas existem na estrutura, mas não contêm nenhuma lógica de implementação, retornando `nil` diretamente. Isso indica que são funcionalidades planejadas, mas ainda não desenvolvidas.
*   **Quebra de Contrato com Domínio (Descoberta na análise de `domain/card.go`):** O método `GetCards` usa `c.GetOwnerID()` para filtrar as cartas. No entanto, o modelo `domain.Card` e sua interface `CardInterface` **não possuem** um campo ou método `GetOwnerID`. Isso é uma inconsistência crítica que impede o `GetCards` de funcionar como pretendido. O modelo de domínio precisa ser estendido para associar uma carta a um dono.

---

## `internal/services/services.go`

### Função Prevista
Este arquivo funciona como um "catálogo de contratos" para a camada de serviço. Sua única responsabilidade é definir as interfaces (`UserServiceInterface`, `CardsServiceInterface`, `MatchServiceInterface`) que descrevem a funcionalidade de negócio disponível na aplicação. Ao centralizar as interfaces, o projeto promove o desacoplamento, permitindo que as camadas superiores (como a API) dependam de abstrações, e não de implementações concretas.

### Componentes Notáveis (Existentes ou Implícitos)
*   **`UserServiceInterface`**: Define o contrato para as operações de usuário: `Register` e `Login`.
*   **`CardsServiceInterface`**: Define o contrato para as operações com cartas: `GetCards`, `BuyPack`, `OfferTrade`, e `AcceptTrade`. Os métodos correspondem exatamente aos implementados (ou em stub) em `cards.go`.
*   **`MatchServiceInterface`**: Define o contrato para o ciclo de vida de uma partida: `StartMatch`, `JoinMatch`, `SurrenderMatch`, `MakeMove`, `CheckOpponentMove`, e `EndMatch`.

### Relações e Interdependências Esperadas
*   **Implementações de Serviço:** Cada interface neste arquivo espera uma implementação concreta em um arquivo correspondente dentro do mesmo pacote (`services/`). Por exemplo, `UserServiceInterface` deve ser implementado por um `UserService` em `users.go`.
*   **Consumidores de Serviço:** Os principais consumidores dessas interfaces serão os handlers da API (`internal/api/handlers.go`) e, crucialmente, o `EventHandler` que é usado pela FSM do Raft. O `main.go` será responsável por injetar as implementações concretas dos serviços nesses consumidores.
*   **Camada de Domínio (`domain`):** Todas as assinaturas de métodos nas interfaces utilizam tipos do pacote `domain` (como `domain.UserInterface`), reforçando que os serviços são os componentes responsáveis por orquestrar a lógica usando os modelos de domínio.

### Observações sobre Incompletude
*   **Completo para seu Propósito:** O arquivo é puramente declarativo (só contém interfaces) e, para essa função, está completo. Ele estabelece um blueprint claro e sólido para as funcionalidades que a camada de serviço deve oferecer.

---

## `internal/services/users.go`

### Função Prevista
Este arquivo fornece a implementação concreta da `UserServiceInterface`, cuidando da lógica de negócio para autenticação de usuários. Suas responsabilidades são o registro de novos usuários no sistema e a validação de credenciais para login.

### Componentes Notáveis (Existentes ou Implícitos)
*   **`UserService`**: A struct que implementa a interface e contém a dependência do repositório de usuários.
*   **`NewUserService(...)`**: O construtor que injeta o repositório de usuários, seguindo o padrão de injeção de dependência.
*   **`Register(username, password string)`**: Implementa o fluxo de registro: gera um novo ID, cria um objeto `domain.User` e o persiste através do repositório.
*   **`Login(username, password string)`**: Implementa o fluxo de login: consulta o repositório por um usuário com o nome e senha correspondentes.

### Relações e Interdependências Esperadas
*   **Contrato de Serviço (`services.go`):** Implementa a interface `UserServiceInterface` ali definida.
*   **Camada de Dados (`data`):** Depende de `data.Repository[domain.UserInterface]` para criar e consultar os dados de usuários.
*   **Camada de Domínio (`domain`):** Instancia `domain.User` e depende de métodos da interface `domain.UserInterface`, como `GetUsername()` e, crucialmente, `CheckPassword()`. Esta dependência mostra uma boa separação de responsabilidades: o serviço não sabe *como* a senha é verificada, apenas delega essa ação ao objeto de domínio.
*   **Camada de API (`api`/`handlers`):** Espera-se que os métodos deste serviço sejam chamados por handlers que expõem endpoints como `/register` e `/login`.

### Observações sobre Incompletude
*   **Funcionalmente Completo (Com Ressalvas):** A lógica de `Register` e `Login` está presente e parece funcional.
*   **Grave Falha de Segurança:** O serviço lida com senhas em texto puro. O método `Register` armazena o `password` diretamente no objeto `User`, que presumivelmente é salvo como está no repositório. Em uma aplicação real, seria mandatório o uso de hashing e salting para as senhas.
*   **Ambiguidade no Login:** A função `Login` retorna `(nil, nil)` quando o usuário não é encontrado. Seria uma prática melhor retornar um erro explícito e distinto, como `ErrInvalidCredentials`, para evitar ambiguidade no código que consome o serviço.

---

## `internal/services/match.go`

### Função Prevista
Este arquivo implementa o `MatchService`, o cérebro por trás da lógica de jogo. Ele gerencia o ciclo de vida completo de uma partida, desde sua criação e a entrada de jogadores, passando pelo processamento de jogadas e desistências, até a finalização e declaração de um vencedor.

### Componentes Notáveis (Existentes ou Implícitos)
*   **`MatchService`**: A struct do serviço, que notavelmente depende de três repositórios diferentes: `matchRepo`, `cardsRepo`, e `usersRepo`. Isso indica que a lógica de uma partida é complexa e precisa acessar e validar dados de múltiplos domínios (partidas, cartas e usuários).
*   **`NewMatchService(...)`**: O construtor que injeta todas as dependências de repositório.
*   **`StartMatch(userID string)`**: Cria uma nova partida com um jogador, definindo o estado como "esperando oponente".
*   **`JoinMatch(userID, matchID string)`**: Permite que um segundo jogador entre em uma partida existente, alterando seu estado para "em andamento".
*   **`SurrenderMatch(userID, matchID string)`**: Finaliza uma partida quando um jogador desiste, declarando o outro como vencedor.
*   **`MakeMove(...)`**: Stub para a lógica de um jogador realizar uma jogada (ex: usar uma carta).
*   **`CheckOpponentMove(...)`**: Stub que sugere um mecanismo para um jogador consultar a última jogada do oponente, característico de um jogo de turnos.
*   **`EndMatch(matchID string)`**: Stub para finalizar a partida, que deve conter a lógica para determinar o vencedor com base nas regras do jogo.

### Relações e Interdependências Esperadas
*   **Contrato de Serviço (`services.go`):** Implementa a interface `MatchServiceInterface`.
*   **Camada de Dados (`data`):** Acoplamento forte com a camada de dados, necessitando de acesso a múltiplos repositórios. A lógica de `MakeMove` (ainda não implementada) certamente precisaria do `cardsRepo` para validar se o jogador possui a carta jogada e do `usersRepo` para buscar informações dos jogadores.
*   **Camada de Domínio (`domain`):** Orquestra objetos `domain.Match`, modificando seu estado através de métodos como `SetStatus` e `SetPlayer2ID`.
*   **Camada de Eventos (`api`/`handlers`):** Este serviço é o alvo principal para eventos de jogo que chegam via MQTT (ex: um evento `make_move` dispararia o método `MakeMove` neste serviço).

### Observações sobre Incompletude
*   **Lógica de Setup Implementada:** As funcionalidades que gerenciam o estado da partida (começar, entrar, desistir) estão parcialmente implementadas e parecem funcionais para alterar o estado básico da partida.
*   **Lógica de Jogo Ausente:** As partes mais críticas e complexas, que definem o jogo em si, são apenas stubs. `MakeMove`, `CheckOpponentMove`, e a lógica de determinação de vencedor em `EndMatch` estão completamente ausentes. O coração da funcionalidade do jogo ainda precisa ser construído.
*   **Vencedor Fixo:** O método `EndMatch` atualmente retorna uma string fixa `"winner_id"`, confirmando que a lógica para calcular o resultado real da partida não foi implementada.
*   **INCOMPATIBILIDADE ARQUITETURAL CRÍTICA (Descoberta na análise de `domain/match.go`):** Este serviço opera sobre um conceito de partida (com `Player1ID`, `Player2ID`, `Status`) que é totalmente diferente e incompatível com o modelo `Match` definido na camada de domínio (que usa uma lista de `Players` e um histórico de `Moves`). O serviço, como está escrito, não pode interagir com o modelo de domínio `domain.Match`, indicando uma severa divergência de design que precisa ser resolvida.

---

## `internal/cluster/http.go`

### Função Prevista
Este arquivo define um transporte HTTP para a comunicação entre os nós do cluster. Sua função é dupla:
1.  **Servidor:** Expor endpoints HTTP (`/raft/join`, `/raft/command`) para que outros nós possam se juntar ao cluster e para que comandos sejam encaminhados ao líder.
2.  **Cliente:** Fornecer métodos (`JoinCluster`, `ForwardCommand`) para que o nó local possa fazer requisições a esses endpoints em outros nós.

Essencialmente, ele implementa a mecânica de encaminhamento de comandos (requests) e de adição de novos membros ao cluster sobre um protocolo conhecido (HTTP), desacoplando essa comunicação da camada de consenso do Raft.

### Componentes Notáveis (Existentes ou Implícitos)
*   **`GinHttpTransport`**: A struct principal que encapsula um servidor Gin (`router`), um cliente Resty (`client`) e uma referência ao nó Raft (`raftNode`).
*   **`NewGinHttpTransport(...)`**: O construtor para criar e inicializar o transporte.
*   **`Start()`**: Método para configurar as rotas e iniciar o servidor HTTP em background.
*   **`JoinCluster(...)`**: Método cliente que envia uma requisição a outro nó para se juntar ao seu cluster.
*   **`ForwardCommand(...)`**: Método cliente que encaminha um evento (comando) para o líder do cluster via HTTP.
*   **`handleJoin(...)`**: O handler HTTP que recebe uma requisição de `JoinCluster` e usa a referência do `raftNode` para adicionar o requisitante como um `Voter`.
*   **`handleCommand(...)`**: O handler HTTP que recebe um comando encaminhado e o aplica ao log do Raft local através de `raftNode.Apply()`.

### Relações e Interdependências Esperadas
*   **Nó Raft (`*raft.Raft`):** É a dependência mais crítica. Toda a lógica dos handlers (`handleJoin`, `handleCommand`) se resume a invocar métodos no objeto `raftNode`.
*   **Coordenador (`coordinator.go`):** O coordenador em um nó seguidor usará o método `ForwardCommand` deste transporte para enviar comandos recebidos para o líder atual do cluster.
*   **Serviço de Descoberta (`discovery.go`):** O serviço de descoberta usará o método `JoinCluster` para tentar se juntar a um cluster existente ao encontrar um novo par na rede.
*   **`cmd/main.go`:** O ponto de entrada da aplicação é responsável por instanciar `GinHttpTransport` e chamar `Start()`.
*   **Bibliotecas Externas:** Dependência explícita do `gin-gonic/gin` para o servidor e `go-resty/resty` para o cliente HTTP.

### Observações sobre Incompletude
*   **Esqueleto Arquitetural:** Assim como `main.go`, este arquivo é um blueprint. Todos os métodos e handlers são compostos apenas por pseudocódigo e comentários que descrevem a lógica pretendida.
*   **Nenhuma Implementação Real:** Nenhum código funcional de servidor ou cliente HTTP está presente. As rotas não são registradas, o servidor não é iniciado e nenhuma requisição é feita. O arquivo serve como um guia detalhado para a implementação futura.

---

## `internal/cluster/transport.go`

### Função Prevista
Este arquivo serve como um "arquivo de cabeçalho" para o transporte de cluster HTTP. Sua função é puramente definicional:
1.  Define os **DTOs (Data Transfer Objects)** `JoinRequest` e `CommandRequest`, que são as estruturas de dados a serem serializadas em JSON para a comunicação de rede.
2.  Define a interface `ClusterTransportInterface`, que estabelece o contrato que qualquer implementação de transporte de cluster (como `GinHttpTransport`) deve seguir.

### Componentes Notáveis (Existentes ou Implícitos)
*   **`JoinRequest`**: Estrutura que carrega o ID e o endereço de um nó que deseja se juntar ao cluster.
*   **`CommandRequest`**: Estrutura que carrega os bytes brutos de um evento a ser aplicado no log do Raft.
*   **`ClusterTransportInterface`**: A interface chave que desacopla o resto do sistema da implementação específica do transporte. Exige os métodos `Start`, `JoinCluster` e `ForwardCommand`.

### Relações e Interdependências Esperadas
*   **`http.go`**: O `GinHttpTransport` em `http.go` é a implementação pretendida desta interface. Os handlers em `http.go` usarão os DTOs aqui definidos para fazer o parse das requisições JSON.
*   **Consumidores da Interface:** Componentes como o `Coordinator` e o `DiscoveryService` devem depender desta interface (`ClusterTransportInterface`), não da implementação concreta (`GinHttpTransport`), para facilitar testes e futuras modificações.

### Observações sobre Incompletude
*   **Completo para seu Propósito:** Como é um arquivo que contém apenas definições de tipos e interfaces, ele está completo. Ele cumpre com sucesso seu papel de definir os contratos e as estruturas de dados para a camada de transporte do cluster. Não há e nem deveria haver lógica executável aqui.

---

## `internal/cluster/coordinator.go`

### Função Prevista
Este arquivo define o `RaftCoordinator`, um dos componentes mais inteligentes e cruciais da arquitetura. Sua função é ser o único ponto de entrada para qualquer comando que precise ser executado de forma consistente em todo o cluster. Ele implementa a lógica central de um sistema baseado em Raft:
1.  Se o nó atual for o **Líder**, ele aplica o comando ao seu próprio log Raft para replicação.
2.  Se o nó atual for um **Seguidor**, ele encaminha o comando para o nó que é o líder atual.
Isso garante que todas as modificações de estado passem obrigatoriamente pelo líder, mantendo a consistência do sistema.

### Componentes Notáveis (Existentes ou Implícitos)
*   **`CoordinatorInterface`**: Define a única e importante função `Handle(event api.Event)`.
*   **`RaftCoordinator`**: A implementação da interface. Contém uma referência ao nó Raft (`raftNode`) para consultar seu estado e aplicar comandos, e ao transporte do cluster (`transport`) para encaminhar comandos.
*   **`NewRaftCoordinator(...)`**: O construtor que injeta suas dependências (`raftNode` e `transport`).
*   **`Handle(event api.Event)`**: O método principal, cujo pseudocódigo descreve claramente o fluxo de decisão: serializar o evento, verificar o estado do nó (`c.raftNode.State()`), e então aplicar localmente (`raftNode.Apply`) ou encaminhar (`transport.ForwardCommand`).

### Relações e Interdependências Esperadas
*   **Entrada de Eventos (e.g., `mqtt.go`):** O principal chamador do `Coordinator.Handle` será o componente que recebe os comandos do mundo exterior. No `main.go`, está implícito que este é o cliente MQTT: ao receber uma mensagem, ele a transforma em um `api.Event` e a entrega ao coordenador.
*   **Nó Raft (`*raft.Raft`):** O coordenador depende vitalmente do nó Raft para verificar o estado (`State()`), descobrir o líder (`Leader()`) e aplicar comandos (`Apply()`).
*   **Transporte do Cluster (`ClusterTransportInterface`):** Para o encaminhamento, o coordenador depende da *interface* do transporte, não de sua implementação. Isso permite que a forma de encaminhamento (HTTP, gRPC, etc.) seja trocada sem impactar o coordenador.
*   **Eventos da API (`api.Event`):** Recebe e processa o tipo `api.Event`, o que o conecta à camada de definição de eventos da aplicação.

### Observações sobre Incompletude
*   **Lógica em Pseudocódigo:** Toda a lógica essencial do método `Handle` está escrita em comentários como pseudocódigo. A função em si apenas retorna `nil`, significando que a peça central do roteamento de comandos ainda não foi implementada.
*   **Blueprint Claro:** O pseudocódigo é um excelente guia para a implementação, cobrindo os cenários de ser líder ou seguidor e a necessidade de serializar os dados. Ele também aponta corretamente a necessidade de tratar o caso de um cluster estar temporariamente sem líder.

---

## `internal/cluster/fsm.go`

### Função Prevista
Este arquivo define a **Máquina de Estados Finitos (FSM)**, a ponte entre o sistema de consenso Raft e a lógica de negócio da aplicação. A responsabilidade da FSM é receber os "logs" que o Raft já confirmou como consensuais entre os nós e "aplicá-los" ao estado local da aplicação. Em outras palavras, ele traduz um comando replicado (como "registrar usuário X") na chamada de serviço real (`userService.Register(...)`). Como cada nó do cluster executa a mesma FSM sobre os mesmos logs na mesma ordem, o estado da aplicação permanece consistente em todo o cluster.

### Componentes Notáveis (Existentes ou Implícitos)
*   **`ClusterFSM`**: A struct que implementa a interface `raft.FSM`. Sua única dependência é a `api.EventHandlerInterface`, mostrando que ela delega a execução da lógica.
*   **`NewClusterFSM(...)`**: O construtor que injeta o `eventHandler`.
*   **`Apply(log *raft.Log)`**: O método mais importante, chamado pela biblioteca Raft para cada log comitado. O pseudocódigo detalha o fluxo:
    1.  Deserializar os dados do log (`log.Data`) em um `api.Event`.
    2.  Usar um `switch` no `event.Method`.
    3.  Chamar o método correspondente no `eventHandler` (ex: `fsm.eventHandler.OnRegister(*event)`).
*   **`Snapshot()`**: Método da interface Raft para criar um "backup" do estado atual da aplicação, permitindo que o Raft descarte logs antigos.
*   **`Restore()`**: Método da interface Raft para restaurar o estado da aplicação a partir de um `Snapshot`, usado para acelerar a sincronização de nós novos ou atrasados.

### Relações e Interdependências Esperadas
*   **Biblioteca Raft:** O `ClusterFSM` é um componente passivo que é acionado diretamente pela instância do `*raft.Raft` configurada no `main.go`.
*   **`api.EventHandlerInterface`:** É a dependência chave. A FSM atua como um despachante (dispatcher), decodificando o evento e passando-o para o `EventHandler`, que de fato contém a lógica para chamar os serviços de domínio corretos.
*   **`api.Event`:** Confirma que o `Event` é a unidade padrão de comando que é serializada, replicada via Raft, e finalmente executada pela FSM.

### Observações sobre Incompletude
*   **Lógica em Pseudocódigo:** Os três métodos da interface (`Apply`, `Snapshot`, `Restore`) são stubs. A lógica principal de despacho no método `Apply` está descrita em comentários, mas não implementada.
*   **Complexidade de Snapshot/Restore:** O pseudocódigo corretamente aponta que a implementação de `Snapshot` e `Restore` é complexa, pois exige a serialização e desserialização de *todo* o estado da aplicação (os dados de todos os repositórios). A sugestão de um `NoOpSnapshot` é uma abordagem comum para adiar essa complexidade.
*   **Nenhum Tratamento de Erro:** A implementação atual apenas retorna `nil`. A lógica real precisaria de um tratamento de erros robusto, especialmente durante a deserialização do evento.

---

## `internal/cluster/discovery.go`

### Função Prevista
Este arquivo implementa um serviço de descoberta de pares na rede local (peer discovery). Sua função é permitir que os nós do servidor se encontrem automaticamente, sem a necessidade de configurar manualmente uma lista de IPs. Ele utiliza uma estratégia simples e eficaz de broadcast UDP:
1.  **Broadcast:** Cada nó envia periodicamente uma mensagem "mágica" para toda a rede local.
2.  **Listen:** Cada nó escuta por esta mesma mensagem.
Ao receber uma mensagem de outro nó, ele sabe que encontrou um "vizinho" e pode tentar se juntar ao seu cluster.

### Componentes Notáveis (Existentes ou Implícitos)
*   **`DiscoveryService`**: A struct principal que contém as configurações de rede (porta, intervalo) e, crucialmente, um campo de callback `OnPeerDiscovered`.
*   **`NewDiscoveryService(...)`**: O construtor do serviço.
*   **`OnPeerDiscovered func(peerAddress string)`**: Um campo de função (callback). Esta é a peça central do design. O serviço de descoberta não sabe *o que fazer* quando encontra um par; ele apenas invoca este callback, passando o endereço do par encontrado. A lógica de ação é definida externamente.
*   **`Start()`**: Inicia o serviço, disparando as rotinas `listen` e `broadcast` em background.
*   **`listen()`**: Uma goroutine que escuta por pacotes UDP. Ao receber a mensagem correta, invoca o callback `OnPeerDiscovered`.
*   **`broadcast()`**: Uma goroutine que, em um loop com `time.Ticker`, envia a mensagem de descoberta para o endereço de broadcast da rede.

### Relações e Interdependências Esperadas
*   **`cmd/main.go`**: O `main` é responsável por criar o `DiscoveryService` e, mais importante, por **implementar o callback `OnPeerDiscovered`**. O pseudocódigo no `main` deixa claro que a ação a ser tomada é chamar o `httpTransport.JoinCluster(...)`, conectando a descoberta à ação de juntar-se ao cluster.
*   **Transporte do Cluster (`http.go`):** O serviço de descoberta em si é desacoplado do transporte, mas eles são ligados pela lógica no `main`. O `DiscoveryService` encontra, e o `http.go` age.
*   **Biblioteca `net`:** O arquivo depende fortemente do pacote `net` do Go para toda a comunicação UDP.

### Observações sobre Incompletude
*   **Lógica em Pseudocódigo:** Os métodos `listen` e `broadcast`, que contêm a lógica de rede, são inteiramente pseudocódigo. Nenhuma conexão UDP real é estabelecida.
*   **Design Baseado em Callback:** O uso do `OnPeerDiscovered` é um bom padrão de design, pois desacopla a detecção da ação, tornando o componente reutilizável e fácil de testar.
*   **Filtro de Auto-descoberta:** O pseudocódigo menciona a necessidade de filtrar mensagens enviadas pelo próprio nó para evitar que ele "se descubra" continuamente, um detalhe importante para a implementação real.

---

## `internal/domain/card.go`

### Função Prevista
Este arquivo define os modelos de domínio (`domain models`) para os conceitos de `Card` (Carta) e `Pack` (Pacote). Ele é o coração do subdomínio do jogo de cartas, estabelecendo não apenas a estrutura de dados (campos como `ID`, `Type`), mas também a lógica de negócio intrínseca a esses objetos. A lógica de "pedra-papel-tesoura" está corretamente encapsulada aqui, no método `Against`, seguindo princípios de Domain-Driven Design (DDD).

### Componentes Notáveis (Existentes ou Implícitos)
*   **`CardInterface`**: Define o contrato que qualquer tipo de carta deve seguir, com métodos para obter seu ID, tipo, e a lógica de confronto (`Against`).
*   **`PackInterface`**: Define o contrato para um agrupamento de cartas, como um deck ou mão, com métodos para adicionar ou remover cartas.
*   **`Card`**: A implementação concreta de uma carta, com `ID` e `Type`.
*   **`Pack`**: A implementação concreta de um pacote de cartas.
*   **`Against(opponent CardInterface)`**: Método no `*Card` que implementa a regra de negócio principal do jogo. Retorna `1` para vitória, `0` para empate e `-1` para derrota.
*   **`DrawCard(index int)`**: Método em `*Pack` que simula a ação de comprar uma carta de um baralho, removendo-a da coleção e retornando-a.

### Relações e Interdependências Esperadas
*   **Camada de Serviço (`services`):** Os serviços, especialmente `MatchService`, serão os principais consumidores da lógica deste arquivo, usando `Against` para determinar os resultados das jogadas. `CardsService` irá gerenciar a posse e criação desses objetos.
*   **Camada de Dados (`data`):** Os repositórios serão responsáveis pela persistência e recuperação de instâncias de `Card` e `Pack`.
*   **Outros Modelos de Domínio:** O `User` provavelmente terá uma coleção de `CardInterface`, e o `Match` envolverá a interação entre cartas de diferentes usuários.

### Observações sobre Incompletude
*   **Funcionalmente Completo:** Para seu escopo, o arquivo parece completo. As interfaces e implementações são coesas e a lógica de negócio está bem encapsulada.
*   **INCONSISTÊNCIA DETECTADA:** O `CardsService` (em `internal/services/cards.go`) possui um método `GetCards` que filtra as cartas usando `c.GetOwnerID()`. No entanto, nem a `CardInterface` nem a struct `Card` neste arquivo possuem um campo `OwnerID` ou um método `GetOwnerID()`. Isso representa uma quebra de contrato entre a camada de serviço e a camada de domínio, indicando que o modelo `Card` está incompleto para satisfazer as necessidades do serviço.

---

## `internal/domain/user.go`

### Função Prevista
Este arquivo define o modelo de domínio `User` (Usuário). Ele representa a entidade `User` no sistema, encapsulando seus dados (ID, nome de usuário, senha) e comportamentos intrínsecos. Ele também define a `UserInterface` para promover o desacoplamento.

### Componentes Notáveis (Existentes ou Implícitos)
*   **`UserInterface`**: Define o contrato público para um usuário, com métodos para `GetID`, `GetUsername` e `CheckPassword`.
*   **`User`**: A implementação concreta do usuário. Contém os campos `ID`, `Username`, `Password` e, notavelmente, um campo `Cards` do tipo `PackInterface`.
*   **`CheckPassword(password string)`**: Método que contém a lógica para verificar uma senha.

### Relações e Interdependências Esperadas
*   **Camada de Serviço (`services`):** O `UserService` é o principal consumidor deste modelo, criando instâncias de `User` no registro e usando `CheckPassword` durante o login.
*   **Camada de Dados (`data`):** O repositório de usuários (`userRepo`) será responsável por persistir e buscar objetos `User`.
*   **Modelo de Domínio `Pack` (`card.go`):** A presença do campo `Cards PackInterface` cria uma associação direta entre um `User` e um `Pack`, indicando que "um usuário tem um pacote de cartas".

### Observações sobre Incompletude
*   **Funcionalmente Completo (com ressalvas):** O arquivo está completo em termos de estrutura. Todos os métodos da interface estão implementados.
*   **FALHA DE SEGURANÇA CRÍTICA:** O método `CheckPassword` implementa a verificação como `u.Password == password`. Isso confirma a suspeita levantada na análise do `UserService`: as senhas estão sendo armazenadas e comparadas em texto plano, o que é uma prática de segurança inaceitável em qualquer aplicação real.
*   **Relação com Cartas:** A associação de um `User` a uma `PackInterface` é uma decisão de design interessante, implicando que a coleção de cartas de um usuário é gerenciada como um único "pacote". A responsabilidade pela inicialização e manipulação desse pacote recairia sobre o `CardsService`.

---

## `internal/domain/match.go`

### Função Prevista
Este arquivo define o modelo de domínio `Match` (Partida), que encapsula o estado e as regras de uma partida em andamento. Sua responsabilidade é gerenciar os jogadores envolvidos, o progresso dos turnos e determinar o vencedor com base nas jogadas. É o modelo de domínio mais complexo do projeto.

### Componentes Notáveis (Existentes ou Implícitos)
*   **`MatchInterface`**: Define o contrato para interagir com uma partida.
*   **`Match`**: A implementação concreta. Contém uma lista de `Players`, um histórico de `Moves` (uma lista de mapas, onde cada mapa é um turno) e um campo `Winner`.
*   **`winnerMutex`**: Um `sync.RWMutex` para proteger o campo `Winner` de acessos concorrentes, uma boa prática.
*   **`AddPlayer`/`RemovePlayer`**: Métodos para gerenciar os participantes.
*   **`MakeMove(...)`**: Contém a lógica principal do jogo. Ele gerencia os turnos, valida as jogadas e, quando um turno está completo (ambos os jogadores jogaram), ele invoca o método `Against` das cartas para determinar o resultado e **imediatamente define o vencedor da partida**.
*   **`GetWinner()`**: Método para ler o vencedor de forma segura.

### Relações e Interdependências Esperadas
*   **Camada de Serviço (`services`):** O `MatchService` deveria ser o principal orquestrador deste objeto.
*   **Camada de Dados (`data`):** O `matchRepo` deveria persistir e carregar instâncias de `Match`.
*   **Outros Modelos de Domínio:** `Match` é um "agregado" que junta os modelos `User` e `Card`, criando uma forte relação entre eles para formar o núcleo do jogo.

### Observações sobre Incompletude
*   **Lógica de Jogo Falha:** Embora a lógica de `MakeMove` esteja implementada, ela possui uma falha fundamental: a partida termina e o vencedor é decidido na **primeira rodada** bem-sucedida. Um jogo real provavelmente teria múltiplas rodadas, um sistema de pontuação e um critério de vitória diferente.
*   **INCOMPATIBILIDADE ARQUITETURAL CRÍTICA:** Existe uma divergência fundamental entre o `domain.Match` e o `services.MatchService`.
    *   O **serviço** (`MatchService`) pensa em uma partida como tendo campos `Player1ID`, `Player2ID` e um `Status` em string (ex: "ongoing", "player1_won").
    *   O **domínio** (`domain.Match`) pensa em uma partida como uma lista de `Players`, um histórico de `Moves` e um campo `Winner`.
    *   Esses dois conceitos são **completamente incompatíveis**. O `MatchService`, como está escrito, não pode usar o `domain.Match`. Isso sugere que as camadas de serviço e domínio foram desenvolvidas de forma independente e com visões conflitantes da mesma entidade. Esta é a inconsistência arquitetural mais severa encontrada no projeto até agora.

---

## `internal/data/repository.go`

### Função Prevista
Este arquivo define uma interface `Repository` genérica, um componente central da camada de acesso a dados. O objetivo é criar um contrato padrão para operações de persistência (CRUD - Create, Read, Update, Delete) que possa ser usado para qualquer entidade de domínio (`User`, `Card`, `Match`). O uso de Generics (`[T any]`) é um padrão moderno em Go que permite reutilizar a mesma interface para diferentes tipos, evitando código duplicado e garantindo uma API de acesso a dados consistente em toda a aplicação.

### Componentes Notáveis (Existentes ou Implícitos)
*   **`Repository[T any]`**: Uma interface genérica onde `T` será substituído por um tipo de domínio (e.g., `domain.UserInterface`).
*   **Métodos CRUD:** `Create`, `Read`, `Update`, `Delete`.
*   **Métodos de Listagem:** `List` para buscar todos os registros, e `ListBy` para uma busca filtrada customizada através de uma função de predicado.

### Relações e Interdependências Esperadas
*   **Implementações Concretas:** Este arquivo define apenas o "o quê", não o "como". O projeto necessita de implementações concretas desta interface. O pseudocódigo no `main.go` sugere a existência de um `data.NewMemoryRepository`, que seria uma implementação em memória (usando mapas) para testes ou prototipagem rápida. Uma implementação para produção poderia ser, por exemplo, um `PostgresRepository` que se comunica com um banco de dados real.
*   **Camada de Serviço (`services`):** Todos os serviços (`UserService`, `CardsService`, `MatchService`) são os consumidores diretos desta interface. Eles dependem do `Repository` para abstrair completamente a lógica de banco de dados, mantendo o código de negócio limpo.

### Observações sobre Incompletude
*   **Apenas a Interface:** O arquivo está completo para seu propósito de definir um contrato.
*   **Nenhuma Implementação Concreta:** A maior lacuna é a ausência total de uma implementação para esta interface no projeto. Sem ela, nenhum dado pode ser salvo ou recuperado, e nenhum dos serviços pode funcionar.
*   **Import Comentado:** A linha `// import "cod-server/internal/utils"` sugere que pode ter havido uma dependência de um pacote de utilitários que foi removido ou ainda não foi criado.

---

## `internal/api/handlers.go`

### Função Prevista
Este arquivo define o contrato `EventHandlerInterface`. Sua função é desacoplar o componente que processa os eventos (o `EventHandler` concreto) daquele que os origina (a `ClusterFSM`). Ele lista todos os tipos de eventos de negócio que o sistema pode manipular, servindo como um "índice" de todas as ações possíveis que podem ser replicadas via Raft.

### Componentes Notáveis (Existentes ou Implícitos)
*   **`EventHandlerInterface`**: Uma interface que define um método `On...` para cada ação de negócio, como `OnRegister`, `OnBuyPack`, `OnStartMatch`, etc. Cada método aceita um `Event` e retorna um `Event`, sugerindo um padrão de requisição/resposta.

### Relações e Interdependências Esperadas
*   **Implementação Concreta (`internal/handlers/handlers.go`):** A struct `EventHandler` localizada no pacote `handlers` é a implementação pretendida desta interface.
*   **Consumidor da Interface (`cluster/fsm.go`):** O principal consumidor é a `ClusterFSM`. Ela terá uma dependência desta interface, permitindo que a FSM chame os métodos `On...` sem conhecer os detalhes da implementação do handler.
*   **Localização do Arquivo:** A presença desta interface no pacote `api` é uma forte indicação de que o arquivo `internal/handlers/handlers.go` está mal localizado e deveria, logicamente, estar dentro do pacote `api` também, como `internal/api/handlers/handlers.go` ou `internal/api/event_handler.go`.

### Observações sobre Incompletude
*   **Completo para seu Propósito:** Sendo um arquivo puramente de definição de interface, ele está completo e cumpre bem seu papel de definir o contrato do event handler.
*   **Import Comentado:** Possui um fragmento de import comentado (`// "cod-server/internal/"`), que é apenas lixo de código.

---

## `internal/api/mqtt/mqtt.go`

### Função Prevista
Este arquivo define um "Adaptador MQTT" (`MQTTAdapter`). Sua responsabilidade é encapsular a comunicação com um broker MQTT, servindo como a principal "porta de entrada" de comandos externos para o servidor. Ele abstrai os detalhes da biblioteca MQTT, fornecendo uma interface simplificada para conectar, publicar e se inscrever em tópicos.

### Componentes Notáveis (Existentes ou Implícitos)
*   **`MQTTAdapterInterface`**: Define as três operações essenciais: `Connect` para iniciar a conexão, `Publish` para enviar um `api.Event`, e `Subscribe` para escutar um tópico com uma função de callback.
*   **`MQTTAdapter`**: A struct que implementará a interface e que conterá a instância do cliente da biblioteca `paho.mqtt.golang`.

### Relações e Interdependências Esperadas
*   **`cmd/main.go`**: O `main` será responsável por instanciar este adaptador, conectar-se ao broker e, crucialmente, usar o método `Subscribe` para escutar o tópico de ações do jogo (ex: "game/actions"). O callback fornecido ao `Subscribe` irá deserializar a mensagem, transformá-la em um `api.Event` e passá-la para o `RaftCoordinator`.
*   **Coordenador do Cluster (`cluster.go`):** O adaptador é o componente que "alimenta" o coordenador com eventos vindos do mundo exterior, iniciando assim toda a cadeia de processamento de um comando.
*   **`api.Event`**: O método `Publish` que aceita um `api.Event` indica que o sistema também pode enviar eventos para o mundo exterior (e.g., para notificar os clientes do jogo sobre mudanças de estado).
*   **Biblioteca Externa:** Dependência direta da biblioteca `github.com/eclipse/paho.mqtt.golang`.

### Observações sobre Incompletude
*   **Apenas Definições:** O arquivo contém apenas a `interface` e a `struct`. Nenhum dos métodos (`Connect`, `Publish`, `Subscribe`) está implementado.
*   **Construtor Ausente:** Falta uma função `NewMQTTAdapter` para instanciar e configurar o cliente MQTT subjacente (definir o endereço do broker, credenciais, etc.), que seria essencial para seu funcionamento.

---

## `internal/api/event.go`

### Função Prevista
Este arquivo define a estrutura `Event`, a unidade fundamental de comunicação para comandos em todo o sistema. Ele serve como um "envelope" padrão para qualquer ação que precise ser processada de forma consistente (via Raft). Sua estrutura permite encapsular qual método invocar e os dados necessários para essa invocação.

### Componentes Notáveis (Existentes ou Implícitos)
*   **`Event`**: A struct principal, contendo:
    *   `Method`: Uma `string` que atua como um identificador de comando (ex: "register", "make_move").
    *   `Timestamp`: Um carimbo de tempo para o evento.
    *   `Payload`: Um `map[string]any` flexível para carregar os dados específicos do comando.
*   **`Json()`**: Um método para serializar o `Event` para JSON.
*   **`FromJson()`**: Uma função auxiliar para deserializar dados JSON de volta para um `Event`.

### Relações e Interdependências Esperadas
*   **DTO Ubíquo:** Esta struct é o Data Transfer Object (DTO) padrão usado em toda a pilha de processamento de comandos:
    *   Criado no `MQTTAdapter` a partir de mensagens externas.
    *   Recebido pelo `Coordinator`, que o serializa para o Raft.
    *   Lido do log do Raft pela `FSM`, que o deserializa.
    *   Passado para os métodos `On...` do `EventHandler`.

### Observações sobre Incompletude
*   **Funcionalmente Completo:** O arquivo está completo para seu propósito. Ele fornece a estrutura de mensagem e os helpers de serialização necessários.
*   **Payload Flexível vs. Seguro:** A escolha de `map[string]any` para o `Payload` oferece grande flexibilidade, mas sacrifica a segurança de tipos em tempo de compilação. Os componentes que consomem o evento precisam fazer asserções de tipo (`.(string)`, `.(float64)`), o que pode causar pânico em tempo de execução se o payload não for o esperado. Uma abordagem alternativa (e mais segura) seria usar structs específicas para o payload de cada método.