# Plan 007 — Skill compartilhada para branch, commit e pull request

## Related Specification

- [spec.md](spec.md)

## Estratégia de implementação

Criar uma skill autônoma em Markdown com instruções portáveis, focadas em
decisões relevantes do fluxo. Usar `.agents/skills/create-pr/SKILL.md` como
fonte canônica e espelhar o conteúdo nos diretórios de descoberta de Copilot e
Claude Code. Atualizar o changelog e validar frontmatter, nomes e igualdade das
cópias.

## Áreas afetadas

- `.agents/skills/create-pr/SKILL.md`
- `.github/skills/create-pr/SKILL.md`
- `.claude/skills/create-pr/SKILL.md`
- `.spec/specs/007-create-pr-skill/`
- `CHANGELOG.md`

## Marcos

1. Definir o procedimento de branch, commit e PR com salvaguardas de escopo.
2. Tornar a skill disponível nos três ambientes alvo.
3. Atualizar changelog e validar conteúdo e formato.

## Estratégia de validação

- Executar `quick_validate.py` na cópia canônica.
- Comparar as três cópias byte a byte.
- Executar `git diff --check` e revisar o diff.

## Recuperação

Como a mudança é documental, reverter apenas os arquivos criados ou alterados
nesta especificação. Não modificar arquivos de usuário não relacionados.

## Premissas

- `.agents/skills/` é o caminho canônico compartilhado de skills no projeto.
- `.github/skills/` e `.claude/skills/` são mantidos como cópias para descoberta
  nativa, mantendo o formato Agent Skills idêntico.
