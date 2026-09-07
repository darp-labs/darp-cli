# Spec 006 — Harness Manifest

## Status

Approved — implementação limitada ao contrato declarativo e à validação
estrutural interna.

## Contexto

O DARP CLI já inicializa e diagnostica projetos DARP, distribui skills de
governança e registra assets existentes. Os documentos estratégicos do projeto
definem Harness como a composição declarativa de assets, providers,
ferramentas, políticas, workflows e avaliação para engenharia assistida por
agentes.

Ainda não existe um contrato formal para representar essa composição. O
exemplo em `docs/AGENTIC_HARNESS.md` é explicitamente ilustrativo e não deve
ser tratado como schema implementável.

## Problema

Sem um manifesto versionável e provider-agnostic, futuras funcionalidades de
descoberta, validação, integração e registry não possuem uma interface comum.
Implementar execução antes de definir esse contrato criaria acoplamento
prematuro a um agente, provider ou runtime.

## Objetivos

1. Definir um formato declarativo e versionável para um Harness pertencente ao
   repositório.
2. Descrever a composição de assets, providers, ferramentas, políticas,
   workflows, quality gates, avaliação e compatibilidade.
3. Permitir validação estrutural determinística sem executar agentes, modelos,
   ferramentas ou workflows.
4. Preservar extensões futuras e campos desconhecidos sem quebrar leitores
   compatíveis.
5. Estabelecer uma base para futuras specs de Provider Profile, Tool Policy,
   Evaluation Contract, descoberta e registry.

O DARP será o dono do contrato, da validação e da composição declarativa. Não
será o runtime responsável por executar agentes, providers, MCP servers ou
workflows. Um runtime separado poderá consumir o manifesto no futuro.

## Não objetivos

- implementar `darp harness run` ou qualquer executor;
- executar inferência, MCP servers, comandos, workflows ou quality gates;
- detectar automaticamente a stack do repositório;
- integrar um provider, agente ou runtime específico;
- resolver dependências ou instalar assets;
- criar registry, marketplace, cache ou sincronização externa;
- definir o conteúdo detalhado de Provider Profile, Tool Policy ou Evaluation
  Contract além das referências necessárias no manifesto;
- executar ou orquestrar o Harness, mesmo que por meio de uma subcomanda;
- substituir ou alterar o contrato atual de `darp.yml`.

## Contrato do manifesto

### Localização

O manifesto canônico do Harness reside em:

```text
.darp/harness.yaml
```

O arquivo é parte do repositório e deve ser adequado a versionamento em Git.
Não há fallback implícito para arquivos fora desse caminho nesta versão.

### Estrutura mínima

Um manifesto válido possui um mapping YAML com estes campos obrigatórios:

```yaml
apiVersion: darp.dev/v1alpha1
kind: Harness

metadata:
  name: example-harness
  version: 0.1.0
```

`apiVersion` e `kind` identificam o contrato. `metadata.name` identifica o
Harness dentro do projeto e `metadata.version` identifica sua versão
publicável. Os quatro valores são obrigatórios, não vazios e devem ser strings.

O manifesto pode declarar os seguintes blocos opcionais:

```yaml
assets: {}
providers: {}
tools: {}
policy: {}
workflow: {}
qualityGates: {}
evaluation: {}
compatibility: {}
```

O conteúdo desses blocos deve permanecer declarativo. O contrato não atribui
comportamento de execução aos valores declarados.

### Referências e agnosticidade

- Assets são referenciados por identidade ou caminho conforme o contrato de
  assets vigente; o manifesto não copia, instala ou modifica assets.
- Providers e modelos são identificadores declarativos. A presença de um
  provider preferencial não torna outros providers obrigatórios nem vincula o
  Harness a um fornecedor.
- Ferramentas, MCP servers, permissões e limites devem ser expressos como
  intenção e política, não como autorização automática para execução.
- Workflows, quality gates e avaliações são referências ou requisitos
  verificáveis; sua execução será definida por specs posteriores.
