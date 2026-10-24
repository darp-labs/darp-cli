# Spec 008 — Harness Lifecycle

## Status

Approved — escopo e critérios aprovados pelo responsável do projeto. Implementação autorizada somente conforme esta especificação, o plano e as tarefas aprovados.

## Contexto

A Spec 006 estabeleceu `.darp/harness.yaml` como manifesto canônico do Harness DARP e implementou parsing e validação estrutural no pacote `internal/project/harness`. O contrato é declarativo: ler ou validar o manifesto não executa agentes, providers, ferramentas, MCP servers, workflows ou quality gates.

A CLI já possui `darp init` e `darp doctor`. O `darp init` inicializa a estrutura geral do projeto; o `darp doctor` executa diagnóstico somente leitura dos contratos DARP existentes. A implementação atual do Doctor possui oito checks, cobrindo configuração, estrutura, workflows, skills, templates, governança, compatibilidade de versão e assets. Não existe integração explícita do Harness Manifest no fluxo atual do Doctor.

O roadmap aponta Harness Lifecycle como próxima fase e menciona inspeção, inicialização, validação, atualização e manutenção. Esta especificação delimita a primeira entrega executável desse ciclo, sem tentar implementar todo o ciclo futuro.

## Problema

O DARP já possui um contrato de Harness legível e validável por biblioteca, mas a CLI ainda não oferece uma interface explícita e consistente para o usuário inspecionar o estado do manifesto, inicializá-lo de forma segura e solicitar sua validação. O diagnóstico geral também não informa claramente o estado do Harness.

## Objetivos

1. Expor operações locais e explícitas para inspecionar, inicializar e validar o Harness.
2. Reutilizar o pacote `internal/project/harness` como fonte única das regras de parsing e validação estrutural.
3. Integrar ao `darp doctor` um diagnóstico do manifesto, sem duplicar a validação e sem alterar o comportamento somente leitura.
4. Definir resultados determinísticos para manifesto ausente, válido, inválido e erro de leitura.
5. Preservar arquivos existentes e evitar efeitos colaterais silenciosos.
6. Documentar comandos, resultados e limites do ciclo de vida implementado.

## Não objetivos

- Atualizar, migrar, reparar ou reformatar automaticamente manifestos existentes.
- Gerar projeções para `AGENTS.md`, Copilot, Claude ou outros agentes.
- Descobrir automaticamente a stack ou ativar recursos encontrados.
- Instalar, remover, executar ou atualizar Skills, Agent Plugins, MCP servers, providers ou workflows.
- Criar Registry, SDK, Runtime ou executor.
- Resolver dependências externas, acessar a rede ou consultar serviços remotos.
- Alterar o contrato da Spec 006 ou redefinir `darp.yml`.
- Tornar a presença do manifesto obrigatória para que `darp init` ou o diagnóstico geral funcionem em projetos que ainda não inicializaram um Harness.
- Criar um novo ADR global; as decisões arquiteturais relevantes já estão registradas nos ADRs existentes.

## Terminologia

- **Manifesto ausente:** `.darp/harness.yaml` não existe.
- **Manifesto válido:** o arquivo existe e não apresenta diagnósticos estruturais segundo a implementação da Spec 006.
- **Manifesto inválido:** o arquivo foi lido, mas o parser/validador retornou diagnósticos.
- **Erro de leitura:** o arquivo não pôde ser lido por motivo diferente de ausência, por exemplo permissão ou erro de I/O.
- **Inicialização:** criação explícita de um manifesto mínimo, sem inferir tecnologias, agentes ou providers.
- **Inspeção:** apresentação do estado do manifesto e de seus metadados conhecidos, sem escrita.

## Requisitos funcionais

### R1 — Comando de Harness

A CLI MUST reconhecer o grupo de comandos `darp harness` e apresentar ajuda específica para os subcomandos suportados. Comandos e argumentos desconhecidos MUST retornar erro de uso não zero, sem modificar arquivos.

### R2 — Inspeção

`darp harness inspect` MUST inspecionar somente o caminho canônico `.darp/harness.yaml`.

