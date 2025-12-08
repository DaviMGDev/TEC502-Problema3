# Descrição do Jogo: "Cards of Destiny" (CoD)

Este documento descreve as regras, mecânicas e o fluxo do jogo "Cards of Destiny", com base na análise da arquitetura do servidor e em esclarecimentos adicionais.

## 1. Conceito Principal

"Cards of Destiny" é um jogo de cartas colecionáveis para dois jogadores (1v1), onde o combate é resolvido com a mecânica clássica de Pedra-Papel-Tesoura. O jogo inclui um sistema de progressão e de economia, onde os jogadores podem colecionar, evoluir e trocar cartas para montar a sua coleção.

## 2. As Entidades do Jogo

O jogo é composto por quatro entidades principais que interagem entre si.

-   **Utilizador (User):** Representa um jogador no sistema. A entidade `User` armazena as informações de login e, mais importante, a sua **coleção de cartas pessoal**.

-   **Carta (Card):** É a unidade fundamental do jogo. Cada carta é definida por:
    -   **Tipo (Type):** Uma de três opções: `Pedra`, `Papel` ou `Tesoura`.
    -   **Nível (Level):** Um valor numérico que representa a força da carta. O nível é usado como **critério de desempate** em combates entre cartas do mesmo tipo.

-   **Pacote (Package):** Um contentor que agrupa um conjunto de cartas. Os pacotes são a principal forma de os jogadores obterem novas cartas.

-   **Partida (Match):** Regista um confronto 1v1 entre dois jogadores. Armazena os IDs dos jogadores envolvidos, as cartas que cada um jogou e quem foi o vencedor da partida.

## 3. O Ciclo de Vida do Jogador e a Economia

O jogo possui um ciclo de vida e uma economia focados na aquisição e troca de cartas.

1.  **Registo e Login:** Um novo jogador cria uma conta no sistema para ter acesso ao jogo.

2.  **Aquisição de Cartas:** A principal forma de obter cartas é através da compra de pacotes.
    -   **Pacote Inicial:** Presume-se que um novo jogador receba um pacote de cartas inicial.
    -   **Compra de Pacotes:** Um jogador pode solicitar um novo pacote de cartas a qualquer momento. Esta ação é **gratuita** e não consome recursos do jogador, apenas o stock de pacotes disponíveis no servidor. Os pacotes contêm cartas com níveis gerados aleatoriamente.

3.  **Gestão da Coleção:** Um jogador pode consultar a qualquer momento as cartas que possui na sua coleção pessoal.

4.  **Trocas (Swapping):** Os jogadores podem trocar cartas entre si para otimizar os níveis da sua coleção. A regra para a troca é estrita:
    -   A troca só é permitida entre **cartas do mesmo tipo** (ex: um jogador pode trocar a sua `Pedra Nível 5` pela `Pedra Nível 8` de outro jogador, mas não por uma `Tesoura Nível 8`).

## 4. Fluxo de uma Partida (Gameplay)

Uma partida é um evento rápido e direto, consistindo de uma única jogada para determinar o vencedor.

*Nota: A implementação da lógica do serviço de jogo (`GameService`) está atualmente em falta no código, pelo que este fluxo descreve o comportamento pretendido com base no design e nos esclarecimentos fornecidos.*

1.  **Matchmaking (Procura de Partida):**
    -   O matchmaking é gerido por uma **fila interna** no servidor.
    -   Para entrar na fila, um jogador envia um pedido para iniciar um jogo (o evento `StartGame`).
    -   Se a fila estiver vazia, o jogador é adicionado e fica a aguardar um oponente.
    -   Quando um segundo jogador entra na fila, o servidor automaticamente forma uma parida entre os dois. Não existe um sistema de pareamento por nível ou ranking; a regra é **"o primeiro a entrar joga com o segundo"**.
    -   O servidor então cria o objeto `Match` e notifica ambos os jogadores que a partida começou.

2.  **Fazer uma Jogada (`play_card`):**
    -   A revelação das jogadas é **simultânea**. O servidor aguarda que ambos os jogadores enviem a sua jogada. A ordem em que os eventos `MakeMove` chegam não influencia o resultado.
    -   Cada jogador escolhe uma carta da sua coleção e envia o evento de jogada.

3.  **Resolução do Confronto:**
    -   Assim que o servidor recebe a segunda jogada, a partida é imediatamente resolvida:
        -   **Regra Principal:** A vitória é decidida pela regra clássica do Pedra-Papel-Tesoura.
        -   **Regra de Desempate:** Se ambos os jogadores jogarem cartas do mesmo tipo (ex: Pedra vs. Pedra), o vencedor é o jogador cuja carta tiver o **Nível (Level) mais alto**. Se os níveis também forem iguais, a partida termina em **empate**.

4.  **Consequências e Fim da Partida:**
    -   O resultado (vitória, derrota ou empate) é enviado para ambos os jogadores.
    -   **Atualmente, o resultado não é guardado** para um histórico de longo prazo.
    -   **As cartas usadas na partida não são consumidas nem sofrem qualquer tipo de penalidade ou cooldown.** Elas retornam à coleção do jogador e podem ser usadas na partida seguinte.
    -   Um jogador pode também render-se (`player_surrender`), concedendo a vitória imediata ao oponente.

## 5. Modelo de Interação

Toda a interação entre o cliente de um jogador e o servidor é orientada a eventos e utiliza o protocolo MQTT.

-   **Comunicação:** O cliente publica eventos para tópicos MQTT específicos no servidor, dependendo da ação que deseja executar.
    -   Exemplos de tópicos de pedido: `cod/request/login`, `cod/request/start_match`.
-   **Formato da Mensagem:** Todas as mensagens seguem a estrutura padronizada `protocol.Event` (em JSON), que contém o método a ser invocado e os dados da operação.
-   **Fluxo de Resposta:** Após processar um pedido, o servidor publica um `protocol.Event` de resposta num tópico específico, informando o cliente do resultado da sua ação.
