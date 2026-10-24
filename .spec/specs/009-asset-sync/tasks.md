# Tasks 009 — Sincronização explícita do inventário de assets

## Status

Approved — tasks revisadas e aprovadas em 2026-10-19, alinhadas à Spec 009 e
ao Plan 009. Adendo autorizado em 2026-10-23 para validação conforme o padrão
Agent Skills. Cada task só é marcada como concluída após sua validação
correspondente. Implementação concluída em 2026-10-23 após a validação final.

## Related Plan

- [Plan 009 — Sincronização explícita do inventário de assets](plan.md)
- [Spec 009 — Sincronização explícita do inventário de assets](spec.md)

## Regras para o implementador

- Não iniciar implementação antes da aprovação formal deste plano e destas
  tasks.
- Seguir a Spec 009; não alterar retroativamente o contrato de `darp init`.
- Manter `darp init` aditivo e preservar os registros existentes quando for
  reexecutado.
- `darp sync` é a única operação nova que reconcilia `skills`, `assets` e o
  ledger de ownership.
- Não remover entrada fora dos padrões reconhecidos ou fora do ledger após a
  migração inicial especificada.
- Erros de leitura geram warning/ignored e impedem remoções no escopo afetado.
- Não alterar os arquivos dos assets, executar seu conteúdo, acessar rede ou
  modificar `.darp/harness.yaml`.
- Manter todas as tasks desmarcadas até executar e registrar a validação
  correspondente. Tarefas bloqueadas permanecem desmarcadas e registram motivo
  e evidência.

## Task List

- [x] **T1 — Confirmar prontidão SDD.** Confirmar aprovação da Spec 009, Plan
      009 e Tasks 009; registrar quaisquer divergências antes de editar código.
- [x] **T2 — Adicionar discovery de skills exclusivo para sync.** Reconhecer
      filhos diretos de `.agents/skills/<nome>/SKILL.md`, preservar o nome
  literal sem normalização, manter separados os resultados que alimentam
  `skills` e `assets`, e preservar a saída de `discovery.Discover` usada por
  `init`. Testar padrões, diretórios ausentes, arquivos inválidos, conflitos,
  symlinks e ordenação determinística. O adendo Agent Skills é rastreado em
  T15.
- [x] **T3 — Propagar falhas parciais de leitura.** Distinguir ausência opcional
      de erro de leitura; incluir caminho/escopo e motivo em itens ignorados;
      produzir warnings e testar que erros, symlinks, estruturas incompletas,
      candidatos não regulares e conflitos tornam o escopo desconhecido e não
      permitem remoções nele.
- [x] **T4 — Implementar modelo do plano de sync.** Criar resultados
      determinísticos de adicionados, removidos, mantidos, ignorados, warnings e
      falhas, sem efeitos de filesystem e com identidade coerente com as Specs
      004/005/009.
- [x] **T5 — Implementar ownership e migração legada.** Ler/escrever
      `darp.sync.managed` no `darp.yml`; adotar apenas entradas legadas que
      correspondam exatamente às regras de discovery; preservar entradas fora
      do contrato. Usar o nome literal das skills; se uma chave existir com
      outro caminho, manter seu valor, não adotá-la e reportar conflito. Testar
      skills e assets legados, inclusive registros já ausentes.
- [x] **T6 — Implementar reconciliação preservadora.** Planejar inclusão,
      manutenção e remoção; preservar campos desconhecidos e dados adicionais
      de registros mantidos; não remover itens sob escopos com erro de leitura;
      testar conflitos, YAML inválido, chaves duplicadas e idempotência byte a
      byte.
- [x] **T7 — Implementar escrita segura.** Aplicar o manifesto e ledger numa
      substituição atômica única; testar falha de escrita, limpeza de temporário
      e preservação byte a byte do original.
- [x] **T8 — Adicionar CLI `darp sync`.** Implementar roteamento, ajuda,
      argumentos desconhecidos, `--dry-run`, relatório ordenado e códigos de
      saída `0` para conclusão inclusive parcial e `1` para falhas. Se
      `darp.yml` estiver ausente, executar o bootstrap padrão de `init` e depois
      reconciliar; caso contrário, não executar bootstrap. Em `--dry-run` sem
      `darp.yml`, simular o resultado do bootstrap em memória e reconciliar
      sobre ele, sem criar nem alterar arquivos ou diretórios.
