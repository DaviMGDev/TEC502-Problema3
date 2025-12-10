# Metodologia para Criação de Testes Unitários e de Integração

Esta metodologia detalha as etapas que seriam seguidas para desenvolver uma suíte robusta de testes para o projeto do servidor do jogo "Cards of Destiny", abrangendo tanto testes unitários quanto de integração.

## 1. Análise Preliminar e Definição da Estratégia

Antes de começar a escrever qualquer teste, é crucial estabelecer uma base sólida de compreensão e planeamento:

*   **Identificação de Ferramentas de Teste:**
    *   Verificação de ficheiros `_test.go` existentes no projeto, o que indicaria padrões de teste já estabelecidos.
    *   Análise do ficheiro `go.mod` para identificar o uso de bibliotecas de teste externas (ex: `testify/assert`, `testify/mock`). Na ausência destas, o foco seria o pacote `testing` padrão do Go.
*   **Definição do Escopo por Tipo de Teste:**
    *   **Testes Unitários:** O objetivo é validar a lógica de componentes individuais (funções, métodos) em completo isolamento das suas dependências.
        *   **Alvos:** Pacotes como `domain` (ex: `Card.Against`), `utils` (ex: `ParseCommand`, `SafeMap`), e os métodos de cada `service` (ex: `userService.Register`).
        *   **Técnica:** Utilização extensiva de "mocks" (objetos que simulam o comportamento de dependências reais, como `data.Repository`) para garantir que o teste se concentra exclusivamente na unidade sob validação, tornando-o rápido e focado.
    *   **Testes de Integração:** O objetivo é verificar a interação e colaboração entre múltiplos componentes ou camadas do sistema.
        *   **Alvos:** Fluxos que envolvem a comunicação entre `handlers`, `services` e `data` (ex: um pedido completo de registo de utilizador).
        *   **Técnica:** Utilização das implementações reais dos componentes que estão a ser integrados. Componentes externos (como o broker MQTT) podem ser "mockados" ou simulados para controlar o ambiente de teste.

## 2. Priorização e Planeamento da Implementação

A implementação dos testes seria faseada para maximizar a eficiência e a cobertura:

1.  **Prioridade Alta: Testes Unitários para Lógica Base e Funcional.**
    *   **Pacote `domain`:** Cobrir exaustivamente a lógica de `Card.Against` com todas as combinações e o desempate por nível.
    *   **Pacote `utils`:** Testar a robustez das funções utilitárias como `ParseCommand`, `SafeMap` e `SafeList`.
    *   **Pacotes `services` (Funcionais):** Desenvolver testes unitários completos para `UserService` e `CardService`, utilizando mocks para o repositório.
2.  **Prioridade Média: Testes para Lógica Incompleta (Abordagem TDD).**
    *   Para os componentes ainda não implementados (ex: `GameService`, `HandlersImplementation`), eu escreveria os ficheiros de teste primeiro. Os testes iriam falhar, claro, mas eles serviriam como uma **especificação viva** do que o código precisa de fazer. À medida que a lógica fosse implementada, os testes seriam atualizados para passar. Esta é uma forma eficaz de aplicar Test-Driven Development (TDD).
3.  **Prioridade Baixa: Testes de Integração.**
    *   Uma vez que os testes unitários fornecessem uma rede de segurança para os componentes individuais, seriam desenvolvidos testes de integração para os fluxos essenciais (ex: registo de utilizador do `handler` ao `repositório`).

## 3. Execução e Validação

O processo de execução seria iterativo:

1.  **Ficheiros de Teste:** Para cada pacote ou componente a ser testado, seriam criados ficheiros com o sufixo `_test.go` (ex: `user_service_test.go`).
2.  **Casos de Teste:** Dentro destes ficheiros, seriam escritos casos de teste para cenários de sucesso, falha, condições de limite e tratamento de erros.
3.  **Execução:** Os testes seriam executados utilizando o comando padrão do Go: `go test ./...` (ou `go test ./internal/services/...` para testes mais direcionados).
4.  **Refatoração:** Os testes seriam mantidos e refatorados à medida que o código-fonte evoluísse, garantindo que a suíte de testes permanecesse relevante e eficaz.
