# Spec 005 — Descoberta e registro de assets existentes

## Status

Aprovada e implementada — contrato de descoberta, schema,
identidade e validacao definidos nesta spec e no
[ADR 001](adr/001-asset-registry-schema-and-identity.md).

## Related Specifications

- [Spec 004 — Scaffold de skills de governança](../004-init-governance-skills/spec.md)
- [Spec 002 — `darp doctor`](../002-doctor/spec.md)
- [ADR 001 — Schema e identidade de assets descobertos](adr/001-asset-registry-schema-and-identity.md)

## Contexto

O DARP CLI gerencia assets reutilizáveis para desenvolvimento assistido por IA,
mas um repositório existente pode conter instruções em locais convencionais de
diversas ferramentas, como `.github/`, `.claude/`, `CLAUDE.md`, `.codex/`,
`AGENTS.md`, `.cursor/` e `.gemini/`.

Atualmente, o `darp init` cria apenas a estrutura padrão do DARP. Ele não
descobre nem registra arquivos já presentes no repositório.

Esta demanda é separada da Spec 004: a Spec 004 trata do scaffold determinístico
de contratos DARP; esta spec trata da descoberta de conteúdo preexistente e da
sua representação no `darp.yml`.

## Problema

Ao inicializar ou atualizar um repositório existente, o usuário não recebe uma
visão registrada dos assets que já fazem parte das convenções de agentes e
ferramentas do projeto. Isso dificulta a composição, auditoria e futura gestão
desses assets pelo DARP.

## Objetivos

1. Detectar, de forma determinística e sem rede, assets reconhecidos pelo DARP.
2. Registrar os assets descobertos no `darp.yml` sem sobrescrever arquivos ou
   configurações existentes.
3. Preservar compatibilidade com projetos que não possuem esses arquivos.
4. Informar o que foi descoberto, ignorado ou considerado ambíguo.
5. Manter a descoberta independente de provedor, linguagem e framework.

## Não objetivos

- executar instruções, workflows ou assets descobertos;
- interpretar o conteúdo semântico dos arquivos;
- instalar ferramentas ou acessar a rede;
- mover, copiar, converter ou sobrescrever assets existentes;
- inferir a stack da aplicação;
- registrar automaticamente todo arquivo encontrado no repositório;
- resolver automaticamente conflitos de identidade entre ferramentas;
- remover assets registrados automaticamente quando o arquivo fisico nao existir
  mais.

## Escopo

A primeira versão avalia somente os padrões explicitamente suportados abaixo. A
descoberta registra arquivos de asset, nao diretorios genericos.

| Familia | Padrao reconhecido | Tipo |
| --- | --- | --- |
| GitHub | `.github/copilot-instructions.md` | `instruction` |
| GitHub | `.github/instructions/**/*.instructions.md` | `instruction` |
| GitHub | `.github/prompts/**/*.prompt.md` | `prompt` |
| GitHub | `.github/skills/**/SKILL.md` | `skill` |
| GitHub | `.github/hooks/*.json` | `hook` |
| GitHub | `.github/agents/*.md` | `agent` |
| GitHub | `AGENTS.md`, `CLAUDE.md`, `GEMINI.md` | `instruction` |
| Claude | `CLAUDE.md`, `.claude/CLAUDE.md` | `instruction` |
| Claude | `.claude/commands/**/*.md` | `command` |
| Claude | `.claude/skills/**/SKILL.md` | `skill` |
| Claude | `.claude/agents/**/*.md` | `agent` |
| Claude | `.claude/settings.json` | `hook` |
| Codex | `AGENTS.md` | `instruction` |
| Codex | `.codex/skills/**/SKILL.md` | `skill` |
| Cursor | `.cursor/rules/**/*.mdc` | `rule` |
| Cursor | `.cursor/skills/**/SKILL.md` | `skill` |
| Cursor | `.cursor/hooks.json` | `hook` |
| Cursor | `AGENTS.md` | `instruction` |
| Gemini | `GEMINI.md` | `instruction` |
| Gemini | `.gemini/commands/**/*.toml` | `command` |
| Gemini | `.gemini/skills/**/SKILL.md` | `skill` |
| Gemini | `.gemini/settings.json` | `config` |
| Gemini | `.gemini/extensions/*/gemini-extension.json` | `extension` |

Arquivos pessoais ou locais, como `.claude/settings.local.json`, nao entram no
escopo da v1.

## Contrato de registro

Assets descobertos serao registrados em uma nova secao `assets` do `darp.yml`.
Cada entrada deve conter `path`, `family` e `type`.

