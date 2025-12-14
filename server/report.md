# Relatório de Problemas e Inconsistências

Este relatório detalha os problemas, falhas de design e riscos identificados durante a análise arquitetural do projeto. Cada seção aborda um problema específico, seu impacto e uma recomendação para corrigi-lo.

---

### Problema: Incompatibilidade Arquitetural Crítica na Entidade 'Match'

*   **Descrição Detalhada:**
    A camada de serviço (`services`) e a camada de domínio (`domain`) têm duas definições completamente diferentes e incompatíveis para a entidade 'Match' (Partida).
    - O `MatchService` (`internal/services/match.go`) foi escrito para operar em uma struct de partida implícita que possui campos como `Player1ID`, `Player2ID` e `Status` (uma string como "ongoing", "player1_won", etc.). Ele realiza operações como `match.SetStatus("ongoing")`.
    - O `domain.Match` (`internal/domain/match.go`), por outro lado, modela uma partida com uma lista de `Players` (`[]UserInterface`), um histórico de `Moves` (`[]map[string]CardInterface`) e um campo `Winner`. A lógica de negócio para determinar o vencedor está contida dentro do próprio modelo de domínio.
    Essa discrepância significa que as duas camadas não podem se comunicar. O `MatchService`, como está escrito, não pode usar e orquestrar o `domain.Match`, que deveria ser seu principal objetivo.

*   **Localização(ões) no Código:**
    - **Visão do Serviço:** `internal/services/match.go` (especificamente nos métodos `StartMatch`, `JoinMatch`, `SurrenderMatch`).
    - **Visão do Domínio:** `internal/domain/match.go` (a struct `Match` e seu método `MakeMove`).

*   **Impacto e Risco:**
    **Gravíssimo.** Este é o problema arquitetural mais severo do projeto. Ele quebra a espinha dorsal do design em camadas, tornando a lógica de negócio para partidas completamente não funcional. A camada de serviço não consegue usar o objeto de domínio que deveria manipular, resultando em um beco sem saída para a implementação. Qualquer tentativa de fazer o `MatchService` funcionar falhará imediatamente.

*   **Cenário de Manifestação:**
    O problema se manifestaria no momento em que um desenvolvedor tentasse implementar a lógica real no `MatchService` (que hoje é parcialmente stub/implementado de forma divergente). Por exemplo, ao tentar buscar uma partida do repositório (`matchRepo.Get(matchID)`), o serviço receberia um objeto `domain.Match`, mas tentaria chamar métodos inexistentes como `SetStatus()` ou acessar campos como `Player1ID`, causando um erro de compilação imediato.

*   **Recomendação Preliminar:**
    Uma das duas visões precisa ser descartada, e a outra, adotada por ambas as camadas. A abordagem do `domain.Match` (com histórico de jogadas e lógica encapsulada) é mais rica e alinhada com os princípios do Domain-Driven Design. A recomendação é **refatorar completamente o `internal/services/match.go`** para que ele opere sobre o modelo `domain.Match`. O serviço deve parar de gerenciar campos como `Status` e, em vez disso, chamar os métodos do objeto de domínio (`match.AddPlayer`, `match.MakeMove`, etc.) para manipular o estado da partida.

---

### Problema: Falha de Segurança Crítica no Manuseio de Senhas

*   **Descrição Detalhada:**
    O sistema armazena, transmite e compara senhas de usuário em texto plano, sem qualquer forma de hashing ou salting. O `UserService` recebe a senha como uma string, a coloca diretamente no objeto `domain.User`, e o método `CheckPassword` no `domain.User` compara a senha recebida com a senha armazenada usando uma simples igualdade de strings (`u.Password == password`).

*   **Localização(ões) no Código:**
    - **Armazenamento:** `internal/services/users.go`, no método `Register`.
    - **Verificação:** `internal/domain/user.go`, no método `CheckPassword`.

*   **Impacto e Risco:**
    **Crítico.** Este é o tipo de falha que torna uma aplicação inviável para produção. Se a camada de persistência (o banco de dados ou mesmo um arquivo de snapshot) for comprometida, todas as senhas de todos os usuários serão expostas imediatamente. Isso representa um risco de segurança massivo para os usuários e para a reputação do sistema.

