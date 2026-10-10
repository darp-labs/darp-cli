# Tasks 007 — Skill compartilhada para branch, commit e pull request

## Related Plan

- [plan.md](plan.md)

## Task List

- [x] Escrever a skill canônica com fluxo profissional de branch, commit e PR
      via `gh`, incluindo tratamento de estado sujo, validações, template,
      escopo, erros e limites de autorização.
- [x] Espelhar a skill nos diretórios `.github/skills/` e `.claude/skills/` e
      confirmar que as três cópias são idênticas.
- [x] Atualizar `CHANGELOG.md` em `Unreleased`.
- [x] Executar o validador da skill, `git diff --check` e revisar o diff final.

## Validation Checklist

- [x] A skill está disponível nos três ambientes solicitados.
- [x] A skill usa template de PR do repositório quando disponível.
- [x] Validações documentais passam e as cópias permanecem idênticas.

## Notes and Assumptions

- As tasks foram marcadas somente após a validação correspondente passar.
- Não fazer push, abrir PR ou executar outros efeitos remotos ao criar esta
  skill.
- Evidências finais: `/usr/bin/python3
  /home/darp/.codex/skills/.system/skill-creator/scripts/quick_validate.py
  .agents/skills/create-pr` concluiu com `Skill is valid!`; `cmp` confirmou
  igualdade byte a byte das três cópias; `git diff --check --` sobre os três
  caminhos da skill, `CHANGELOG.md` e os artefatos desta spec passou. A
  revisão manual confirmou que os requisitos R1–R9 e os critérios de aceitação
  estão cobertos.
- Revisões do ciclo: arquitetura `PASS` (skill documental independente, sem
  dependências ou acoplamento de execução); documentação `PASS` (três caminhos
  de descoberta e changelog conferidos); testes/validação `PASS` para a skill
  e compatibilidade `PASS` (nenhuma interface de runtime alterada); release
  `PASS` (entrada em `Unreleased`, sem mudança incompatível).
- A invocação global `git diff --check` identifica trailing whitespace em
  `docs/HARNESS-STRATEGY.md`, alteração preexistente e fora do escopo desta
  spec. Ela foi preservada; a verificação dos arquivos da spec 007 passou.