- [x] **T9 — Cobrir regressão e compatibilidade.** Testar a remoção da skill
      `.github` e inclusão de duas skills `.agents`; garantir que `init` mantém
      seu comportamento aditivo e não adota/modifica skills DARP customizadas,
      `doctor` permanece somente leitura e Harness não é alterado nem exigido
      como dependência.
- [x] **T10 — Atualizar documentação e changelog.** Descrever `darp sync`,
      ownership, migração legada, dry-run, warnings, códigos de saída e limites;
      atualizar `CHANGELOG.md` em `Unreleased` somente com o comportamento
      implementado.
- [x] **T11 — Executar testes e quality gates.** Executar testes focados,
      `go test ./...`, `make verify` e `git diff --check`; registrar resultados
      e corrigir falhas dentro do escopo aprovado.
- [x] **T12 — Revisão de arquitetura.** Revisar coesão entre discovery, sync,
      persistência e CLI; confirmar que ownership não duplica lógica e que a
      escrita é segura; registrar findings e bloqueios.
- [x] **T13 — Revisão de documentação e compatibilidade.** Conferir exemplos,
      comandos, mensagens, links e alinhamento com Specs 002, 004, 005, 008 e
      009; não atualizar o roadmap nesta entrega.
- [x] **T14 — Revisão final e conclusão.** Confirmar evidências de todas as
      tasks e gates, atualizar status somente após validação, registrar decisão
      de review e release conforme o lifecycle do projeto.
- [x] **T15 — Validar formato Agent Skills.** Validar frontmatter YAML, campos
  obrigatórios e opcionais, sintaxe e correspondência de `name` com o
  diretório. Reportar skills fora do padrão como ignoradas, impedir remoção
  no escopo não validado e cobrir regras com testes. Atualizar README e
  CHANGELOG apenas conforme o comportamento aprovado.

## Validation Checklist

- [x] Spec 009, Plan 009 e Tasks 009 estão aprovados antes da implementação.
- [x] Skills em `.agents/skills/` e assets externos são descobertos e escritos
      nos campos corretos.
- [x] Ownership legado é migrado conforme as regras aprovadas e persistido em
      `darp.sync.managed`.
- [x] Inclusões, remoções, itens mantidos, ignorados e warnings são reportados
      de forma determinística.
- [x] Erros de leitura não causam remoção no escopo não observado.
- [x] Symlinks, estruturas incompletas, candidatos não regulares e conflitos
      também bloqueiam remoções no escopo correspondente.
- [x] Nomes de skills são preservados literalmente sem normalização; uma chave
      manual conflitante é mantida e reportada.
- [x] Skills registradas cumprem o formato Agent Skills, inclusive frontmatter,
  campos obrigatórios, limites aplicáveis e correspondência de nome.
- [x] Skills fora do padrão são reportadas como ignoradas e não causam remoção
  de entradas gerenciadas no escopo não validado.
- [x] Sem `darp.yml`, `sync` executa o bootstrap padrão de `init` e depois
      reconcilia o inventário.
- [x] `--dry-run` não modifica nenhum arquivo.
- [x] `sync --dry-run` sem `darp.yml` simula bootstrap e reconciliação em
      memória, exibe o plano equivalente ao modo aplicado e não cria nem altera
      qualquer arquivo ou diretório que o bootstrap produziria.
- [x] Manifesto inválido e falhas de escrita preservam o arquivo original.
- [x] Segunda sincronização sem mudanças não altera os bytes de `darp.yml`.
- [x] Nenhum arquivo de skill ou asset é alterado ou executado.
- [x] `darp init`, `darp doctor`, help e version mantêm seus contratos; Harness
      não é alterado nem é uma dependência dos comandos ainda não implementados
      da Spec 008.
- [x] Testes focados, `go test ./...`, `make verify` e `git diff --check` passam.
- [x] Revisões de arquitetura, documentação, compatibilidade e qualidade estão
      registradas com evidências.

## Notas e premissas

