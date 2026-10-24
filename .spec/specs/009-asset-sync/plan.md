# Plan 009 — Sincronização explícita do inventário de assets

## Status

Approved — plano revisado e aprovado em 2026-10-19, alinhado à Spec 009.
Adendo autorizado em 2026-10-23 para validar skills conforme o padrão Agent
Skills publicado em https://agentskills.io/specification.md. Implementação
autorizada somente conforme a Spec, este plano e as tasks aprovadas.
Implementação concluída em 2026-10-23; os gates finais passaram conforme
evidenciado em `tasks.md`.

## Related Specification

- [Spec 009 — Sincronização explícita do inventário de assets](spec.md)
- [Spec 005 — Descoberta e registro de assets existentes](../005-asset-discovery/spec.md)
- [Spec 004 — Scaffold de skills de governança](../004-init-governance-skills/spec.md)
- [Spec 002 — `darp doctor`](../002-doctor/spec.md)
- [Spec 008 — Harness Lifecycle](../008-harness-lifecycle/spec.md), somente
  para compatibilidade da CLI; esta feature não implementa Harness.

## Estratégia de implementação

Manter separadas descoberta, planejamento da reconciliação, persistência e
apresentação da CLI. A descoberta continua somente leitura. Um serviço de sync
constrói um plano determinístico com adições, remoções, itens mantidos,
ignorados e avisos; somente depois de validar todo o plano poderá gravar o
`darp.yml`.

Adicionar uma operação de descoberta específica para `sync` que reconheça
skills DARP em `.agents/skills/<nome>/SKILL.md`, mantendo-as distintas de assets
externos: as skills vão para `skills`; os padrões externos da Spec 005 vão para
`assets`. Não alterar a saída de `discovery.Discover` consumida por `darp init`;
testar que `init` não registra nem modifica skills customizadas em
`.agents/skills/`. Tratar essa raiz como opcional e não percorrer subdiretórios
arbitrários abaixo dela. Validar frontmatter e campos conforme a Spec 009/R13;
reportar skills incompatíveis como ignoradas e preservar ownership dentro do
escopo não validado.

Usar o nome literal do diretório como chave e caminho, sem normalização, após
validar o formato Agent Skills. Ao adotar entradas legadas, exigir um nome que
seja um único componente e o caminho canônico exato; nunca normalizar `..`, `/`
ou outros segmentos. Se a chave já existir em `skills` apontando para outro
caminho, preservar o valor manual, não adotar a entrada como gerenciada e
reportar o conflito.

Propagar erros de listagem, stat e inspeção como itens ignorados com motivo
estável e warning no relatório. A descoberta deve indicar o caminho/escopo que
falhou, para que o reconciliador não interprete uma raiz ou subárvore inacessível
como remoção. Symlinks continuam sem ser seguidos. Erros de leitura, symlinks,
estruturas incompletas, candidatos não regulares e conflitos tornam o escopo
correspondente desconhecido e impedem remoções nesse escopo.

O serviço de sync deve ler o YAML como árvore para preservar campos
desconhecidos, comentários quando possível e dados adicionais de registros
mantidos. Na primeira execução, adotar as entradas legadas que correspondam
exatamente aos contratos de discovery, gerar `darp.sync.managed` e reconciliar
as ausentes somente em escopos lidos com sucesso. Nas execuções seguintes, o
ledger namespaced no próprio `darp.yml` define quais identidades são gerenciadas.
Não inferir ownership de registros fora desses contratos.

Adicionar `internal/project/sync` como serviço de planejamento e aplicação, sem
reutilizar acidentalmente os efeitos de escrita de `init`. Evitar duplicar
validação estrutural do caminho/classificação: compartilhar ou expor somente as
regras necessárias do pacote `discovery`, mantendo o parser de `darp.yml`
responsável por validar o formato editável e sua árvore YAML. Se a implementação
identificar que compartilhar o mecanismo atual de atualização YAML exige
refatoração maior,
registrar a decisão no review arquitetural antes de ampliar o escopo.

Adicionar `darp sync` e `--dry-run` na CLI, com stdout determinístico para o
relatório, código `0` para conclusão inclusive parcial com warnings e código `1`
para falha segura/operacional. Se `darp.yml` estiver ausente, executar o
bootstrap padrão por `internal/project/init` e depois reconciliar o inventário;
se existir, não executar o bootstrap. Com `--dry-run` e sem `darp.yml`, simular
em memória o estado produzido pelo bootstrap e reconciliar sobre esse estado,
sem criar ou alterar arquivos ou diretórios. A escrita aplicada substitui
`darp.yml` de forma atômica e preserva o original em erro. Não alterar o
comportamento de `darp init`, Doctor, help/version ou Harness; não depender da
implementação pendente dos comandos Harness da Spec 008.

## Áreas afetadas

- `internal/project/discovery/` — mecanismo específico de descoberta de skills
  DARP para `sync`, validação Agent Skills, erros/ignorados com escopo e
  cobertura de testes; preservar a saída de `discovery.Discover` usada por
  `init`.
- `internal/project/sync/` — novo plano de alterações, ledger de ownership,
  reconciliação, modo dry-run e escrita segura de `darp.yml`.
- `internal/cli/` — roteamento, ajuda, argumentos, relatório e código de saída
  de `darp sync`.
- `internal/project/asset/` — somente se for necessário compartilhar tipos ou
  identidade de assets; skills DARP não devem ser confundidas com assets
  externos registrados em `assets`.
- testes unitários e de integração de discovery, sync e CLI; testes de
  compatibilidade para `init` e Doctor.
- `README.md` e `CHANGELOG.md` (`Unreleased`).

Não alterar `.darp/harness.yaml`, o comportamento de `darp init`, os arquivos
originais dos assets nem o roadmap nesta etapa.