- Se o manifesto não existir, informar `MISSING` de forma clara, sem criá-lo.
- Se existir e for válido, informar `VALID` e exibir, no mínimo, `apiVersion`, `kind`, `metadata.name` e `metadata.version`.
- Se existir, mas for inválido, informar `INVALID` e exibir os diagnósticos determinísticos fornecidos pelo pacote Harness.
- Se ocorrer erro de leitura, informar o erro sem classificá-lo como manifesto inválido.
- O comando MUST NOT escrever arquivos, acessar a rede ou inspecionar/executar recursos referenciados.

### R3 — Inicialização explícita

`darp harness init` MUST criar `.darp/harness.yaml` apenas mediante invocação explícita.

- O manifesto inicial deve conter somente os campos mínimos exigidos pela Spec 006: `apiVersion`, `kind`, `metadata.name` e `metadata.version`.
- O nome deve ser derivado deterministicamente do nome do diretório do projeto e normalizado para um identificador válido. A regra de normalização deve ser definida no plano e coberta por testes. Se o nome do diretório não produzir um identificador válido após a normalização, o comando deve falhar com mensagem clara, sem solicitar entrada interativa e sem criar arquivos.
- A versão inicial do Harness deve ser `0.1.0`.
- A inicialização MUST NOT detectar ou presumir stack, agente, provider, tool, workflow ou capability.
- A inicialização MUST NOT sobrescrever ou modificar um manifesto existente.
- Se o manifesto já existir, o comando deve falhar de forma segura e indicar o caminho existente.
- A criação deve ser exclusiva, recusando sobrescrita inclusive em caso de concorrência. A estratégia de escrita deve evitar arquivo parcial; falhas devem ser testadas e não podem modificar arquivos não relacionados.
- O comando pode criar o diretório `.darp/` se estiver ausente, mas não deve alterar outros arquivos.
- O comando MUST NOT executar ou ativar qualquer recurso.

### R4 — Validação explícita

`darp harness validate` MUST validar o manifesto no caminho canônico reutilizando o pacote `internal/project/harness`, sem replicar regras.

- Manifesto válido resulta em saída de sucesso e código `0`.
- Manifesto inválido resulta em diagnósticos úteis e código não zero.
- Manifesto ausente resulta em código não zero e mensagem específica.
- Erros de leitura resultam em código não zero e mensagem específica.
- A validação MUST NOT escrever arquivos, acessar a rede ou executar referências declaradas.

### R5 — Integração com Doctor

O `darp doctor` MUST incluir um check de Harness que reutilize o pacote `internal/project/harness`.

- Manifesto presente e válido: `PASS`.
- Manifesto presente e estruturalmente inválido: `FAIL`, com diagnósticos úteis.
- Manifesto ausente: `WARNING`, pois o Harness ainda pode não ter sido inicializado e sua ausência não deve invalidar os contratos DARP existentes.
- Erro de leitura que impeça diagnóstico confiável: `FAIL`.
- O check MUST NOT duplicar regras de parsing/validação.
- O diagnóstico geral MUST permanecer somente leitura e não deve inicializar ou corrigir o Harness.

### R6 — Consistência de saída e códigos de saída

Os subcomandos devem usar mensagens estáveis, legíveis e adequadas a scripts. `inspect` retorna `0` quando consegue concluir a inspeção dos estados `MISSING`, `VALID` ou `INVALID`; erro operacional de leitura ou uso inválido retorna não zero. `validate` retorna `0` somente para manifesto válido; manifesto ausente, inválido ou erro de leitura retorna não zero. `init` retorna `0` após criação bem-sucedida e não zero se não puder criar com segurança. Todos esses comportamentos devem ser testados.

### R7 — Compatibilidade

A implementação MUST preservar os comandos e contratos existentes de `darp init`, `darp doctor`, `darp help` e `darp version`, exceto pela adição do novo check do Harness ao Doctor e pela inclusão de ajuda para o grupo/subcomandos novos.

O manifesto continua separado de `darp.yml`. Não deve ser criado fallback para outro caminho.

### R8 — Documentação

