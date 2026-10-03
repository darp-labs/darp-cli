# DARP Harness Strategy

> Direção estratégica revisada do DARP Harness após análise do mercado de agent runtimes, coding agents e padrões de interoperabilidade.

**Status:** Direção estratégica revisada  
**Escopo:** arquitetura, posicionamento e limites de produto; não implica implementação imediata.

## 1. Posicionamento

O DARP deve ocupar a camada transversal de **Harness Engineering para repositórios agent-ready**. Não deve competir com Agent Runtimes ou coding agents.

> **DARP prepara, analisa, valida, compõe e mantém o ambiente do repositório; o agente escolhido executa o trabalho.**

Posicionamento recomendado:

> **DARP — CLI para Harness Engineering e repositórios preparados para agentes, independente do agente utilizado.**

## 2. Limites do produto

O DARP não deve se tornar:

- Agent Runtime;
- coding assistant;
- IDE;
- clone de Codex, Claude Code, GitHub Copilot, Gemini, OpenHands ou outros coding agents;
- runtime de modelos;
- formato proprietário de Skills ou Plugins;
- marketplace obrigatório.

O DARP não deve implementar como responsabilidade central loops de agentes, sessões, inferência, sandbox, runtime de ferramentas, execução de MCP servers, orquestração obrigatória, roteamento inteligente de modelos ou execução operacional de workflows.

A execução continua responsabilidade do agente ou runtime escolhido pelo projeto.

## 3. Problema transversal

Os ecossistemas atuais possuem mecanismos semelhantes, mas não idênticos:

```text
Instructions
Skills
Agents
Plugins
Hooks
MCP
Configuration
Policies
Workflows
```

O DARP deve resolver problemas que atravessam esses ecossistemas:

- diagnosticar se um repositório está preparado para agentes;
- detectar inconsistências entre ambientes de agentes;
- validar capabilities e AI Assets;
- manter intenção declarativa sem apagar formatos nativos;
- projetar ou migrar ambientes entre agentes;
- analisar riscos de supply chain;
- manter a qualidade do Harness ao longo do tempo.

## 4. Interoperabilidade primeiro

A ordem de preferência é:

1. padrão aberto e interoperável;
2. formato nativo do ecossistema de destino;
3. contrato DARP somente quando houver necessidade transversal real não atendida pelos anteriores.

O DARP **não deve criar formato proprietário de Skill ou Plugin apenas para uniformizar diferenças**.

- Skills mantêm o formato suportado pelo ecossistema de destino.
- Plugins mantêm formato e estrutura originais.
- MCP continua seguindo o protocolo MCP.
- Hooks e instructions permanecem específicos do target quando sua semântica for específica.

## 5. Não existe abstração universal obrigatória

```text
Skill != Plugin
Plugin != Agent
Agent != MCP
Hook != Instruction
```

O DARP pode normalizar metadados comuns como identidade, versão, origem, capabilities, compatibilidade, dependências e permissões declaradas, mas não deve apagar diferenças semânticas relevantes.

Quando a equivalência não for segura, o recurso permanece específico do target.

## 6. Harness como contrato declarativo

O Harness DARP é contrato declarativo, não runtime.

O manifesto canônico continua sendo:

```text
.darp/harness.yaml
```

A Spec 006 continua válida: o manifesto é versionável, provider-agnostic, validável estruturalmente e não executa agentes, modelos, ferramentas, MCP ou workflows.

A existência do manifesto não concede autorização para executar comandos, acessar rede ou utilizar secrets.

## 7. Adapters

Diferenças entre ecossistemas devem ser isoladas em adapters:

```text
                    DARP Core
                       |
              Harness / Asset Model
                       |
       +---------------+---------------+
       |               |               |
   Codex Adapter   Claude Adapter   Copilot Adapter
       |               |               |
   Native files    Native files     Native files
```

O core não deve ser contaminado por regras específicas de cada agente. Novos targets devem poder ser adicionados principalmente pela camada de adapter.

## 8. DARP Doctor como diferencial

O `darp doctor` deve evoluir progressivamente para diagnosticar o estado agent-ready do repositório.

Exemplo conceitual:

```text
DARP Doctor

Repository
✓ Git
✓ Java
✓ Quarkus
✓ Maven
✓ CI

Harness
✓ .darp/harness.yaml
✓ Repository knowledge
⚠ Verification commands incomplete
⚠ Architecture documentation missing

Agent Targets
✓ Codex compatible
✓ Claude compatible
⚠ Copilot projection missing

Consistency
✓ No conflicting instructions detected
⚠ Duplicate skill detected
```

Os critérios concretos devem ser definidos em specs próprias.

## 9. Validation e Migration

### Harness Validation

Deve futuramente validar manifesto, referências, consistência, compatibilidade, duplicação, conflitos, segurança declarativa, adapters e artefatos derivados.

### Harness Migration

Deve futuramente projetar um ambiente de um target para outro, preservando o que é portável e identificando explicitamente perdas e incompatibilidades.

```text
Claude Environment -> DARP Analysis -> Cross-Agent Model -> Codex Projection
```

Exemplo futuro:

```bash
darp harness migrate --from claude --to codex
```

A implementação exige especificação própria.

## 10. AI Assets

AI Assets são recursos de ecossistema, não um formato proprietário obrigatório do DARP.

```text
Skill | Agent | Instruction | Plugin | Hook | MCP | Prompt | Workflow | Policy | Template
```

