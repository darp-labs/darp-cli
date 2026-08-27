# DARP Agentic Harness

> Alinhamento estratégico do DARP CLI com a evolução dos agentes de IA para engenharia de software.

**Status:** Proposta estratégica  
**Data:** 2026-08-15  
**Escopo:** arquitetura e direção de produto; não implica implementação imediata.

---

## 1. Conclusão estratégica

O mercado está convergindo para uma nova camada de desenvolvimento:

```text
Modelos
   ↓
Agentes
   ↓
Harness
   ↓
Tools / MCP / Skills
   ↓
CI/CD / Evaluation / Governance
   ↓
Produção
```

A principal conclusão é:

> **O diferencial está deixando de ser apenas qual LLM é utilizado e passando a ser o sistema que conecta agente, contexto, ferramentas, segurança, feedback e execução.**

O DARP deve aproveitar esse movimento sem se transformar em outro Copilot, Claude Code, Codex ou Cursor.

### Posicionamento proposto

> **DARP é uma camada de assets, governança e harness para engenharia de software agentic, independente do modelo utilizado.**

O modelo deve ser substituível.

---

## 2. Alinhamento com a visão atual do DARP

A visão existente do DARP como plataforma/pacote para assets de IA continua válida.

```text
DARP
 |
 +-- Prompts
 +-- Instructions
 +-- Skills
 +-- Personas
 +-- MCP Servers
 +-- Workflows
 +-- Templates
 +-- Context Packages
 |
 +-- Harness definitions
 +-- Provider profiles
 +-- Tool policies
 +-- Evaluation contracts
```

O Harness não substitui os assets. Ele os **compõe** em um ambiente reproduzível de execução.

---

## 3. O que é um Harness

Um Harness é a composição declarativa dos elementos necessários para um agente executar uma classe de tarefas de engenharia com segurança e critérios objetivos.

Um Harness pode definir:

- modelo/provedor preferencial;
- fallback de modelos;
- instruções;
- skills;
- MCP servers;
- ferramentas;
- permissões;
- contexto do repositório;
- workflows;
- comandos de teste;
- quality gates;
- avaliação;
- limites de execução;
- requisitos de segurança;
- observabilidade.

Conceitualmente:

```text
Harness
 |
 +-- Provider / Model
 +-- Instructions
 +-- Skills
 +-- MCP Servers
 +-- Tools
 +-- Permissions
 +-- Context
 +-- Workflow
 +-- Quality Gates
 +-- Evaluation
 +-- Fallback Strategy
```

O Harness **não é o modelo**.

---

## 4. DARP deve ser multi-modelo

O mesmo Harness deve poder trabalhar com diferentes providers:

```text
                    DARP Harness
                         |
          +--------------+--------------+
          |              |              |
       Claude          Codex         Local LLM
          |              |              |
      Anthropic         OpenAI        Ollama/vLLM
```

Exemplo:

```text
Tarefa simples
    → modelo local

Implementação Java
    → Codex

Arquitetura complexa
    → Claude

Código sensível
    → modelo local

Falha do provider primário
    → fallback
```

A primeira implementação não deve criar um roteador inteligente de modelos. Primeiro deve existir uma **abstração declarativa de providers**.

---

## 5. Modelos locais

Modelos open-weight tornam o suporte a execução local relevante para o DARP.

Possíveis backends:

- Ollama
- llama.cpp
- vLLM
- endpoints compatíveis com OpenAI API
- outros runtimes locais

### Princípio

O DARP **não deve executar inferência de modelos** na primeira fase.

Ele deve apenas declarar e integrar providers/runtimes.

Isso preserva o princípio de não transformar o DARP em um runtime de LLM.

---

## 6. MCP como camada de ferramentas

MCP deve continuar sendo um asset de primeira classe.

```text
Harness
 |
 +-- MCP GitHub
 +-- MCP Database
 +-- MCP Kubernetes
 +-- MCP Documentation
 +-- MCP Observability
```

O manifesto de um MCP futuramente poderá declarar:

- identidade;
- versão;
- transporte;
- endpoint;
- capabilities;
- secrets necessários;
- permissões;
- requisitos de ambiente;
- compatibilidade;
- metadados de segurança.

