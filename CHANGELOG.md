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

### Adicionado

- Descoberta automática de assets preexistentes durante `darp init`, com
  registro aditivo e seguro na nova seção `assets` do `darp.yml`. A descoberta
  é local, determinística e somente leitura, reconhece padrões das famílias
  GitHub, Claude, Codex, Cursor e Gemini, ignora symlinks e reporta entradas
  já registradas e caminhos ambíguos sem sobrescrever nada.
- `darp doctor` agora valida os assets registrados (existência, caminho
  relativo dentro do projeto, família e tipo esperados); asset ausente é
  reportado como `WARNING` sem bloquear o projeto.

### Corrigido

- `darp init` agora ignora arquivos não suportados dentro dos diretórios
  convencionais de assets, evitando falhas durante a descoberta.
- A descoberta agora centraliza a classificação oficial dos assets, reporta
  conflitos semânticos por categoria e nome normalizado e impede seu registro
  automático.
- `darp doctor` agora rejeita combinações incompatíveis de caminho, família e
  tipo, além de symlinks registrados manualmente.
- `darp init` agora valida a configuração e os assets descobertos antes de
  criar contratos, limpa arquivos temporários quando a substituição falha e
  desfaz alterações parciais após falhas de persistência.
- `darp doctor` agora classifica um diretório registrado no lugar de um asset
  como `FAIL`; somente um asset ausente continua sendo `WARNING`.

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
