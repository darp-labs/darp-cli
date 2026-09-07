# DARP CLI — Roadmap

> **Status:** Documento de direção do projeto
> **Última atualização:** 2026-08-08
> **Objetivo:** Registrar as decisões arquiteturais e a evolução planejada do DARP CLI, preservando o contexto entre diferentes sessões de desenvolvimento.

---

## 1. Visão do projeto

O **DARP CLI** tem como objetivo evoluir de uma CLI de inicialização e diagnóstico de projetos para uma plataforma capaz de preparar, estruturar e manter repositórios **agent-ready** — ambientes preparados para desenvolvimento assistido por agentes de IA.

O DARP deve ser **agnóstico ao agente, à linguagem e à stack**, sempre que possível.

A visão de longo prazo é:

```text
                         DARP Ecosystem
                               │
              ┌────────────────┴────────────────┐
              │                                 │
          DARP CLI                         DARP Registry
              │                                 │
              ▼                                 ▼
        DARP Harness                      Assets / Catalog
              │                                 │
              ├───────────────┐         ┌───────┼────────┐
              │               │         │       │        │
          Knowledge       Workflow    Skills  Agents    MCP
              │               │
              └───────────────┴───────────────┐
                                              │
                                              ▼
                                         Repository
                                              │
                    ┌───────────────┬─────────┼─────────┐
                    ▼               ▼         ▼         ▼
                  Codex          Claude    Copilot    Outros
```

A ideia central é que o **harness pertence ao projeto**, e não a um agente específico.

---

# 2. Princípio central

O DARP deve aplicar os princípios de **Harness Engineering** ao desenvolvimento de software com agentes.

O objetivo não é simplesmente gerar prompts ou arquivos de instruções.

O objetivo é criar um ambiente no qual agentes possam:

* entender o projeto;
* descobrir sua arquitetura;
* consultar conhecimento persistente;
* executar tarefas;
* verificar seu próprio trabalho;
* manter estado entre sessões;
* seguir convenções e políticas;
* utilizar ferramentas disponíveis;
* evoluir o projeto de forma incremental;
* trabalhar com diferentes agentes sem depender de um fornecedor específico.

### Princípio

> **O DARP prepara o ambiente. O agente executa o trabalho.**

Ou, em termos da filosofia de Harness Engineering:

> **Humans steer. Agents execute.**

---

# 3. Estado atual do DARP CLI

Os seguintes comandos já foram implementados:

```text
darp init
darp version
darp help
darp doctor
```

Esses comandos constituem a base inicial da CLI.

O `darp init` é responsável pela inicialização do projeto.

O `darp doctor` é responsável pelo diagnóstico e deverá, futuramente, também ser capaz de avaliar a qualidade do Harness.

---

# 4. Próxima grande capacidade: DARP Harness

A próxima evolução planejada é a criação do conceito de:

## DARP Harness

O Harness será uma estrutura padronizada que prepara um repositório para desenvolvimento assistido por agentes.

A primeira funcionalidade relacionada será:

```bash
darp harness init
```

Essa funcionalidade deverá:

1. analisar o repositório;
2. identificar sua estrutura;
3. detectar tecnologias e ferramentas quando possível;
4. identificar documentação existente;
5. identificar mecanismos de build e testes;
6. identificar configurações de CI/CD;
7. criar ou complementar a estrutura de conhecimento;
8. criar estado persistente do projeto;
9. criar mecanismos de verificação;
10. gerar o manifesto do Harness.

---

# 5. Harness agnóstico

O Harness **não deve ser específico de Codex, Claude Code, GitHub Copilot, Cursor, Gemini ou qualquer outro agente**.

O modelo conceitual será:

```text
                    DARP Harness
                         │
          ┌──────────────┼──────────────┐
          │              │              │
      Knowledge       Workflow      Verification
          │              │              │
          └──────────────┼──────────────┘
                         │
                  Agent Interface
                         │
          ┌──────────────┼──────────────┐
          ▼              ▼              ▼
        Codex          Claude        Copilot
```

