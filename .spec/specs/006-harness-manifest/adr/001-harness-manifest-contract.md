# ADR 001 — Contrato independente para o Harness Manifest

## Status

Accepted — aprovado para implementação como contrato declarativo.

## Contexto

O DARP precisa representar a composição de assets e requisitos de execução
agentic sem se tornar um runtime de modelos ou depender de um fornecedor. O
projeto já possui `darp.yml` para configuração e registro do projeto e possui
contratos de assets locais. Misturar a definição do Harness nesses contratos
criaria acoplamento e dificultaria a evolução independente.

Esta decisão estabelece também um limite de produto: o DARP é responsável pelo
contrato, empacotamento, validação e governança declarativa do Harness. Um
runtime futuro poderá consumir esse contrato, mas execução de agentes,
providers, MCP servers, workflows, permissões, segredos, observabilidade e
avaliação operacional não pertencem ao DARP CLI nesta fase.

## Decisão

1. O manifesto será um arquivo separado em `.darp/harness.yaml`.
2. O formato usará YAML com identificação explícita:

   ```yaml
   apiVersion: darp.dev/v1alpha1
   kind: Harness
   metadata:
     name: example-harness
     version: 0.1.0
   ```

3. Assets, providers, ferramentas, políticas, workflows, quality gates,
   avaliação e compatibilidade serão dados declarativos e extensíveis. Os
   blocos cuja semântica ainda não possui contrato próprio serão apenas
   preservados e validados estruturalmente.
4. A primeira versão suportará validação estrutural, não execução.
5. Providers serão identificadores intercambiáveis; nenhum provider será
   obrigatório por decisão implícita do formato.
6. Campos desconhecidos serão aceitos e preservados na representação de leitura
   para forward compatibility. A preservação de comentários e da ordem textual
   original não é uma garantia desta versão.
7. Mudanças incompatíveis exigirão nova versão de `apiVersion`, sem migração
   silenciosa.

## Alternativas consideradas

### Incorporar o Harness em `darp.yml`

Rejeitada porque mistura configuração do projeto com composição declarativa de
execução e aumenta o acoplamento entre `init`, `doctor`, assets e Harness.

### Usar um formato específico de provider ou agente

Rejeitada porque viola a Constituição provider-agnostic e impediria a mesma
definição de ser usada com Codex, Claude, Gemini ou runtimes locais.

### Implementar primeiro `darp harness run`

Rejeitada porque exigiria decisões prematuras sobre execução, permissões,
segredos, observabilidade e segurança. Também contrariaria o direcionamento
estratégico de definir o contrato antes do executor.

### Criar um CLI ou runtime separado imediatamente

Adiada porque ainda não existe uma responsabilidade operacional concreta que
justifique um segundo produto. O contrato permanece no DARP, enquanto um
futuro runtime poderá ser separado quando houver requisitos reais de execução,
permissões, segredos e observabilidade.

### Usar JSON ou um formato proprietário

Rejeitada nesta fase porque YAML já é usado pelos contratos do projeto e é
legível para configuração versionada. O contrato continua baseado em dados
portáveis, sem exigir uma implementação específica de runtime.

## Consequências

### Positivas

- separação clara entre projeto e Harness;
- evolução independente dos contratos;
- portabilidade entre providers e agentes;
- validação segura sem execução;
- espaço para extensões futuras sem quebrar leitores compatíveis.

### Negativas

- existirão dois arquivos YAML relacionados ao projeto;
- os blocos opcionais precisarão de specs próprias para evitar ambiguidade;
- preservar campos desconhecidos exige cuidado no parser e na escrita futura;
- o contrato não promete preservação textual de comentários ou ordem dos campos;
- a validação estrutural não comprova disponibilidade ou segurança de recursos.

## Follow-up Actions

- definir os contratos de Provider Profile, Tool Policy e Evaluation Contract em
  specs separadas quando houver necessidade concreta;
- definir a interface de validação e sua integração com `darp doctor` antes de
  adicionar uma subcomanda;
- manter a primeira implementação como biblioteca interna, sem subcomando
  próprio e sem execução;
- definir migração entre versões de `apiVersion` antes de suportar atualização
  automática;
- avaliar provenance, assinaturas e lockfiles em uma spec de supply chain.
