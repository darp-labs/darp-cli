# Spec 007 — Skill compartilhada para branch, commit e pull request

## Status

Completed — implementation and validation complete.

## Contexto

O repositório já distribui skills agnósticas em `.agents/skills/`. A criação de
branches, commits e pull requests pelo GitHub CLI exige orientações
reutilizáveis, que também possam ser descobertas pelo GitHub Copilot e pelo
Claude Code.

## Problema

Sem um fluxo comum, agentes podem criar branches e commits pouco claros, incluir
alterações alheias ao pedido ou ignorar o template de pull request do repositório.

## Objetivos

1. Criar uma skill reutilizável que conduza a criação de branch, commit e PR
   usando `gh`.
2. Orientar práticas profissionais de Git, revisão e comunicação do PR sem
   depender de linguagem, framework ou provedor de IA.
3. Usar e preencher o template de PR existente no repositório, quando houver.
4. Tornar a skill descobrível no Codex, GitHub Copilot e Claude Code.

## Não objetivos

- Automatizar publicação, merge, aprovação ou release do PR.
- Impor um gerenciador de dependências, conjunto de testes ou convenção local
  não estabelecida pelo repositório.
- Fazer push ou abrir PR sem pedido explícito do usuário para esse fluxo.

## Requisitos

- R1. A skill deve seguir o formato aberto `SKILL.md`, com frontmatter `name` e
  `description`.
- R2. O fluxo deve inspecionar estado Git, convenções locais, mudanças e
  validações aplicáveis antes de preparar um commit.
- R3. O agente deve isolar somente alterações pertencentes ao pedido, revisar
  o diff staged e evitar incluir segredos ou arquivos alheios.
- R4. Branches e commits devem ter nomes claros e profissionais; commits devem
  usar convenção existente e, na ausência dela, mensagem imperativa e concisa.
- R5. Push e abertura do PR devem usar `gh`; nunca usar push forçado.
- R6. Antes de criar o PR, procurar templates suportados pelo GitHub no
  repositório e preencher o template escolhido sem deixar comentários ou
  placeholders de instrução.
- R7. Se não houver template, compor um corpo conciso com contexto, mudanças,
  validação e riscos relevantes.
- R8. Criar cópias idênticas da skill nos caminhos de descoberta dos três
  ambientes: `.agents/skills/`, `.github/skills/` e `.claude/skills/`.
- R9. Não adicionar dependências nem executar ações além do fluxo solicitado.

## Critérios de aceitação

- [x] Os três caminhos contêm um `SKILL.md` idêntico e válido.
- [x] As instruções incluem branch, commit, push, PR via `gh` e tratamento do
      template existente.
- [x] As instruções preservam alterações existentes e não promovem merge.
- [x] `quick_validate.py` valida a skill canônica e `git diff --check` passa
      nos arquivos desta especificação.
- [x] O changelog registra a skill como mudança de governança disponível aos
      contribuidores.

## Riscos e questões em aberto

- O conteúdo duplicado pode divergir no futuro; registrar que as três cópias
  devem permanecer idênticas e comparar seu conteúdo na validação.
- O local de descoberta pode variar por versão e instalação. Manter os três
  caminhos de projeto usados pelos ambientes alvo.
