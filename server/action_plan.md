# Plano de Ação para Correção do Projeto

Este plano descreve os passos para resolver os problemas identificados no `report.md`, em ordem de prioridade.
**Status Geral do Plano:** [ ] Pendente

---

## Prioridade 1: Correções Estruturais e de Contrato

*Objetivo: Resolver as inconsistências fundamentais entre as camadas para que o projeto tenha uma base coesa e compilável.*

1.  **Refatoração do Modelo de Domínio `Card` (Correção de Contrato Quebrado)**
    *   **Descrição:** Adicionar a noção de posse ao modelo `Card` para satisfazer a necessidade do `CardsService`.
    *   **Arquivo(s) a modificar:** `internal/domain/card.go`
    *   **Progresso:**
        - [ ] Adicionar campo `OwnerID string` à struct `Card`.
        - [ ] Adicionar método `GetOwnerID() string` à interface `CardInterface`.
        - [ ] Implementar o método `GetOwnerID()` na struct `*Card` para retornar o novo campo.

2.  **Harmonização do `MatchService` com `domain.Match` (Correção de Incompatibilidade Crítica)**
    *   **Descrição:** Refatorar o `MatchService` para que ele opere exclusivamente sobre o modelo `domain.Match`, abandonando sua implementação conflitante e passando a usar o modelo de domínio como a única fonte da verdade.
    *   **Arquivo(s) a modificar:** `internal/services/match.go`
    *   **Progresso:**
        - [ ] Reescrever `StartMatch` para criar uma instância de `domain.Match`, adicionar o primeiro jogador via `AddPlayer`, e salvá-la no repositório.
        - [ ] Reescrever `JoinMatch` para buscar um `domain.Match` do repositório e usar o método `AddPlayer`.
        - [ ] Reescrever `SurrenderMatch` para buscar um `domain.Match`, determinar o vencedor com base em quem desistiu, e atualizar o estado da partida.
        - [ ] Implementar `MakeMove` para buscar o `domain.Match` do repositório e invocar seu método `MakeMove`.
        - [ ] Remover os métodos/lógica de `CheckOpponentMove` e `EndMatch` do serviço, pois essa lógica deve residir no `domain.Match`.

3.  **Consolidação da Estrutura de Pacotes**
    *   **Descrição:** Mover o arquivo de implementação do `EventHandler` para o pacote `api`, onde sua interface já reside, para melhorar a coesão e clareza da estrutura do projeto.
    *   **Arquivo(s) a modificar:** `internal/handlers/handlers.go` (mover), `internal/api/handlers.go` (renomear).
    *   **Progresso:**
        - [ ] Mover o arquivo `internal/handlers/handlers.go` para `internal/api/event_handler.go`.
        - [ ] Renomear o arquivo `internal/api/handlers.go` (que contém a interface) para `internal/api/interfaces.go` para refletir melhor seu conteúdo.
        - [ ] Atualizar todos os `import` no projeto que foram quebrados por esta mudança (principalmente em `cmd/main.go` e `internal/cluster/fsm.go`).

---

## Prioridade 2: Correção da Falha de Segurança

*Objetivo: Eliminar o armazenamento de senhas em texto plano.*

1.  **Implementação de Hashing de Senhas**
    *   **Descrição:** Substituir o armazenamento e comparação de senhas em string pura por um mecanismo de hashing seguro como o `bcrypt`.
    *   **Arquivo(s) a modificar:** `internal/services/users.go`, `internal/domain/user.go`.
    *   **Progresso:**
        - [ ] Adicionar a biblioteca `golang.org/x/crypto/bcrypt` às dependências do projeto.
        - [ ] Em `UserService.Register`, antes de criar o `domain.User`, gerar um hash da senha recebida usando `bcrypt.GenerateFromPassword`.
        - [ ] Armazenar o **hash** da senha (como string) no campo `Password` do `domain.User`, não a senha original.
        - [ ] Modificar o método `domain.User.CheckPassword` para receber a senha em texto plano e comparar com o hash armazenado usando `bcrypt.CompareHashAndPassword`.

