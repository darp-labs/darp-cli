# ADR 002 — Limites e interoperabilidade do DARP Harness

## Status

Accepted — decisão estratégica para orientar a evolução do DARP Harness e das próximas especificações.

## Contexto

O mercado passou a usar *harness* para responsabilidades diferentes. Alguns projetos tratam Harness como runtime completo de agentes; outros como o ambiente de engenharia que fornece contexto, ferramentas, políticas, validação e feedback ao agente.

Ao mesmo tempo, Codex, Claude Code, GitHub Copilot, Gemini, OpenHands, DeepSeek Harness e outros ecossistemas evoluem formatos e runtimes próprios para agents, skills, plugins, hooks, MCP, configuração e execução.

Como o DARP é inicialmente mantido por uma única pessoa, competir diretamente com esses runtimes ou duplicar seus formatos aumentaria custo, acoplamento e risco de obsolescência.

A oportunidade do DARP é resolver um problema transversal: preparar, analisar, validar, compor e manter repositórios para diferentes agentes, sem assumir a execução de nenhum deles.

## Decisão

### 1. DARP não será um Agent Runtime

O DARP não implementará como objetivo central:

- loop de agente;
- gerenciamento de sessões;
- inferência de modelos;
- sandbox de execução;
- runtime de ferramentas;
- execução de MCP servers;
- orquestração obrigatória de agentes;
- substituição de Codex, Claude Code, GitHub Copilot, Gemini, OpenHands, DeepSeek Harness ou outros coding agents.

A execução permanece responsabilidade do agente ou runtime escolhido pelo projeto.

### 2. DARP será uma camada transversal de Harness Engineering

O foco do DARP deverá ser:

- descoberta e diagnóstico do repositório;
- descrição declarativa do estado desejado do Harness;
- validação estrutural e de consistência;
- geração e manutenção de artefatos agent-ready;
- composição de capabilities e AI Assets;
- compatibilidade entre agentes;
- migração e projeção entre formatos compatíveis;
- governança, provenance e supply chain de assets;
- verificação da qualidade do ambiente preparado para agentes.

### 3. Interoperabilidade é princípio de primeira classe

Quando existir padrão aberto ou formato amplamente adotado que resolva o problema, o DARP deverá preferi-lo a um formato proprietário.

Ordem de preferência:

1. padrões abertos e especificações interoperáveis;
2. formatos nativos dos ecossistemas-alvo;
3. contratos DARP somente quando houver necessidade transversal não atendida adequadamente pelos formatos existentes.

O DARP não deve criar uma abstração própria apenas para uniformizar diferenças cosméticas.

### 4. DARP não substitui Agents, Skills, Plugins, Hooks ou MCP

Esses mecanismos pertencem aos ecossistemas dos agentes. O DARP poderá compreendê-los, validar, catalogar, compor ou instalar quando houver contrato explícito, mas não deverá exigir conversão para um formato DARP proprietário.

Em particular:

- Skills devem manter o formato suportado pelo ecossistema de destino;
- Plugins devem manter seu formato e estrutura originais;
- MCP Servers continuam seguindo o protocolo MCP;
- instruções nativas devem ser geradas ou mantidas pelo adapter correspondente.

### 5. O Harness DARP é declarativo, não runtime

`.darp/harness.yaml` representa intenção e composição declarativa. Ele não concede por si só permissão para executar comandos, acessar rede, utilizar secrets, iniciar MCP servers ou executar agentes.

Qualquer execução futura exigirá contrato próprio, modelo de segurança, permissões, observabilidade e critérios de aprovação.

### 6. Adapters isolam diferenças entre ecossistemas

Regras específicas de agentes não devem ficar espalhadas pelo core. Integrações deverão ser isoladas em adapters que conheçam os formatos e capacidades de cada ecossistema.

```text
                       DARP Core
                           |
                 Harness / Asset Model
                           |
          +----------------+----------------+
          |                |                |
       Codex Adapter   Claude Adapter   Copilot Adapter
          |                |                |
       Native files     Native files     Native files
```

Novos agentes devem poder ser adicionados sem alterar o modelo central quando suas diferenças puderem ser isoladas no adapter.

### 7. Não existe promessa de abstração total

DARP não deve esconder diferenças semânticas relevantes entre ecossistemas.