## Marcos

### M1 — Regras de discovery e estados de leitura

- Implementar discovery de skills diretas em `.agents/skills/` e manter os
  padrões externos da Spec 005.
- Distinguir diretórios opcionais ausentes de erros reais de leitura.
- Reportar symlinks, formatos/estruturas inválidos e erros de I/O como
  ignorados; erros de I/O também produzem warning.
- Tornar observável o escopo de leitura parcial para evitar remoções inseguras.

### M2 — Modelo de plano e ownership

- Definir tipos de plano determinísticos para inclusão, manutenção, remoção,
  ignorados, warnings e falhas.
- Implementar adoção legada por identidade/path conforme Spec 009.
- Persistir o ledger `darp.sync.managed` no `darp.yml`.
- Preservar entradas não gerenciadas, campos desconhecidos e dados adicionais
  de registros mantidos.
- Tratar as quatro skills de governança como itens do espelho, sem proteção
  especial contra remoção.

### M3 — Reconciliação segura do manifesto

- Construir e validar o plano antes de qualquer escrita.
- Bloquear remoção no escopo desconhecido por erro de leitura, symlink,
  estrutura incompleta, candidato não regular ou conflito; continuar com outros
  escopos somente quando o resultado global puder ser persistido com segurança.
- Garantir idempotência byte a byte quando não houver alterações.
- Implementar substituição atômica do único arquivo `darp.yml`, com rollback ou
  preservação do original em falhas.

### M4 — CLI `darp sync`

- Adicionar o comando independente e a opção `--dry-run`.
- Apresentar resultados determinísticos de adicionados, removidos, mantidos,
  ignorados e warnings.
- Aplicar códigos `0` para conclusão inclusive parcial e `1` para falha.
- Manter `darp init`, `darp doctor` e Harness compatíveis e sem novos efeitos.

### M5 — Testes, documentação e qualidade

- Cobrir migração inicial, ownership, sync parcial, preservação, falhas,
  idempotência, dry-run e comportamento da CLI.
- Atualizar README e `CHANGELOG.md` em `Unreleased` após implementação.
- Revisar arquitetura, compatibilidade e documentação; executar `make verify` e
  `git diff --check`.

## Estratégia de validação

- Unit tests de descoberta para cada padrão externo vigente, skills em
  `.agents/skills/`, frontmatter Agent Skills válido/inválido, diretórios
  ausentes, estruturas incompletas, conflitos, ordenação e erros de leitura.
- Teste de regressão de `init` provando que skills DARP customizadas em
  `.agents/skills/` não são adicionadas a `darp.yml` nem modificadas.
- Teste de integração de `sync` sem `darp.yml`, provando que o bootstrap padrão
  é criado e o inventário é reconciliado em seguida.
- Teste de integração de `sync --dry-run` sem `darp.yml`, comparando o plano ao
  da execução aplicada e verificando que nenhum arquivo ou diretório do
  bootstrap é criado ou alterado.
- Unit tests do plano para adoção de registros legados válidos e ausentes,
  preservação de registros desconhecidos, aplicação do ledger e remoção apenas
  de identidades gerenciadas.
- Testes de sync para o caso real: remover skill GitHub, adicionar duas skills
  em `.agents/skills/`, executar sync e validar o inventário e ownership
  resultantes.
- Testes de escopo parcial assegurando que erro de leitura, symlink, estrutura
  incompleta, arquivo não regular ou conflito gera warning/ignored e impede
  remoções das identidades no escopo não observado.
- Testes de filesystem que comparem bytes e snapshots antes/depois de dry-run,
  falha de parsing, falha de escrita e segunda execução idempotente.
- Testes de CLI para argumentos desconhecidos, saída ordenada e códigos `0`/`1`.
- Testes de regressão confirmando que `darp init` não remove registros nem
  adota skills DARP customizadas, que o Doctor continua somente leitura e que
  Harness não é alterado nem é dependência desta feature.
- Executar testes focados e `go test ./...`; em seguida `make verify` e
  `git diff --check`.

## Recuperação ou rollback

Uma falha de parsing ou validação não pode gravar nada. Uma falha durante a
persistência deve preservar o `darp.yml` anterior byte a byte e remover arquivos
temporários criados pela operação. Como o ledger de ownership fica no mesmo
manifesto, inventário e ledger devem ser gravados na mesma substituição atômica.

O rollback da feature consiste em remover o comando, o serviço e a integração
nova de discovery. O rollback não deve apagar entradas já sincronizadas nem
remover `darp.sync.managed` automaticamente; qualquer migração reversa futura
exigiria decisão explícita.

## Premissas

- A Spec 009 e este plano serão aprovados antes de implementação.
- Em `--dry-run` sem `darp.yml`, o estado inicial do bootstrap é montado em
  memória; nenhuma operação de escrita do `init` é chamada.
- A decisão de adotar entradas legadas reconhecíveis é aplicada somente a
  identidades que correspondam exatamente às regras de discovery aprovadas.
- Entradas fora dos contratos conhecidos e chaves de topo desconhecidas são
  mantidas sem inferir ownership.
- Skills válidas são filhos diretos de `.agents/skills/` com `SKILL.md` regular.
- Avisos por leitura parcial não impedem sync seguro de outras raízes e resultam
  em código `0`.
- A operação não precisa de rede, execução de assets, Harness ou Registry.
- Skills válidas usam o nome literal do diretório, sem normalização; conflito
  com chave já existente em `skills` preserva o valor manual e gera relatório.
- Se `darp.yml` estiver ausente, `sync` chama o bootstrap padrão de `init` antes
  da reconciliação. Isso não altera o comportamento de invocações de `init`.
- O roadmap permanece inalterado nesta feature.
