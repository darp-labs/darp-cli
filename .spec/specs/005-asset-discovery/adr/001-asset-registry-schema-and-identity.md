# ADR 001 — Schema e identidade de assets descobertos

## Status

Accepted — aprovada com a Spec 005 e o Plan 005.

## Contexto

A Spec 005 registra assets preexistentes de varias ferramentas de assistente de
codigo. O contrato atual de `darp.yml` ja usa `skills` para skills DARP sob
`.agents/skills/<nome>`, entao registrar assets externos nessa mesma secao
criaria conflito de semantica e uma falsa promessa de execucao.

Tambem existem caminhos que podem ser reconhecidos por mais de uma familia,
como `AGENTS.md`, e nomes que podem colidir entre comandos, skills ou regras de
ferramentas diferentes.

## Decisao

1. Assets externos serao registrados em uma nova secao `assets` do `darp.yml`.
2. Cada asset registrado tera `path`, `family` e `type`.
3. `path` e o caminho relativo do arquivo reconhecido, nunca um diretorio pai.
4. `type` descreve o tipo semantico do asset, como `skill`, `instruction`,
   `prompt`, `command`, `hook`, `agent`, `rule`, `config` ou `extension`.
5. A identidade de persistencia e o trio `(path, family, type)`.
6. O `darp init` adiciona apenas entradas novas e preserva o `darp.yml`
   existente. Entradas ja registradas sao informadas como `already registered`.
7. O CLI nao escolhe automaticamente entre candidatos ambiguos. Se um mesmo
   caminho fisico ou nome semantico gerar conflito entre familias ou tipos, a
   nova entrada conflitante deve ser reportada e nao persistida automaticamente.
8. Symlinks nao sao seguidos nem registrados na v1.
9. `darp doctor` valida existencia, caminho relativo dentro do projeto e tipo
   esperado. Asset ausente e warning.
10. A identidade semantica e `(categoria, nome normalizado)`, usando o tipo
  como categoria. Skills e extensions usam o diretorio pai como nome; os
  demais tipos usam o nome do arquivo sem extensao ou sufixo convencional,
  com comparacao sem diferenca entre maiusculas e minusculas.
11. Symlinks registrados manualmente sao invalidos e devem resultar em `FAIL`
  no `darp doctor`; arquivos ausentes continuam como `WARNING`.

Exemplo:

```yaml
assets:
  - path: .github/instructions/caveman-mode.instructions.md
    family: github
    type: instruction
  - path: .claude/commands/prompt.md
    family: claude
    type: command
```

## Alternativas consideradas

- Reusar `skills`: rejeitado porque `skills` ja possui contrato especifico para
  skills DARP em `.agents/skills/<nome>`.
- Registrar diretorios de familia, como `.github/` ou `.claude/`: rejeitado
  porque gera ruido e nao identifica o asset real.
- Escolher uma familia preferida para arquivos compartilhados, como
  `AGENTS.md`: rejeitado porque esconderia ambiguidades entre ferramentas.
- Seguir symlinks dentro do projeto: rejeitado na v1 para manter a descoberta
  simples, segura e previsivel.

## Consequencias

- O schema diferencia assets externos de skills DARP.
- A descoberta e aditiva e idempotente.
- Ambiguidades exigem decisao explicita futura do usuario ou comando dedicado.
- `darp doctor` passa a poder alertar sobre assets registrados ausentes sem
  bloquear o projeto inteiro.
- A v1 pode deixar de registrar automaticamente alguns candidatos validos em
  troca de preservar seguranca e previsibilidade.

## Follow-up actions

- Implementar parser/persistencia YAML segura seguindo o precedente da Spec
  004.
- Documentar os estados reportados por `darp init`: encontrado, adicionado, ja
  registrado, ignorado, ambiguo/conflitante e ausente.
- Definir em spec futura um comando explicito para reconciliar, remover ou
  resolver assets ambiguos.