O DARP poderá posteriormente gerar adaptações específicas para determinados agentes, mas essas adaptações deverão ser derivadas do Harness e não constituir fontes independentes de verdade.

---

# 6. Harness como propriedade do repositório

O Harness deve ser tratado como parte do próprio projeto.

A intenção é que um desenvolvedor possa clonar um repositório e encontrar nele tudo que um agente precisa para começar a trabalhar de maneira consistente.

Conceitualmente:

```text
Repository
│
├── source code
├── tests
├── infrastructure
├── documentation
│
└── DARP Harness
    ├── knowledge
    ├── state
    ├── workflow
    ├── verification
    ├── policies
    └── agent adapters
```

O Harness não deve depender de um estado mantido exclusivamente fora do Git.

---

# 7. Repositórios existentes

Quando o repositório já possui código, o DARP deverá realizar uma análise inicial antes de gerar o Harness.

Exemplos de artefatos que podem ser analisados:

```text
pom.xml
build.gradle
package.json
pyproject.toml
requirements.txt
go.mod
Cargo.toml

Dockerfile
docker-compose.yml

.github/
.gitlab/
src/
app/
backend/
frontend/
tests/
test/
docs/
```

A detecção deverá ser extensível.

O DARP não deve assumir que a presença de um único arquivo determina toda a stack do projeto.

---

# 8. Detecção de stack

O Harness deverá ser adaptável à stack encontrada.

Exemplo conceitual:

```yaml
stack:
  language: java
  framework: quarkus
  build: maven
  test: junit
```

Entretanto, detecção não deve ser tratada automaticamente como verdade absoluta.

Quando apropriado, o DARP poderá representar:

* tecnologia detectada;
* nível de confiança;
* informações desconhecidas;
* informações que precisam ser confirmadas.

Exemplo conceitual:

```yaml
detection:
  confidence: high

detected:
  - java
  - quarkus
  - maven

unknown:
  - deployment
```

A especificação definitiva desse mecanismo será definida na próxima spec.

---

# 9. Repositório vazio

O DARP também deverá funcionar quando não existir código.

Exemplo:

```bash
mkdir meu-projeto
cd meu-projeto

darp harness init
```

Nesse cenário, o DARP deverá criar um **Generic Harness**.

O Harness genérico:

* não deve assumir uma linguagem;
* não deve assumir um framework;
* não deve assumir um sistema de build;
* não deve assumir um banco de dados;
* não deve assumir uma ferramenta de CI/CD;
* não deve assumir um agente de IA.

Ele deverá fornecer apenas as estruturas e contratos necessários para que o projeto possa evoluir.

---

# 10. Knowledge Layer

Um princípio importante adotado é evitar concentrar todo o conhecimento do projeto em um único arquivo de instruções.

O `AGENTS.md`, quando utilizado, deve funcionar principalmente como um **ponto de entrada/mapa**, enquanto o conhecimento detalhado deverá ser organizado em documentos próprios.

Conceitualmente:

```text
AGENTS.md

docs/
├── index.md
├── architecture/
│   └── index.md
├── decisions/
│   └── index.md
├── development/
│   └── index.md
└── generated/
    └── index.md
```

A estrutura final ainda será definida na especificação do Harness.

---

# 11. Estado persistente

Agentes de longa duração precisam de mecanismos para preservar contexto entre sessões.

Por isso o DARP deverá possuir uma camada de estado persistente.

Conceitualmente:

```text
.darp/
└── state/
    ├── progress.yaml
    └── tasks.yaml
```

Essa camada poderá futuramente armazenar informações como:

* progresso;
* tarefas;
* estado de execução;
* handoffs;
* descobertas;
* checkpoints;
* informações necessárias para continuidade.