## Evidências desta implementação

- T1: Spec, Plan e Tasks declaram aprovação formal; as decisões locais
  complementares foram consideradas na implementação.
- T10: `README.md` e `CHANGELOG.md` descrevem o comando e seu comportamento.
- T11: `GOCACHE=/tmp/darp-go-cache make verify` passou (gofmt, `go test ./...`
  e go vet); `git diff --check` passou.
- T12: Architecture Review — PASS. Discovery de sync, reconciliação/persistência
  e CLI têm responsabilidades separadas; `init` continua usando o discovery
  legado. A persistência substitui um único manifesto por rename atômico.
- T13: Documentation/compatibility review — PASS. O README lista o comando e
  documenta dry-run, ownership, bootstrap e códigos de saída. `init`, Doctor e
  Harness não receberam alterações de runtime.
- T3, T4, T9 e T14 foram concluídas em 2026-10-23; ver evidências abaixo.
- T15 foi autorizada como adendo à Spec em 2026-10-23 e concluída após a
  implementação e validação do contrato Agent Skills.
- T15: `go test ./internal/project/sync ./internal/cli` e `go test ./...`
  passaram; `go vet ./...`, `git diff --check` e diagnósticos do VS Code também
  passaram. Cobertura adicionada para frontmatter ausente/malformado, nome
  incompatível, limites/tipos opcionais, skill externa incompleta, caminho
  legado não canônico e preservação de ownership em skill inválida.
- T3: testes de I/O injetado cobrem `ReadDir`, `ReadFile`, `Lstat` e erro de
  `WalkDir`; candidatos não regulares, symlinks, estruturas incompletas e
  conflitos semânticos são reportados e não permitem remoções no escopo
  desconhecido. Escopos legíveis continuam sendo reconciliados.
- T4: `BuildPlan` é somente leitura e determinístico; `Apply` é separado,
  revalida o manifesto contra mudanças posteriores e devolve falhas estruturadas
  no resultado. Testes cobrem estabilidade, falha de manifesto e plano obsoleto.
- T9: teste integrado remove skill GitHub ausente e adiciona duas skills DARP;
  testa bootstrap aplicado sem manifesto, preservação byte a byte das skills e
  de `.darp/harness.yaml`, warnings com exit code `0` e falha de manifesto com
  exit code `1`. `init` aditivo e Doctor somente leitura também têm regressões.
- T14: Architecture Review — PASS; responsabilidades de planejamento,
  discovery, aplicação e CLI estão separadas, sem dependência de Harness.
  Documentation/Compatibility Review — PASS; README e changelog correspondem
  ao comportamento validado e o roadmap ficou inalterado. Release Review — PASS
  para `Unreleased`; nenhuma versão ou publicação foi solicitada. Review final —
  PASS, sem bloqueios conhecidos.
- Gate final em 2026-10-23: `GOCACHE=/tmp/darp-go-cache make verify` passou
  (gofmt, `go test ./...`, lint com 0 issues); `git diff --check` passou.

- A primeira sincronização pode remover registros legados ausentes se sua
  identidade corresponder exatamente aos padrões reconhecidos. `--dry-run`
  deve tornar essa migração revisável antes da gravação.
- O ledger armazena ownership, não histórico nem cópia paralela do inventário.
- As quatro skills de governança não recebem exceção de proteção contra
  remoção; o sync espelha o filesystem.
- A descoberta ampliada de `.agents/skills/` pertence ao fluxo de `sync`; o
  `init` mantém a descoberta externa e o contrato aditivo atual.
- Nomes de skills são literais. Colisões com valores manuais são preservadas e
  reportadas, sem sobrescrita ou adoção de ownership.
- Se `darp.yml` estiver ausente, `sync` executa o bootstrap padrão de `init`
  antes de reconciliar.
- Se `darp.yml` estiver ausente em `--dry-run`, simular o estado do bootstrap
  em memória; não chamar operações de escrita de `init` nem alterar o
  filesystem.
- O roadmap permanece sem alteração conforme decisão aprovada.
- A Spec 008 está aprovada, mas sua implementação não está concluída no estado
  atual; esta feature não implementa nem completa comandos Harness.