*   **Cenário de Manifestação:**
    O problema existe permanentemente no estado atual do código. Ele se torna uma catástrofe no momento em que qualquer forma de persistência de dados é implementada e o primeiro usuário se registra. Qualquer pessoa com acesso aos dados armazenados (desenvolvedor, administrador de banco de dados, ou um invasor) pode ler as senhas.

*   **Recomendação Preliminar:**
    A prática padrão da indústria deve ser implementada.
    1.  **Hashing na Criação:** No `UserService.Register`, antes de salvar o usuário, a senha deve ser processada por um algoritmo de hashing de via única, forte e lento, como **bcrypt** ou **Argon2**. O que é salvo no banco de dados é o hash, não a senha.
    2.  **Comparação com Hash:** O método `domain.User.CheckPassword` deve ser modificado para aceitar a senha em texto plano e o hash armazenado, e usar a função de comparação da biblioteca de hashing (ex: `bcrypt.CompareHashAndPassword`) para verificar a correspondência.

---

### Problema: Contrato Quebrado entre `CardsService` e `domain.Card`

*   **Descrição Detalhada:**
    O `CardsService` possui um método `GetCards` que tenta filtrar as cartas de um usuário específico. A lógica de filtro implementada é `c.GetOwnerID() == userID`. No entanto, ao inspecionar a camada de domínio, a `CardInterface` e sua implementação `Card` não possuem um campo `OwnerID` nem um método `GetOwnerID`.

*   **Localização(ões) no Código:**
    - **Consumidor da Informação:** `internal/services/cards.go`, no método `GetCards`.
    - **Fonte Ausente da Informação:** `internal/domain/card.go`, na definição da `CardInterface` e da `struct Card`.

*   **Impacto e Risco:**
    **Alto.** O código não compila. A funcionalidade principal do `CardsService` (obter as cartas de um usuário) é impossível de ser executada, pois depende de um método que não existe. Isso demonstra uma falta de sincronia entre o design da camada de serviço e o design da camada de domínio.

*   **Cenário de Manifestação:**
    Este problema é detectado em tempo de compilação. Assim que um desenvolvedor tenta compilar o projeto, o compilador Go apontará que `c.GetOwnerID` não é um método do tipo `domain.CardInterface`.

*   **Recomendação Preliminar:**
    O modelo de domínio `Card` precisa ser a fonte da verdade. A recomendação é estender o domínio para incluir a posse da carta.
    1.  Adicionar um campo `OwnerID string` à struct `domain.Card`.
    2.  Adicionar um método `GetOwnerID() string` à `CardInterface` e implementá-lo no `*Card`.
    Desta forma, o modelo de domínio passa a conter a informação que a camada de serviço corretamente identificou como necessária.

---

### Problema: Ausência de Implementações de Componentes Essenciais

*   **Descrição Detalhada:**
    O projeto, embora tenha uma arquitetura bem delineada por interfaces e pseudocódigo, é em grande parte não funcional devido à ausência de implementações concretas para vários componentes cruciais. A aplicação não pode ser construída e executada em seu estado atual. As principais ausências são:
    - **Repositório de Dados:** Não existe uma implementação da interface `data.Repository`. Sem isso, nenhuma entidade (`User`, `Card`, `Match`) pode ser salva ou recuperada.
    - **Orquestração Principal:** O `cmd/main.go` é inteiramente pseudocódigo, então nenhum serviço ou servidor é de fato iniciado.
    - **Lógica de Cluster:** Todo o pacote `cluster` (`coordinator.go`, `http.go`, `fsm.go`, `discovery.go`) contém apenas pseudocódigo em seus métodos principais. A lógica de consenso, encaminhamento de comandos e descoberta de nós não está implementada.
    - **Manipulação de Eventos:** Os métodos `On...` do `EventHandler` são stubs vazios, significando que, mesmo que um evento chegasse à FSM, nenhuma ação de negócio seria executada.
    - **Entrada de Dados:** O `MQTTAdapter` não tem implementação, impedindo a conexão com um broker e o recebimento de comandos externos.

*   **Localização(ões) no Código:**
    - `internal/data/` (falta de `memory_repository.go` ou similar)
    - `cmd/main.go`
    - `internal/cluster/*.go`
    - `internal/handlers/handlers.go`
    - `internal/api/mqtt/mqtt.go`

*   **Impacto e Risco:**
    **Total.** O sistema como um todo não funciona. É um esqueleto arquitetural, não um software funcional. O risco é que a implementação de cada um desses componentes pode revelar novos problemas ou desafios não previstos no design inicial.