O formato definitivo ainda será definido.

---

# 12. Verification Layer

Um Harness não deve apenas orientar o agente.

Ele deve permitir verificar se o trabalho produzido está correto.

O DARP deverá possuir uma camada de verificação capaz de representar mecanismos como:

```text
build
test
lint
format
static analysis
security checks
custom validation
```

Quando possível, o DARP deverá detectar automaticamente os comandos existentes no projeto.

Exemplos:

```bash
mvn test
npm test
pytest
go test ./...
```

O DARP não deverá inventar comandos inexistentes.

---

# 13. Evolução do `darp doctor`

O `darp doctor`, já existente, deverá futuramente evoluir para também avaliar o Harness.

Conceitualmente:

```bash
darp doctor
```

poderá produzir algo como:

```text
✓ Harness manifest
✓ Repository instructions
✓ Repository knowledge
✓ Project state
✓ Verification commands
✓ Git configuration
✓ CI configuration

⚠ Architecture documentation missing
⚠ Automated tests not detected
```

Poderá existir posteriormente um modo específico:

```bash
darp doctor --harness
```

A decisão sobre comandos, opções e critérios ficará para uma futura especificação.

---

# 14. DARP Harness Specification

Antes da implementação de `darp harness init`, deverá ser criada uma especificação formal.

A próxima spec deverá definir pelo menos:

* estrutura do Harness;
* manifesto;
* diretórios;
* arquivos;
* responsabilidades de cada arquivo;
* detecção de stack;
* Harness genérico;
* camada de conhecimento;
* camada de estado;
* workflow;
* verification;
* políticas;
* compatibilidade com agentes;
* adapters;
* comportamento em repositórios existentes;
* comportamento em repositórios vazios;
* idempotência;
* atualização de Harness existente;
* regras para não sobrescrever arquivos do usuário;
* testes;
* critérios de aceitação.

### Próxima especificação planejada

```text
spec/002-harness.md
```

A implementação **não deve começar antes da definição dessa especificação**.

---

# 15. DARP Harness Manifest

O Harness deverá possuir um manifesto próprio.

A estrutura abaixo é apenas conceitual e **não constitui ainda o contrato definitivo**:

```yaml
version: "1"

harness:
  type: generic

project:
  detection: automatic

knowledge:
  root: docs/

workflow:
  state: .darp/state/

verification:
  enabled: true

agents:
  compatible: any

marketplaces:
  enabled: true
```

O formato definitivo deverá ser estabelecido na `spec/002-harness.md`.

---

# 16. Source of Truth

Um princípio arquitetural importante será evitar múltiplas fontes conflitantes de instruções.

O DARP deverá possuir uma fonte de verdade para o Harness.

Arquivos específicos de agentes poderão existir como **adapters/projeções**.

Exemplo conceitual:

```text
DARP Harness
      │
      ├── AGENTS.md
      ├── CLAUDE.md
      ├── GitHub Copilot instructions
      └── outros adapters
```

Esses arquivos não devem divergir conceitualmente do Harness.

O objetivo é evitar situações como:

```text
AGENTS.md              → Maven
CLAUDE.md              → Gradle
Copilot instructions  → Maven
```

A informação deve ser derivada de uma fonte comum sempre que possível.

---

# 17. Assets de IA

Após a implementação inicial do Harness, o DARP deverá evoluir para trabalhar com **AI Assets**.

O conceito de Asset deverá ser genérico.

Exemplos:

```text
skill
agent
instruction
workflow
hook
prompt
MCP
policy
template
```

A definição exata dos tipos de Asset será feita em uma futura especificação.

---

# 18. Marketplace

O DARP deverá futuramente suportar marketplaces de AI Assets.

A intenção é permitir que o projeto utilize assets provenientes de diferentes ecossistemas.

Exemplo conceitual:

```yaml
marketplaces:
  - name: awesome-copilot
    type: git
    source: github/awesome-copilot

  - name: anthropic
    type: git
    source: <anthropic-marketplace>

  - name: darp
    type: darp
    source: <darp-registry>
```

Os marketplaces não deverão ser necessariamente controlados pelo DARP.

O DARP deverá possuir uma arquitetura capaz de consumir diferentes fontes.

---

# 19. DARP Marketplace

No futuro poderá existir um marketplace oficial do DARP.

Conceitualmente:

```text
DARP Marketplace
│
├── Skills
├── Agents
├── Instructions
├── Workflows
├── Hooks
├── MCP
├── Policies
└── Harness Templates
```

Esse marketplace poderá utilizar o DARP Registry como infraestrutura.

Entretanto:

> **Marketplace não faz parte da primeira implementação do Harness.**

A arquitetura deve apenas ser preparada para essa evolução.

---

# 20. DARP Registry

O conceito de DARP Registry continuará sendo uma evolução futura.

O Registry poderá funcionar como catálogo e distribuição de:

```text
Harnesses
Assets
Skills
Agents
MCP
Templates
Policies
```

Visão conceitual:

```text
                    DARP Registry
                         │
             ┌───────────┴───────────┐
             │                       │
         Harnesses                 Assets
             │                       │
             │          ┌────────────┼────────────┐
             │          │            │            │
             │       Skills       Agents         MCP
             │
             ▼
          DARP CLI
```

---

# 21. DARP CLI e agentes

O DARP não pretende substituir agentes de desenvolvimento.

Agentes como:

* Codex;
* Claude Code;
* GitHub Copilot;
* Cursor;
* Gemini;
* outros agentes futuros;

continuarão sendo responsáveis pela execução do trabalho.

O DARP deverá fornecer o ambiente no qual esses agentes trabalham.

Portanto:

```text
Agent ≠ Harness
```

e:

```text
DARP ≠ Agent
```

A relação pretendida é:

```text
DARP Harness
      │
      ▼
   Agent
      │
      ▼
 Repository
```

---

# 22. O que o DARP não deve se tornar

As seguintes direções devem ser evitadas:

### 22.1 Gerador de prompts

O DARP não deve se resumir a gerar prompts para agentes.

### 22.2 Framework específico de agente

O DARP não deve ser construído exclusivamente para Codex, Claude, Copilot ou qualquer outro agente.

### 22.3 Framework de multi-agent

Planner/worker/reviewer e outras arquiteturas multi-agent poderão existir futuramente, mas não fazem parte do núcleo inicial do Harness.

### 22.4 Orquestrador obrigatório

O Harness não deve exigir que o DARP controle diretamente a execução dos agentes.

### 22.5 Dependência de um marketplace

O DARP deve funcionar sem depender de um marketplace externo.

---

# 23. Roadmap de evolução

A sequência planejada passa a ser:

```text
FASE 1 — Fundação
│
├── darp init
├── darp version
├── darp help
└── darp doctor
│
▼
FASE 2 — Harness
│
├── DARP Harness Specification
├── darp harness init
├── stack detection
├── Generic Harness
├── Knowledge Layer
├── State Layer
└── Verification Layer
│
▼
FASE 3 — Harness Quality
│
├── evolução do darp doctor
├── harness validation
├── consistency checks
└── health/quality checks
│
▼
FASE 4 — AI Assets
│
├── asset model
├── asset discovery
├── asset installation
└── asset management
│
▼
FASE 5 — Marketplaces
│
├── marketplace model
├── external marketplaces
├── marketplace configuration
└── asset synchronization
│
▼
FASE 6 — DARP Registry
│
├── registry
├── catalog
├── publishing
├── versioning
└── distribution
│
▼
FASE 7 — DARP SDK
│
├── programmatic access
├── Harness APIs
├── Asset APIs
└── Registry APIs
│
▼
FASE 8 — Advanced Agent Capabilities
│
├── workflows
├── orchestration
├── multi-agent patterns
└── advanced automation
```