A documentação de usuário e `CHANGELOG.md` em `Unreleased` MUST ser atualizados para descrever somente os comportamentos implementados. A documentação deve reforçar que Harness é declarativo e que inspeção/validação não ativa recursos.

## Fluxos de usuário

### Inspecionar um projeto sem Harness

1. O usuário executa `darp harness inspect`.
2. A CLI informa que `.darp/harness.yaml` está ausente.
3. Nenhum arquivo é criado ou alterado.

### Inicializar Harness

1. O usuário executa `darp harness init`.
2. A CLI verifica se o caminho canônico já existe.
3. Se não existir, cria um manifesto mínimo e informa o caminho criado.
4. Se existir, falha sem modificar o arquivo existente.

### Validar manifesto

1. O usuário executa `darp harness validate`.
2. A CLI lê o manifesto canônico e reutiliza a validação da Spec 006.
3. A CLI apresenta sucesso ou os diagnósticos, com o código de saída definido.

### Diagnosticar projeto

1. O usuário executa `darp doctor`.
2. O Doctor executa os checks atuais e o novo check de Harness.
3. A ausência do Harness aparece como `WARNING`; um manifesto inválido aparece como `FAIL`.
4. Nenhum arquivo é criado ou corrigido.

## Critérios de aceitação

- [ ] `darp harness` apresenta ajuda e lista os subcomandos suportados.
- [ ] `darp harness inspect` distingue manifesto ausente, válido, inválido e erro de leitura.
- [ ] `darp harness inspect` não altera o estado do projeto.
- [ ] `darp harness init` cria um manifesto mínimo válido usando o contrato da Spec 006.
- [ ] `darp harness init` não sobrescreve manifesto existente.
- [ ] Falha de inicialização não deixa manifesto parcial nem altera arquivos não relacionados.
- [ ] `darp harness validate` reutiliza o pacote Harness e produz códigos de saída definidos.
- [ ] O Doctor reporta Harness ausente como `WARNING`, válido como `PASS`, inválido como `FAIL` e erro de leitura como `FAIL`.
- [ ] Os checks do Doctor continuam somente leitura e os checks preexistentes mantêm seu comportamento.
- [ ] Há testes para os estados de manifesto, erros de I/O, preservação de arquivos e códigos de saída.
- [ ] `make verify` passa.
- [ ] `git diff --check` passa nos arquivos alterados.
- [ ] README/documentação e `CHANGELOG.md` refletem apenas funcionalidades implementadas.

## Riscos e decisões de implementação

1. A normalização de `metadata.name` deve ser determinística, não interativa e testada para nomes comuns e casos-limite. Se o resultado não for válido, falhar sem efeitos colaterais.
2. A inicialização deve usar criação exclusiva e uma estratégia de escrita segura contra concorrência e falhas de persistência; nunca sobrescrever arquivo existente nem deixar manifesto parcial.
3. O Doctor hoje organiza checks como lista ordenada. A implementação deve adicionar o check Harness sem remover ou alterar a semântica dos checks existentes e atualizar testes que dependam da contagem/ordem.
4. O manifesto pode conter campos desconhecidos preservados pelo modelo interno. A Spec 008 não deve reemitir o manifesto nem perder esses campos.
5. A ausência do manifesto é `WARNING`, mantendo compatibilidade com projetos que ainda não adotaram Harness.

## Dependências

- Spec 006 — Harness Manifest.
- `internal/project/harness`.
- `internal/project/doctor`.
- `internal/cli`.
- Constituição, `AGENTS.md` e lifecycle operacional do repositório.

## Referências

- [Spec 006 — Harness Manifest](../006-harness-manifest/spec.md)
- [ADR 001 — Harness Manifest Contract](../006-harness-manifest/adr/001-harness-manifest-contract.md)
- [Roadmap](../../../docs/ROADMAP.md)
- [ADR-002 — Ecosystem Architecture and Agent Interoperability](../../../docs/ADR/ADR-002-darp-ecosystem-architecture-and-agent-interoperability.md)
- [ADR-002 — Cross-Agent Harness Boundaries](../../../docs/ADR/ADR-002-cross-agent-harness-boundaries.md)
