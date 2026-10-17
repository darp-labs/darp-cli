# Plan 009 — Sincronização explícita do inventário de assets

## Status

Draft — preparado a partir da Spec 009 aprovada. O plano deve ser revisado e
aprovado antes do início da implementação.

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

Estender o pacote de descoberta para reconhecer skills DARP em
`.agents/skills/<nome>/SKILL.md`, mantendo-as distintas de assets externos: as
skills vão para `skills`; os padrões externos da Spec 005 vão para `assets`.
Tratar `.agents/skills/` como raiz opcional e não percorrer subdiretórios
arbitrários abaixo dela.

Propagar erros de listagem, stat e inspeção como itens ignorados com motivo
estável e warning no relatório. A descoberta deve indicar o caminho/escopo que
falhou, para que o reconciliador não interprete uma raiz ou subárvore inacessível
como remoção. Symlinks continuam sem ser seguidos.

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
para falha segura/operacional. A escrita aplicada substitui `darp.yml` de forma
atômica e preserva o original em erro. `darp init`, `darp doctor` e o pacote
Harness não recebem comportamento novo nesta entrega.

## Áreas afetadas

- `internal/project/discovery/` — descoberta das skills DARP, classificação,
  erros/ignorados com escopo e cobertura de testes.
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
- Bloquear remoção no escopo afetado por erro de leitura e continuar com outros
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
  `.agents/skills/`, diretórios ausentes, estruturas inválidas, conflitos,
  ordenação e erros de leitura.
- Unit tests do plano para adoção de registros legados válidos e ausentes,
  preservação de registros desconhecidos, aplicação do ledger e remoção apenas
  de identidades gerenciadas.
- Testes de sync para o caso real: remover skill GitHub, adicionar duas skills
  em `.agents/skills/`, executar sync e validar o inventário e ownership
  resultantes.
- Testes de escopo parcial assegurando que erro de leitura gera warning/ignored
  e não remove itens sob a raiz ou subárvore inacessível.
- Testes de filesystem que comparem bytes e snapshots antes/depois de dry-run,
  falha de parsing, falha de escrita e segunda execução idempotente.
- Testes de CLI para argumentos desconhecidos, saída ordenada e códigos `0`/`1`.
- Testes de regressão confirmando que `darp init` não remove registros, que o
  Doctor continua somente leitura e que Harness não é alterado.
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
- A decisão de adotar entradas legadas reconhecíveis é aplicada somente a
  identidades que correspondam exatamente às regras de discovery aprovadas.
- Entradas fora dos contratos conhecidos e chaves de topo desconhecidas são
  mantidas sem inferir ownership.
- Skills válidas são filhos diretos de `.agents/skills/` com `SKILL.md` regular.
- Avisos por leitura parcial não impedem sync seguro de outras raízes e resultam
  em código `0`.
- A operação não precisa de rede, execução de assets, Harness ou Registry.
- O roadmap permanece inalterado nesta feature.