---

# 24. Prioridade atual

A prioridade imediata é:

## `DARP Harness Specification`

Antes de implementar:

```bash
darp harness init
```

deveremos definir claramente:

```text
O que é um Harness?
        ↓
Qual é sua estrutura?
        ↓
Qual é seu manifesto?
        ↓
Como detectar o projeto?
        ↓
Como lidar com projeto vazio?
        ↓
Como representar conhecimento?
        ↓
Como representar estado?
        ↓
Como representar verificação?
        ↓
Como manter compatibilidade com qualquer agente?
        ↓
Como evoluir sem quebrar o projeto?
```

Somente depois dessa definição deverá ser criado o prompt de implementação para o VSCode/Codex.

---

# 25. Decisões registradas em 2026-08-08

As seguintes decisões foram tomadas:

1. **Harness Engineering será incorporado ao DARP CLI.**

2. O conceito será chamado inicialmente de **DARP Harness**.

3. O Harness será **agnóstico ao agente de IA**.

4. O Harness será uma propriedade do **repositório/projeto**, e não de um agente específico.

5. O DARP deverá funcionar tanto em repositórios existentes quanto em repositórios vazios.

6. Repositórios vazios receberão um **Generic Harness**.

7. Repositórios existentes deverão passar por **detecção de stack e estrutura**.

8. `AGENTS.md` não será tratado como o Harness inteiro.

9. O conhecimento do projeto deverá ser organizado em uma **Knowledge Layer**.

10. O Harness deverá possuir uma **State Layer** para permitir continuidade entre sessões.

11. O Harness deverá possuir uma **Verification Layer**.

12. O `darp doctor` deverá futuramente verificar a qualidade do Harness.

13. O DARP deverá possuir um **manifesto próprio do Harness**.

14. Arquivos específicos de agentes deverão ser tratados como **adapters/projeções**, quando necessários.

15. O DARP deverá permanecer independente de Codex, Claude, Copilot, Cursor, Gemini ou outros agentes.

16. O conceito de **AI Assets** será incorporado futuramente.

17. Assets poderão incluir, entre outros:

    * skills;
    * agents;
    * instructions;
    * workflows;
    * hooks;
    * prompts;
    * MCP;
    * policies;
    * templates.

18. O DARP deverá futuramente suportar **marketplaces**.

19. O marketplace não será implementado junto com o primeiro Harness.

20. A arquitetura deverá ser preparada para um futuro **DARP Marketplace**.

21. O **DARP Registry** continuará como componente futuro do ecossistema.

22. Multi-agent orchestration não faz parte do primeiro Harness.

23. O próximo artefato técnico a ser produzido será:

```text
spec/002-harness.md
```

24. A implementação de `darp harness init` deverá começar somente após a aprovação dessa especificação.

---

# 26. Referências conceituais

As decisões deste roadmap foram influenciadas principalmente pela evolução recente do conceito de Harness Engineering e pelas experiências publicadas por fornecedores de agentes.

Referências principais:

* OpenAI — Harness Engineering
* OpenAI — Agents SDK
* Anthropic — Effective Harnesses for Long-Running Agents
* Anthropic — Harness Design for Long-Running Application Development
* GitHub — Awesome Copilot e ecossistema de plugins/AI assets

Essas referências servem como **inspiração arquitetural**, não como dependências do DARP.

O DARP deverá possuir sua própria especificação e seus próprios contratos.

---

# 27. Princípio final

A direção estratégica do DARP pode ser resumida em:

> **Transformar qualquer repositório em um ambiente estruturado, verificável e evolutivo para desenvolvimento com agentes de IA, independentemente do agente utilizado.**

Ou, de forma mais curta:

```text
DARP
=
Agent-Ready Development Environment
```

O Harness é o primeiro grande passo para atingir essa visão.