```yaml
assets:
  - path: .github/instructions/caveman-mode.instructions.md
    family: github
    type: instruction
  - path: .claude/commands/prompt.md
    family: claude
    type: command
```

`path` e sempre o caminho relativo do asset reconhecido, nunca o diretorio pai.
`type` usa o tipo semantico do asset (`skill`, `instruction`, `prompt`,
`command`, `hook`, `agent`, `rule`, `config` ou `extension`).

## Fluxos

### Inicialização com descoberta

1. O usuário executa `darp init`.
2. O CLI examina somente padrões suportados nas raizes conhecidas do projeto.
3. O CLI registra resultados novos no `darp.yml` e informa um resumo.
4. Nenhum asset encontrado é executado ou modificado.

### Atualização de projeto registrado

O CLI preserva registros existentes, adiciona somente entradas novas conforme
as regras de identidade definidas nesta spec e informa conflitos ou arquivos
removidos sem apagá-los silenciosamente.

Se o `darp.yml` ja existir, ele nao sera sobrescrito. A atualizacao e aditiva:
novos assets reconhecidos podem ser adicionados, entradas ja registradas serao
informadas como `already registered`, e conflitos ou ambiguidades serao
informados sem escolha automatica.

Se o `darp.yml` estiver invalido ou em formato que nao permita edicao segura, o
CLI nao deve altera-lo. O usuario devera remover ou corrigir a configuracao
antiga e executar `darp init` novamente, ou fazer backup, executar o init e
reconciliar manualmente com o arquivo antigo.

### Repositório sem assets suportados

O comando termina normalmente, não cria registros artificiais e informa que
nenhum asset suportado foi encontrado.

## Requisitos

- R1. A descoberta deve ser somente local e determinística.
- R2. Caminhos registrados devem ser relativos e permanecer dentro do projeto.
- R3. A descoberta não deve executar comandos, scripts ou acessar a rede.
- R4. Assets existentes devem ser preservados integralmente.
- R5. A atualização do `darp.yml` deve ser aditiva e segura, seguindo a Spec
  004 quando aplicável.
- R6. O resultado deve distinguir encontrado, já registrado, ignorado e
  ambíguo/conflitante.
- R7. O comportamento deve ser testável em repositórios vazios, monorepos e
  projetos com múltiplas famílias.
- R8. A descoberta deve ocorrer automaticamente durante `darp init`.
- R9. O registro deve ser somente por assets reconhecidos, nunca por diretorios
  genericos.
- R10. Symlinks nao devem ser seguidos nem registrados na v1; devem ser
  informados como ignorados.
- R11. Arquivos fisicamente removidos devem permanecer registrados ate acao
  explicita futura; `darp doctor` deve informa-los como ausentes.
- R12. O `darp doctor` deve validar existencia, caminho relativo dentro do
  projeto e tipo esperado dos assets registrados. Asset ausente deve ser
  `warning`, nao falha bloqueante.
- R13. A descoberta deve usar raizes conhecidas do projeto: a raiz atual do
  `darp init`, os diretorios convencionais de familias suportadas nessa raiz e
  arquivos raiz suportados. A v1 nao atravessa diretorios arbitrarios de
  monorepo.

## Compatibilidade e arquitetura atual

O `darp.yml` atualmente possui uma seção `skills`, e o `darp doctor` exige que
as entradas dessa seção apontem para `.agents/skills/<nome>`. Portanto, esta
spec não pode registrar diretamente `.github/`, `.claude/`, `AGENTS.md` ou
outros caminhos externos como skills sem alterar esse contrato.

Assets externos serao registrados em `assets`, nao em `skills`. A secao
`skills` permanece reservada para skills DARP sob `.agents/skills/<nome>`.
Assets externos nao sao convertidos para skills DARP nesta spec.

## Identidade e conflitos

A identidade de persistencia de um asset e o trio `(path, family, type)`.

O CLI nao deve sobrescrever nem escolher automaticamente entre candidatos
ambiguos. Quando o mesmo caminho fisico for reconhecido por mais de uma
familia ou tipo, ou quando houver nomes semanticos duplicados que a ferramenta
externa trataria como conflito, o CLI deve informar a duplicidade/ambiguidade e
nao persistir a nova entrada conflitante automaticamente.

A identidade semantica e o par `(categoria, nome normalizado)`. A categoria e o
tipo do asset. Para `skill` e `extension`, o nome e o diretorio que contem o
arquivo reconhecido; para os demais tipos, e o nome do arquivo sem a extensao
ou sufixo convencional (`.md`, `.mdc`, `.toml`, `.json`, `.prompt.md` ou
`.instructions.md`). A comparacao ignora maiusculas e minusculas. Candidatos
com a mesma identidade semantica, mesmo em caminhos diferentes ou familias
diferentes, sao conflitantes e nao sao persistidos automaticamente.

