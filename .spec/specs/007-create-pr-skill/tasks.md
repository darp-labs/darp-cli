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
- [ ] Executar o validador da skill, `git diff --check` e revisar o diff final.

## Validation Checklist

- [x] A skill está disponível nos três ambientes solicitados.
- [x] A skill usa template de PR do repositório quando disponível.
- [ ] Validações documentais passam e as cópias permanecem idênticas.

## Notes and Assumptions

- Manter todas as tasks desmarcadas até a validação correspondente passar.
- Não fazer push, abrir PR ou executar outros efeitos remotos ao criar esta
  skill.
- Bloqueio de validação: `quick_validate.py` não executou porque o Python
  disponível não tem o módulo `yaml` (`ModuleNotFoundError: No module named
  'yaml'`). A verificação manual do frontmatter/instruções, a comparação das
  três cópias e `git diff --check` passaram; a task de validação permanece
  desmarcada até o validador poder ser executado.
