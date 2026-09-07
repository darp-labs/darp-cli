# Plan 004 — Scaffold de skills de governança no `darp init`

## Status

Approved and implemented — validated on 2026-08-27.

## Related Specifications

- [`spec.md`](spec.md)
- [`../003-governance/spec.md`](../003-governance/spec.md)
- [`adr/001-bootstrap-assets-and-config-update.md`](adr/001-bootstrap-assets-and-config-update.md)

## Estratégia

Estender o serviço de inicialização para tratar os contratos de governança como
um conjunto determinístico de ativos padrão. O serviço deve continuar usando a
abstração de filesystem existente, criando somente arquivos ausentes e
aplicando uma atualização aditiva e segura ao `darp.yml` quando necessário.

A seleção de conteúdo não deve depender de inspeção do projeto-alvo. A fonte
canônica é `internal/project/init/assets/`, incorporada com `go:embed`; os
contratos equivalentes na raiz são sincronizados e verificados por teste.

## Dependências e ordem

1. Confirmar Spec 003 aprovada e aprovar Spec 004, Plan 004 e ADR 001.
2. Consolidar os conteúdos canônicos no diretório de ativos embarcados.
3. Estender a abstração de filesystem e implementar a atualização segura do
   `darp.yml` conforme o ADR 001.
4. Estender o serviço de bootstrap e preservar idempotência.
5. Adicionar testes unitários, de reparo, preservação e fixtures agnósticas.
6. Atualizar documentação e changelog, e validar com `darp doctor`.

## Áreas afetadas

- `internal/project/init/`
- testes de `internal/project/init/`
- `internal/project/init/assets/`
- `README.md` e `CHANGELOG.md`
- documentação da Spec 001, se o contrato de preservação for alterado

Não devem ser alterados pelo escopo normal:

- `internal/project/doctor/`, salvo ajustes de teste estritamente necessários;
- `.darp/` e `.agents/skills/` do repositório sem atualizar também a fonte
  canônica correspondente;
- comportamento de workflows e execução de skills;
- skills específicas existentes, como `security-review`.

## Marcos

### M1 — Contratos canônicos

Definir e versionar os conteúdos completos dos quatro `SKILL.md`, lifecycle e
quality gates, verificando neutralidade de stack.

### M2 — Configuração aditiva

Implementar a inclusão das quatro entradas em `darp.yml` para projetos novos e
o reparo seguro de configurações existentes.

### M3 — Scaffold e preservação

Criar diretórios e arquivos ausentes, atualizar somente placeholders históricos
exatos e preservar ativos customizados.

### M4 — Testes

Cobrir projeto novo, projeto parcial, reexecução, placeholders históricos,
ativos customizados, YAML inválido, campos desconhecidos, skills adicionais e
falhas de filesystem.

### M5 — Compatibilidade e documentação

Validar `darp doctor`, atualizar documentação e confirmar que nenhuma etapa
executa comandos ou detecta o tipo do projeto.

## Estratégia de validação

- `go test ./...`;
- testes do serviço com filesystem temporário;
- comparação do conteúdo mínimo dos quatro `SKILL.md`;
- fixtures sem código, Java/Spring, Python/FastAPI, Go e monorepo;
- execução de `darp doctor` em cada fixture scaffoldada;
- segunda execução do `init` sem diferenças nos arquivos existentes;
- inspeção do diff para garantir que `security-review` não foi alterada;
- `git diff --check`;
- verificação de sincronização entre ativos embarcados e contratos versionados
  na raiz.

## Recuperação

O serviço deve falhar sem apagar arquivos. Em caso de erro na validação ou
atualização do `darp.yml`, manter o arquivo original, não criar ativos nessa
execução e reportar a causa. Depois que a configuração for validada e salva,
um erro ao criar ativos não deve remover os ativos já criados; a próxima
execução deve conseguir reparar os ausentes.

## Premissas

- As quatro skills da Spec 003 são os templates padrão oficiais.
- O `darp doctor` continua validando skills encontradas no diretório, além das
  entradas registradas na configuração.
- Um `darp.yml` existente inválido ou não editável impede a criação de ativos
  nessa execução e permanece intacto.
- O workflow `implement.yaml` permanece sem execução automática.
- A neutralidade é obtida por descoberta contextual dentro das skills, não por
  lógica de detecção de stack no `darp init`.