Entradas com o mesmo `(path, family, type)` ja presentes no `darp.yml` devem ser
informadas como `already registered` e nao duplicadas.

## Monorepo e exclusoes

A v1 nao tera exclusoes configuraveis porque a descoberta nao faz varredura
ampla no repositorio. Ela avalia apenas raizes conhecidas e padroes
explicitamente suportados.

Em monorepos, `darp init` considera como projeto o diretorio em que esta sendo
executado. Ele nao sobe nem desce para descobrir workspaces ou pacotes
arbitrarios. Caso a raiz do monorepo seja inicializada, os padroes reconhecidos
sob as raizes conhecidas dessa raiz serao avaliados.

## Fontes de padrões

Os padroes foram definidos a partir de documentacao publica dos assistentes e
de contratos ja usados por este projeto:

- GitHub Copilot customization cheat sheet:
  https://docs.github.com/en/copilot/reference/customization-cheat-sheet
- GitHub Copilot agent skills:
  https://docs.github.com/en/copilot/how-tos/copilot-on-github/customize-copilot/customize-cloud-agent/add-skills
- Claude Code memory, commands, skills and hooks:
  https://code.claude.com/docs/en/memory,
  https://code.claude.com/docs/en/commands,
  https://code.claude.com/docs/pt/skills,
  https://code.claude.com/docs/en/hooks
- Cursor rules, skills and hooks:
  https://cursor.com/docs/rules,
  https://cursor.com/docs/skills,
  https://cursor.com/docs/hooks
- Gemini CLI context, commands, settings and extensions:
  https://github.com/google-gemini/gemini-cli/blob/main/docs/cli/gemini-md.md,
  https://github.com/google-gemini/gemini-cli/blob/main/docs/cli/custom-commands.md,
  https://github.com/google-gemini/gemini-cli/blob/main/docs/cli/settings.md,
  https://github.com/google-gemini/gemini-cli/blob/main/docs/extensions/reference.md
- OpenAI Developers documentation index for Codex customization and skills:
  https://developers.openai.com/

## Critérios de aceitação

- [x] O contrato final de registro de assets está definido e documentado.
- [x] A lista de padrões suportados e regras de exclusão está definida.
- [x] A descoberta não executa, copia, move ou sobrescreve assets encontrados.
- [x] O `darp.yml` preserva configurações e registros existentes.
- [x] Reexecuções são determinísticas e não geram duplicatas.
- [x] Conflitos, duplicidades, symlinks e arquivos removidos possuem comportamento
      definido e testado.
- [x] O `darp doctor` foi atualizado ou explicitamente declarado compatível.
- [x] Testes cobrem repositório vazio, assets únicos, múltiplas famílias,
      monorepo, caminhos inválidos e falhas de atualização.
- [x] README, changelog e documentação do `darp init` refletem o comportamento.

## Perguntas resolvidas

1. Descoberta automatica durante `darp init`.
2. Nova secao `assets` no `darp.yml`.
3. Registro por arquivo de asset reconhecido.
4. Padroes suportados definidos por familia nesta spec.
5. `AGENTS.md`, `CLAUDE.md` e `GEMINI.md` sao assets independentes por familia;
   quando o mesmo caminho for reconhecido por multiplas familias, isso e
   ambiguidade a informar, nao escolha automatica.
6. Duplicatas, nomes iguais e conflitos devem ser informados sem sobrescrita ou
   selecao automatica.
7. Symlinks ignorados na v1.
8. Arquivos removidos permanecem registrados ate acao explicita futura; doctor
   informa como ausentes.
9. `darp doctor` valida existencia, caminho e tipo; ausencia e warning.
10. A descoberta usa raizes conhecidas e nao atravessa arbitrariamente limites
    de monorepo.
11. Sem exclusoes configuraveis na v1, pois nao ha varredura ampla.
12. ADR necessario e criado nesta spec.

## Riscos

- registrar arquivos demais pode gerar ruído;
- diferentes ferramentas podem usar o mesmo nome para conceitos distintos;
- seguir symlinks pode expor arquivos fora do projeto;
- atualizar YAML preservando configuração local continua sendo sensível;
- descoberta automática pode surpreender usuários ao alterar `darp.yml`;
- tratar instruções como skills pode criar uma falsa promessa de execução.
