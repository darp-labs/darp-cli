# Changelog

Todas as mudanças relevantes do DARP CLI serão documentadas neste arquivo.

O formato segue [Keep a Changelog](https://keepachangelog.com/pt-BR/1.1.0/), e
as versões seguirão [Semantic Versioning](https://semver.org/lang/pt-BR/).

## [Unreleased]

### Adicionado

- Comando `darp init` para inicializar projetos DARP com os contratos, a
  configuração base e a estrutura compartilhada de skills.
- Comando `darp doctor` para diagnosticar projetos DARP sem alterar arquivos,
  incluindo configuração, estrutura, workflows, skills, templates,
  governança e compatibilidade de contratos.
- Comandos `darp --help` e `darp --version`.
- Fluxo de desenvolvimento orientado por especificações, com ciclo de vida,
  governança e documentação para contribuições.

### Alterado

- `darp init` restaura apenas arquivos e diretórios ausentes em projetos
  parciais, preservando arquivos existentes e especificações do projeto.
- `darp init` agora distribui as quatro skills de governança, os contratos
  completos de lifecycle e quality gates, e atualiza `darp.yml` de forma
  aditiva e segura.
- A versão do CLI é calculada a partir do Git durante o build, incluindo
  indicação de desenvolvimento e de worktree modificada quando aplicável.

## Como manter este arquivo

- Registre mudanças voltadas aos usuários em `Unreleased` durante o
  desenvolvimento.
- Use as categorias `Adicionado`, `Alterado`, `Corrigido`, `Removido`,
  `Segurança` e `Descontinuado` somente quando forem necessárias.
- Ao publicar uma versão, transforme `Unreleased` em uma seção versionada com
  a data no formato `AAAA-MM-DD` e crie uma nova seção `Unreleased` acima dela.
- Não registre cada commit: prefira mudanças relevantes para usuários,
  contribuidores e consumidores do CLI.

[Unreleased]: https://github.com/darpbr/darp-cli/compare/HEAD...main