*   **Cenário de Manifestação:**
    O problema se manifesta ao tentar compilar e executar o projeto (`go run cmd/main.go`), o que não resultaria em um servidor funcional, pois a função `main` está vazia. Mesmo que estivesse preenchida, falharia ao tentar instanciar componentes cujas implementações não existem (como o repositório).

*   **Recomendação Preliminar:**
    A implementação deve seguir uma ordem lógica de dependências:
    1.  Implementar um `MemoryRepository` que satisfaça a `data.Repository[T]` para permitir que os serviços operem com dados em memória.
    2.  Implementar os métodos `On...` no `EventHandler` para que ele chame corretamente os métodos dos serviços.
    3.  Implementar a lógica no `main.go` para de fato instanciar e injetar todas as dependências (Repository -> Services -> EventHandler -> FSM).
    4.  Implementar a lógica de `cluster` (FSM, Coordinator, etc.) para habilitar a replicação de estado.
    5.  Implementar o `MQTTAdapter` para conectar o sistema ao mundo exterior.

---

### Problema: Lógica de Jogo Simplista e Falha

*   **Descrição Detalhada:**
    O `domain.Match`, embora contenha a implementação da lógica de jogo, está falho. A função `MakeMove` determina o vencedor da partida inteira com base no resultado de uma única rodada. Assim que um jogador vence uma troca de cartas, ele é declarado o `Winner` da partida, e o jogo efetivamente termina.

*   **Localização(ões) no Código:**
    - `internal/domain/match.go`, no método `MakeMove`.

*   **Impacto e Risco:**
    **Médio a Alto.** Embora não quebre a compilação, torna a experiência de jogo incorreta e trivial. Se esta lógica fosse para produção, o jogo não seria jogável como pretendido (assumindo-se que um jogo de cartas deve ter múltiplas rodadas ou um sistema de pontuação). O risco é entregar uma funcionalidade que não atende aos requisitos básicos do negócio.

*   **Cenário de Manifestação:**
    O problema ocorreria assim que dois jogadores completassem o primeiro turno de uma partida. O vencedor daquele turno seria declarado vencedor da partida inteira.

*   **Recomendação Preliminar:**
    A lógica do `domain.Match` precisa ser expandida para suportar um jogo completo.
    1.  Introduzir um sistema de pontuação ou contagem de vitórias por rodada nos campos da struct `Match`.
    2.  Modificar `MakeMove` para, em vez de definir o `Winner` da partida, apenas registrar o vencedor da *rodada* e atualizar a pontuação.
    3.  Adicionar uma nova lógica (talvez em um método `CheckEndCondition` ou similar) que verifique se a condição de vitória da partida foi atingida (ex: melhor de 3, ou um jogador atingiu X pontos).

---

### Problema: Organização Estrutural de Pacotes Inconsistente

*   **Descrição Detalhada:**
    Existe uma inconsistência na organização dos pacotes relacionados à API e aos handlers. A interface `EventHandlerInterface` está corretamente localizada em `internal/api/handlers.go`, definindo um contrato público do pacote `api`. No entanto, sua implementação, `EventHandler`, reside em um pacote separado e de nível superior, `internal/handlers/`.

*   **Localização(ões) no Código:**
    - **Interface:** `internal/api/handlers.go`
    - **Implementação:** `internal/handlers/handlers.go`

*   **Impacto e Risco:**
    **Baixo.** Este é um problema "cosmético" ou de organização. Não impede o funcionamento do código, mas vai contra as convenções e a clareza da arquitetura. Um desenvolvedor que olha para o pacote `api` espera encontrar tanto a interface quanto suas implementações diretas, ou pelo menos tê-las em um subpacote (como `api/handlers`). A separação atual é confusa.

*   **Cenário de Manifestação:**
    O problema é uma fonte constante de pequena fricção e confusão para desenvolvedores que navegam na base de código, tornando mais difícil entender a relação entre os componentes.

*   **Recomendação Preliminar:**
    Consolidar os pacotes para refletir melhor a relação entre interface e implementação. A melhor abordagem seria mover `internal/handlers/handlers.go` para `internal/api/handlers.go` e renomear o arquivo `internal/api/handlers.go` (que contém a interface) para `internal/api/interfaces.go` ou similar, de modo que o pacote `api` contenha todos os seus componentes relacionados. Ou, alternativamente, mover o `handlers.go` para `internal/api/` diretamente.