- Compatibilidade pode declarar agentes, providers, runtimes ou stacks
  suportados, sem exigir que o DARP execute detecção automática.

### Extensibilidade e preservação

Leitores devem ignorar campos desconhecidos quando puderem preservar a
semântica conhecida. A representação interna deve preservar os nós
desconhecidos para leitura e eventual reemissão. Preservação textual exata de
comentários e ordem dos campos não é requisito desta versão.

Campos reservados para futuras extensões não podem alterar o significado dos
campos definidos nesta spec. Extensões incompatíveis devem usar uma nova versão
de `apiVersion`.

## Validação

Uma validação estrutural deve:

- rejeitar YAML inválido;
- exigir um mapping na raiz;
- exigir `apiVersion: darp.dev/v1alpha1`;
- exigir `kind: Harness`;
- exigir `metadata.name` e `metadata.version` como strings não vazias;
- rejeitar tipos incompatíveis nos campos conhecidos;
- rejeitar caminhos de asset absolutos, com traversal ou fora do projeto;
- rejeitar referências de asset sem identidade suficiente;
- aceitar campos desconhecidos para forward compatibility;
- produzir diagnósticos determinísticos, sem executar conteúdo referenciado.

Validação de disponibilidade de provider, ferramenta, MCP server, modelo,
comando, runtime ou stack está fora desta spec. Ausência desses recursos não
deve ser inferida somente pela leitura do manifesto.

Uma implementação futura deverá distinguir `PASS`, `WARNING` e `FAIL` e usar
saída não zero somente para erro estrutural ou incompatibilidade de contrato.
A definição do comando e de seus códigos de saída pertence a uma spec de
validação posterior ou à evolução explicitamente aprovada do `darp doctor`.

## Fluxos

### Manifesto válido

1. O projeto contém `.darp/harness.yaml` com a estrutura mínima.
2. Um leitor carrega o manifesto sem executar referências.
3. A validação informa sucesso determinístico.

### Manifesto inválido

1. O arquivo contém YAML inválido, raiz incompatível, campos obrigatórios
   ausentes ou tipos inválidos.
2. A validação informa o caminho do problema e o tipo de erro.
3. Nenhum arquivo ou recurso externo é alterado ou executado.

### Extensão desconhecida

1. O manifesto contém campos adicionais não conhecidos pelo leitor.
2. Os campos conhecidos continuam sendo validados.
3. A extensão é preservada e não é interpretada como comportamento suportado.

## Critérios de aceitação

- [x] O contrato define `.darp/harness.yaml` como localização canônica.
- [x] O contrato define `apiVersion`, `kind` e `metadata` mínimos.
- [x] A composição declarativa cobre assets, providers, ferramentas, políticas,
      workflows, quality gates, avaliação e compatibilidade.
- [x] O contrato é provider-agnostic e não exige execução de modelos ou agentes.
- [x] Regras de versionamento e compatibilidade futura estão explícitas.
- [x] Campos desconhecidos podem ser preservados sem alterar campos conhecidos.
- [x] Validações obrigatórias, tipos inválidos, referências inseguras e erros
      determinísticos estão definidos.
- [x] A spec não introduz `harness run`, detecção de stack, registry ou
      integrações concretas.
- [x] O contrato não conflita com `darp.yml`, `darp init`, `darp doctor` ou o
      registro de assets existente.
- [x] A spec, o plano, as tasks e o ADR estão consistentes entre si.

## Limite de responsabilidade

A implementação desta spec expõe somente parsing, modelo e validação
estrutural por biblioteca interna. Não haverá subcomando de execução nem
integração com provider, modelo, MCP server, runtime, workflow ou sistema de
avaliação.

## Riscos e questões futuras

- O formato detalhado dos blocos opcionais deve evoluir em specs próprias para
  evitar um manifesto monolítico.
- A validação sem execução não confirma disponibilidade nem segurança real de
  recursos externos; isso exige contratos e gates posteriores.
- O suporte a migrações entre versões de `apiVersion` deverá ser especificado
  antes de qualquer atualização automática.
