# Tasks 009 — Sincronização explícita do inventário de assets

## Status

Draft — tarefas derivadas da Spec 009 e do Plan 009. Todas permanecem
desmarcadas até que a tarefa e sua validação correspondente sejam concluídas.

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

- [ ] **T1 — Confirmar prontidão SDD.** Confirmar aprovação da Spec 009, Plan
      009 e Tasks 009; registrar quaisquer divergências antes de editar código.
- [ ] **T2 — Estender discovery para skills DARP.** Reconhecer apenas filhos
      diretos válidos de `.agents/skills/<nome>/SKILL.md`, manter separados os
      resultados que alimentam `skills` e `assets`, e testar padrões, múltiplas
      skills, diretórios ausentes, arquivos inválidos, conflitos, symlinks e
      ordenação determinística.
- [ ] **T3 — Propagar falhas parciais de leitura.** Distinguir ausência opcional
      de erro de leitura; incluir caminho/escopo e motivo em itens ignorados;
      produzir warnings e testar que escopos inacessíveis não são tratados como
      remoções.
- [ ] **T4 — Implementar modelo do plano de sync.** Criar resultados
      determinísticos de adicionados, removidos, mantidos, ignorados, warnings e
      falhas, sem efeitos de filesystem e com identidade coerente com as Specs
      004/005/009.
- [ ] **T5 — Implementar ownership e migração legada.** Ler/escrever
      `darp.sync.managed` no `darp.yml`; adotar apenas entradas legadas que
      correspondam exatamente às regras de discovery; preservar entradas fora
      do contrato e testar skills e assets legados, inclusive registros já
      ausentes.
- [ ] **T6 — Implementar reconciliação preservadora.** Planejar inclusão,
      manutenção e remoção; preservar campos desconhecidos e dados adicionais
      de registros mantidos; não remover itens sob escopos com erro de leitura;
      testar conflitos, YAML inválido, chaves duplicadas e idempotência byte a
      byte.
- [ ] **T7 — Implementar escrita segura.** Aplicar o manifesto e ledger numa
      substituição atômica única; testar falha de escrita, limpeza de temporário
      e preservação byte a byte do original.
- [ ] **T8 — Adicionar CLI `darp sync`.** Implementar roteamento, ajuda,
      argumentos desconhecidos, `--dry-run`, relatório ordenado e códigos de
      saída `0` para conclusão inclusive parcial e `1` para falhas.
- [ ] **T9 — Cobrir regressão e compatibilidade.** Testar a remoção da skill
      `.github` e inclusão de duas skills `.agents`; garantir que `init` mantém
      seu comportamento aditivo, `doctor` permanece somente leitura e os
      comandos Harness não são alterados.
- [ ] **T10 — Atualizar documentação e changelog.** Descrever `darp sync`,
      ownership, migração legada, dry-run, warnings, códigos de saída e limites;
      atualizar `CHANGELOG.md` em `Unreleased` somente com o comportamento
      implementado.
- [ ] **T11 — Executar testes e quality gates.** Executar testes focados,
      `go test ./...`, `make verify` e `git diff --check`; registrar resultados
      e corrigir falhas dentro do escopo aprovado.
- [ ] **T12 — Revisão de arquitetura.** Revisar coesão entre discovery, sync,
      persistência e CLI; confirmar que ownership não duplica lógica e que a
      escrita é segura; registrar findings e bloqueios.
- [ ] **T13 — Revisão de documentação e compatibilidade.** Conferir exemplos,
      comandos, mensagens, links e alinhamento com Specs 002, 004, 005, 008 e
      009; não atualizar o roadmap nesta entrega.
- [ ] **T14 — Revisão final e conclusão.** Confirmar evidências de todas as
      tasks e gates, atualizar status somente após validação, registrar decisão
      de review e release conforme o lifecycle do projeto.

## Validation Checklist

- [ ] Spec 009, Plan 009 e Tasks 009 estão aprovados antes da implementação.
- [ ] Skills em `.agents/skills/` e assets externos são descobertos e escritos
      nos campos corretos.
- [ ] Ownership legado é migrado conforme as regras aprovadas e persistido em
      `darp.sync.managed`.
- [ ] Inclusões, remoções, itens mantidos, ignorados e warnings são reportados
      de forma determinística.
- [ ] Erros de leitura não causam remoção no escopo não observado.
- [ ] `--dry-run` não modifica nenhum arquivo.
- [ ] Manifesto inválido e falhas de escrita preservam o arquivo original.
- [ ] Segunda sincronização sem mudanças não altera os bytes de `darp.yml`.
- [ ] Nenhum arquivo de skill ou asset é alterado ou executado.
- [ ] `darp init`, `darp doctor` e Harness mantêm seus contratos.
- [ ] Testes focados, `go test ./...`, `make verify` e `git diff --check` passam.
- [ ] Revisões de arquitetura, documentação, compatibilidade e qualidade estão
      registradas com evidências.

## Notas e premissas

- A primeira sincronização pode remover registros legados ausentes se sua
  identidade corresponder exatamente aos padrões reconhecidos. `--dry-run`
  deve tornar essa migração revisável antes da gravação.
- O ledger armazena ownership, não histórico nem cópia paralela do inventário.
- As quatro skills de governança não recebem exceção de proteção contra
  remoção; o sync espelha o filesystem.
- O roadmap permanece sem alteração conforme decisão aprovada.
- A Spec 008 está aprovada, mas sua implementação não está concluída no estado
  atual; esta feature não implementa nem completa comandos Harness.