O DARP pode descobrir, catalogar, validar, analisar, comparar, verificar compatibilidade e instalar quando houver contrato seguro.

**O formato original do asset deve ser preservado.**

## 11. Supply chain

AI Assets são dependências com comportamento. Podem conter scripts, hooks, MCP, comandos, acesso à rede, permissões e instruções que influenciam o agente.

Evoluções futuras devem considerar provenance, publisher, checksums, assinaturas, lockfiles, permissões declaradas, análise de conteúdo, níveis de confiança, instalação reproduzível e auditoria.

**Instalação não deve significar execução automática.**

## 12. Evaluation

Evaluation continua importante, mas é diferente de execução do Harness.

Um futuro Evaluation Contract poderá representar tarefa, fixture, comportamento esperado, testes, quality gates, security checks, limites de custo/tempo e resultado.

A execução concreta de avaliações deverá ser definida separadamente.

## 13. O que não fazer

- Não criar DARP Agent Runtime.
- Não criar DARP Skill Format.
- Não criar DARP Plugin Format.
- Não tornar marketplace uma dependência.
- Não criar abstração universal artificial.
- Não implementar roteamento inteligente de modelos no núcleo.
- Não colocar integrações específicas de agentes diretamente no core.

## 14. Guardrails para novas funcionalidades

Toda nova feature de Harness ou AI Assets deve responder:

1. A capacidade atravessa dois ou mais ecossistemas?
2. Já existe um padrão aberto?
3. Já existe formato nativo que deve ser preservado?
4. A abstração DARP resolve uma necessidade real ou apenas uniformiza diferenças?
5. Pode ser implementada como adapter?
6. Exige execução, rede, secrets, permissões ou observabilidade?
7. Coloca o DARP em competição direta com um produto maior?
8. O custo é compatível com um projeto inicialmente mantido por uma pessoa?

> **Se a feature depender essencialmente de um único agente ou competir diretamente com seu runtime, deve ser reconsiderada.**

## 15. Arquitetura estratégica

```text
                         DARP CLI
                            |
             +--------------+--------------+
             |              |              |
          Discover       Diagnose       Validate
             |              |              |
             +--------------+--------------+
                            |
                         Compose
                            |
                         Generate
                            |
                         Maintain
                            |
                    Harness Contract
                            |
             +--------------+--------------+
             |              |              |
          Adapters       Assets        Policies
             |              |              |
       +-----+-----+        |         Governance
       |     |     |        |
     Codex Claude Copilot  Skills/MCP/etc.
       |     |     |
       +-----+-----+
             |
         Repository
```

## 16. Roadmap estratégico revisado

### Fase atual — Harness Contract

- `.darp/harness.yaml`;
- parsing;
- validação estrutural;
- extensibilidade;
- compatibilidade declarativa.

Essa fase já está coberta pela Spec 006.

### Próxima prioridade

- Harness diagnosis;
- evolução do `darp doctor`;
- consistency checks;
- compatibility model;
- agent adapters;
- validação de projeções.

### Depois

- asset discovery;
- asset validation;
- provenance;
- security analysis;
- safe installation;
- migration/projection.

### Futuro

- Registry;
- SDK;
- evaluation;
- automações avançadas;
- integrações adicionais.

Um runtime próprio de agentes **não faz parte do roadmap padrão** e só deve ser considerado se surgir requisito concreto que não possa ser atendido mantendo o DARP como camada transversal.

## 17. Princípios consolidados

1. **Cross-agent** — capacidades prioritárias atravessam múltiplos agentes.
2. **Agent-agnostic core** — o núcleo não depende de um agente.
3. **Adapter isolation** — diferenças ficam em adapters.
4. **Open standards first** — padrões abertos têm prioridade.
5. **Native format preservation** — formatos externos não são substituídos sem necessidade.
6. **Declarative first** — contratos precedem execução.
7. **Read-only by default** — diagnóstico e validação não alteram o projeto.
8. **Least privilege** — nenhuma declaração concede execução automaticamente.
9. **Composable** — assets e capabilities podem ser compostos.
10. **Reproducible** — ambientes devem ser versionáveis e comparáveis.
11. **Evaluable** — qualidade deve poder ser verificada objetivamente.
12. **Maintainable by one person** — complexidade deve ser proporcional à capacidade de manutenção.
13. **No unnecessary competition** — o DARP complementa, não substitui, os grandes ecossistemas.

## 18. Relação com a Spec 006

A Spec 006 — Harness Manifest continua válida.

Este documento **não altera o contrato implementado**; estabelece limites estratégicos para as próximas evoluções.

A decisão complementar está registrada em:

```text
.spec/specs/006-harness-manifest/adr/002-cross-agent-harness-boundaries.md
```

## 19. Próximos documentos

Nenhuma implementação adicional de Harness deve começar apenas com base neste documento.

O ciclo continua:

```text
Necessidade -> Spec -> ADR quando necessário -> Plan -> Tasks -> Implementation -> Validation
```

As próximas specs devem aplicar os guardrails deste documento e do ADR 002.

## 20. Conclusão

A oportunidade do DARP não está em construir outro agente ou outro runtime. Está em resolver a fragmentação criada pela existência de muitos agentes.

> **DARP deve tornar repositórios preparados para agentes mais portáveis, verificáveis, governáveis e sustentáveis, independentemente do agente utilizado.**

O agente muda. O modelo muda. O runtime muda. O Harness Engineering do repositório pode permanecer como uma camada estável.