MCP deve ser tratado como infraestrutura de produção, não apenas integração experimental.

---

## 7. Skills como ecossistema de pacotes

O crescimento do formato `SKILL.md` reforça diretamente a visão original do DARP como package manager.

Uma Skill deve ser um asset versionado com:

- identidade;
- versão;
- descrição;
- inputs;
- dependências;
- providers suportados;
- agentes suportados;
- ferramentas necessárias;
- permissões;
- scripts;
- documentação;
- evidências de avaliação.

Exemplo:

```text
java-quarkus-review@1.2.0
 |
 +-- instructions
 +-- skills
 +-- MCP dependencies
 +-- validation workflow
 +-- evaluation contract
```

---

## 8. Segurança e supply chain

À medida que skills, MCP servers e plugins passam a influenciar agentes, eles se tornam uma nova superfície de supply chain.

Um asset malicioso pode potencialmente:

- acessar código;
- executar comandos;
- acessar MCP;
- ler credenciais;
- modificar arquivos;
- influenciar decisões do agente.

Portanto, o DARP deve evoluir para suportar:

- provenance;
- identidade do publisher;
- checksums/assinaturas;
- versionamento;
- lockfiles;
- declarações de permissões;
- níveis de confiança;
- security scanning;
- instalação reproduzível;
- auditoria.

### Princípio

> **Um AI asset deve ser tratado como uma dependência de software com comportamento.**

---

## 9. Evaluation como primeira classe

Package managers tradicionais respondem:

> O pacote instala?

DARP deve eventualmente responder:

> O asset/harness produz o resultado de engenharia esperado?

Um Evaluation Contract pode conter:

```text
Evaluation
 |
 +-- Task
 +-- Repository fixture
 +-- Expected behavior
 +-- Tests
 +-- Quality gates
 +-- Security checks
 +-- Cost limit
 +-- Time limit
 +-- Result score
```

Isso permite comparar:

```text
Claude + Harness A
Codex + Harness A
Gemini + Harness A
Local Qwen + Harness A
```

usando a mesma tarefa e os mesmos critérios.

---

## 10. Arquitetura proposta

```text
                         DARP CLI
                            |
               +------------+------------+
               |                         |
            Registry                  Local Cache
               |                         |
               +------------+------------+
                            |
                          Assets
                            |
                    Harness Resolver
                            |
          +-----------------+-----------------+
          |                 |                 |
       Provider           Tools             Policy
          |                 |                 |
       Claude            MCP/Git           Security
       Codex             Shell             Permissions
       Gemini            CI/CD             Limits
       Local             APIs              Audit
          |                 |                 |
          +-----------------+-----------------+
                            |
                       Evaluation
                            |
                     Engineering Result
```

---

## 11. Novos conceitos candidatos

### Harness
Composição declarativa de providers, assets, tools, policies e evaluation.

### Provider Profile
Configuração de um backend de raciocínio e suas capabilities.

### Tool Policy
Regras que definem quais ferramentas podem ser utilizadas.

### Evaluation Contract
Definição reproduzível de como o resultado será validado.

### Agent Workflow
Fluxo reutilizável para uma classe de tarefas.

### Capability
Representação normalizada das capacidades de um asset/provider.

### Compatibility
Declaração dos agentes, providers, runtimes e stacks compatíveis.

---

## 12. Relação com `darp init` e `darp doctor`

### `darp init`

No futuro poderá inicializar:

- baseline do Harness;
- providers;
- policies;
- evaluation contracts;
- MCP/skills declarations.

### `darp doctor`

Poderá validar:

- providers configurados;
- modelos disponíveis;
- MCP servers;
- integridade dos assets;
- conflitos de permissões;
- ferramentas ausentes;
- incompatibilidades de versão;
- readiness para evaluation.

O `darp doctor` deve continuar sendo essencialmente diagnóstico/read-only.

---

## 13. O que DARP não deve se tornar

### Não é um LLM provider
Não deve treinar ou hospedar modelos frontier.

### Não é uma IDE
Deve integrar com IDEs.