---

## Prioridade 3: Implementação da Funcionalidade Central (Stubs)

*Objetivo: Transformar o esqueleto arquitetural em uma aplicação funcional, implementando a lógica que hoje é apenas pseudocódigo ou stubs.*

1.  **Criação do Repositório em Memória**
    *   **Descrição:** Criar uma implementação funcional da interface `data.Repository` que armazena os dados em mapas na memória, permitindo que a aplicação seja executada para testes e desenvolvimento.
    *   **Arquivo(s) a modificar:** `internal/data/memory_repository.go` (novo arquivo).
    *   **Progresso:**
        - [ ] Criar o arquivo `internal/data/memory_repository.go`.
        - [ ] Definir a struct `MemoryRepository[T]` contendo um `map[string]T` e um `sync.RWMutex`.
        - [ ] Criar a função construtora `NewMemoryRepository`.
        - [ ] Implementar o método `Create`.
        - [ ] Implementar o método `Read`.
        - [ ] Implementar o método `Update`.
        - [ ] Implementar o método `Delete`.
        - [ ] Implementar o método `List`.
        - [ ] Implementar o método `ListBy`.

2.  **Implementação do Ponto de Entrada e Orquestração**
    *   **Descrição:** Substituir o pseudocódigo em `cmd/main.go` por código Go real que inicializa e conecta todos os componentes do sistema.
    *   **Arquivo(s) a modificar:** `cmd/main.go`.
    *   **Progresso:**
        - [ ] Instanciar os `MemoryRepository` para cada entidade.
        - [ ] Instanciar os serviços (`UserService`, etc.) injetando os repositórios.
        - [ ] Instanciar o `EventHandler` injetando os serviços.
        - [ ] Instanciar e configurar o `ClusterFSM` injetando o `EventHandler`.
        - [ ] Implementar a configuração e inicialização do nó Raft.
        - [ ] Implementar a inicialização do `GinHttpTransport` e do `RaftCoordinator`.
        - [ ] Implementar a inicialização do `MQTTAdapter` e a subscrição no tópico de eventos.
        - [ ] Implementar a inicialização do `DiscoveryService`.

3.  **Implementação da Lógica de Cluster e Handlers**
    *   **Descrição:** Substituir o pseudocódigo restante por código funcional.
    *   **Arquivo(s) a modificar:** `internal/cluster/*.go`, `internal/api/event_handler.go`.
    *   **Progresso:**
        - [ ] Implementar o `switch` de despacho de eventos no `ClusterFSM.Apply`.
        - [ ] Implementar a lógica "líder/seguidor" no `RaftCoordinator.Handle`.
        - [ ] Implementar os handlers Gin e o cliente Resty no `GinHttpTransport`.
        - [ ] Implementar a lógica de broadcast e escuta UDP no `DiscoveryService`.
        - [ ] Implementar todos os métodos `On...` no `EventHandler` para chamar os serviços apropriados.

---

## Prioridade 4: Refinamento da Lógica de Jogo

*Objetivo: Corrigir a lógica de jogo falha para permitir partidas com múltiplas rodadas.*

1.  **Melhoria do `domain.Match`**
    *   **Descrição:** Modificar o modelo `Match` para suportar um sistema de pontuação e múltiplas rodadas, em vez de terminar na primeira jogada.
    *   **Arquivo(s) a modificar:** `internal/domain/match.go`.
    *   **Progresso:**
        - [ ] Adicionar campos de pontuação à struct `Match` (ex: `Scores map[string]int`).
        - [ ] Modificar `MakeMove` para, em vez de definir um vencedor final, determinar o vencedor da *rodada* e atualizar a pontuação.
        - [ ] Adicionar uma lógica para verificar a condição de fim de jogo (ex: melhor de 3) antes de definir o `Winner`.
        - [ ] Modificar `GetWinner()` para apenas retornar um vencedor final se a condição de fim de jogo for atingida.