`Skill`, `Plugin`, `Agent`, `Instruction`, `Hook`, `MCP Server` e outros tipos podem possuir responsabilidades e modelos de segurança diferentes.

Uma normalização interna só é válida quando preservar a semântica necessária. Quando isso não for possível, o DARP deve representar a capacidade como específica do target.

### 8. DARP deve funcionar sem agente específico instalado

Diagnóstico, análise, validação e preparação não devem exigir Codex, Claude, Copilot, Gemini ou outro agente.

Um repositório pode possuir Harness DARP válido mesmo sem target instalado localmente.

### 9. Execução é evolução opcional

Descrever providers, tools, MCP, workflows ou quality gates no manifesto não implica que o DARP deva executá-los.

Se execução se tornar necessária, deverá ser avaliada como nova responsabilidade e poderá justificar runtime separado.

## Alternativas consideradas

### A. Criar um Agent Harness Runtime

Rejeitada por colocar o DARP em competição direta com runtimes e coding agents de grande escala e aumentar a superfície operacional e de segurança.

### B. Criar formato DARP proprietário para Skills e Plugins

Rejeitada por criar lock-in, dificultar interoperabilidade e duplicar formatos já existentes.

### C. Escolher e executar automaticamente o melhor agente

Rejeitada para o núcleo. Roteamento e execução introduzem dependências, permissões, custo e observabilidade desnecessários ao problema central.

### D. Não possuir contrato DARP próprio

Rejeitada. Um contrato transversal continua necessário para representar o estado desejado do ambiente quando múltiplos agentes e formatos coexistem.

### E. Colocar adapters específicos diretamente no core

Rejeitada para evitar acoplamento e tornar a manutenção de novos targets cara.

## Consequências

### Positivas

- reduz o risco de competição direta com grandes coding agents e runtimes;
- preserva a utilidade do DARP mesmo com mudanças rápidas de fornecedores;
- favorece interoperabilidade;
- permite adoção incremental;
- reduz o custo de manutenção para um projeto mantido por uma única pessoa;
- mantém o Harness como contrato estável enquanto os adapters evoluem;
- cria espaço para funcionalidades diferenciais como doctor, validation, migration, compatibility e governance.

### Negativas

- algumas funcionalidades dependerão da evolução dos formatos externos;
- nem toda capacidade de um agente terá equivalente em outro;
- adapters precisarão ser mantidos conforme os ecossistemas mudarem;
- o DARP poderá precisar reconhecer diferenças entre versões dos targets;
- algumas funcionalidades atraentes de execução permanecerão deliberadamente fora do núcleo.

## Guardrails para futuras especificações

Uma nova feature relacionada a Harness ou AI Assets deverá responder:

1. Essa capacidade atravessa dois ou mais ecossistemas de agentes?
2. Existe um padrão ou formato externo que já resolve o problema?
3. Estamos criando uma abstração DARP porque ela é necessária ou apenas para uniformizar formatos diferentes?
4. A funcionalidade pode ser implementada como adapter sem contaminar o core?
5. A funcionalidade exige execução, permissões, secrets, rede ou observabilidade?
6. Ela cria competição direta com um runtime ou coding agent existente?
7. O custo de manutenção é compatível com um projeto inicialmente mantido por uma pessoa?

Se a resposta indicar dependência de um único agente ou competição direta com seu runtime, a funcionalidade deve ser reconsiderada antes da implementação.

## Impacto no Harness Manifest

Esta decisão não altera o contrato mínimo aprovado na Spec 006 nem o caminho `.darp/harness.yaml`.

Ela estabelece uma regra de evolução: os blocos do manifesto devem representar conceitos transversais e declarativos, enquanto detalhes específicos de agentes, providers, plugins, skills, hooks e runtimes devem ser tratados por contratos ou adapters próprios quando necessários.

## Follow-up Actions

- manter a Spec 006 como contrato declarativo e não executivo;
- revisar futuras specs de Harness contra estes guardrails;
- priorizar descoberta, diagnóstico, validação, compatibilidade e migração antes de runtime próprio;
- priorizar padrões abertos e formatos nativos dos ecossistemas;
- especificar adapters concretos separadamente antes de implementá-los;
- não tornar marketplace ou registry dependência obrigatória do Harness.