### Não é um clone de coding assistant
Não deve replicar Claude Code, Codex ou Copilot feature-by-feature.

### Não é um model runtime
Deve integrar com runtimes locais.

### Não é um workflow engine genérico
Workflows devem permanecer focados em AI-assisted engineering.

### Não é um marketplace na primeira fase
Registry e package semantics devem amadurecer antes de uma marketplace.

---

## 14. Impacto no roadmap

O Harness deve ser introduzido incrementalmente.

### Fase 1 — Foundation

1. Definir conceito de Harness.
2. Definir abstração de provider.
3. Definir composição de assets.
4. Definir vocabulário de segurança e permissões.

### Fase 2 — Specification

5. Harness Manifest Specification.
6. Provider Profile Specification.
7. Tool Policy Specification.
8. Evaluation Contract Specification.

### Fase 3 — Implementation

9. Harness discovery.
10. Harness validation.
11. Provider detection.
12. Local provider metadata/integration.
13. Evaluation support.

### Fase 4 — Ecosystem

14. Registry para Harness packages.
15. Dependency resolution.
16. Provenance/signing.
17. Compatibility metadata.
18. Community Harnesses.

---

## 15. Primeira Spec recomendada

A primeira especificação formal **não deve ser `darp harness run`**.

A recomendação é:

> **Harness Manifest — formato declarativo e provider-agnostic para descrever assets, providers, ferramentas, permissões e quality gates necessários para um workflow agentic de engenharia.**

Motivos:

- estabelece a abstração central;
- pode ser validado sem executar um agente;
- preserva o não-objetivo de executar modelos;
- cria contrato para futuras APIs/CLI;
- permite registry;
- habilita portabilidade entre providers;
- encaixa naturalmente no SDD.

### Exemplo conceitual

```yaml
apiVersion: darp.dev/v1alpha1
kind: Harness

metadata:
  name: java-quarkus-development
  version: 0.1.0

provider:
  preferred: anthropic/claude
  fallback:
    - openai/codex
    - local/qwen-coder

assets:
  skills:
    - java-review
    - quarkus-development
    - testing

  mcpServers:
    - github
    - documentation

tools:
  shell: true
  git: true
  filesystem: workspace

policy:
  network: restricted
  secrets: denied
  requireTests: true

evaluation:
  required:
    - build
    - tests
    - lint
    - security
```

**Importante:** esse YAML é apenas ilustrativo. Não deve ser implementado como schema definitivo antes da Spec.

---

## 16. Princípios arquiteturais

A evolução para Harness adiciona estes princípios:

1. **Model agnostic** — o Harness não depende de um provider.
2. **Asset first** — capabilities são assets reutilizáveis e versionados.
3. **Declarative** — descreve capacidades desejadas.
4. **Reproducible** — o mesmo Harness produz ambiente comparável.
5. **Least privilege** — agentes recebem somente permissões necessárias.
6. **Evaluable** — tarefas possuem critérios objetivos.
7. **Observable** — ações importantes são auditáveis.
8. **Composable** — skills, MCP, providers e workflows compõem-se.
9. **Portable** — assets funcionam em ambientes compatíveis.
10. **Open standards** — priorizar padrões interoperáveis.

---

## 17. Conclusão

A oportunidade estratégica do DARP está em ocupar a camada entre **AI assets e engineering execution**.

```text
                DARP
                 |
        Assets + Harness
                 |
     +-----------+-----------+
     |           |           |
   Claude      Codex      Local
     |           |           |
     +-----------+-----------+
                 |
          MCP / Skills / Tools
                 |
          CI/CD / Evaluation
                 |
              Software
```

A aposta estratégica é:

> **DARP deve tornar capacidades de engenharia agentic tão portáveis, versionáveis, reproduzíveis e governáveis quanto os pacotes de software são hoje.**

### Próximo passo

O próximo artefato recomendado é a **Spec do Harness Manifest**, seguindo o ciclo SDD existente no projeto:

```text
Spec
  ↓
Clarification
  ↓
ADR
  ↓
Plan
  ↓
Tasks
  ↓
Implementation
```

Nenhuma implementação de Harness deve começar antes dessa especificação ser aprovada.
