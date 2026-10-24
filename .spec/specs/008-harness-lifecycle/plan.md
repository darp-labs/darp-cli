# Plan 008 — Harness Lifecycle

## Status

Approved — plano aprovado em conjunto com a Spec 008. Implementar somente o escopo e os critérios definidos nos artefatos aprovados.

## Related Specification

- [spec.md](spec.md)
- [Spec 006 — Harness Manifest](../006-harness-manifest/spec.md)
- [ADR 001 — Harness Manifest Contract](../006-harness-manifest/adr/001-harness-manifest-contract.md)

## Estratégia de implementação

Implementar uma camada fina de CLI sobre a biblioteca existente `internal/project/harness`, sem duplicar parsing ou validação. Adicionar os subcomandos `harness inspect`, `harness init` e `harness validate`. Integrar o resultado do mesmo validador ao serviço `internal/project/doctor`.

Separar leitura/inspeção, validação e criação. As operações de inspeção, validação e Doctor são somente leitura. Somente `harness init` escreve, e apenas no caminho canônico do manifesto. A inicialização deve ser explícita, não destrutiva e segura contra sobrescrita.

Não criar abstração de runtime, adapter, registry ou sistema genérico de lifecycle nesta entrega.

## Áreas afetadas

- `internal/cli/` — roteamento, ajuda, saída e códigos de saída dos subcomandos.
- `internal/project/harness/` — somente extensões estritamente necessárias para distinguir ausência, erro de leitura e diagnósticos, e para criar o manifesto com segurança; não duplicar o parser.
- `internal/project/doctor/` — novo check Harness.
- testes de CLI, Harness e Doctor.
- `README.md` e `CHANGELOG.md`.
- `.spec/specs/008-harness-lifecycle/` — artefatos SDD.

Não alterar o contrato da Spec 006, `darp.yml`, comandos não relacionados, formatos nativos de agentes ou arquivos de assets existentes.

## Contrato de nomes, mensagens e códigos de saída

Esta seção fecha as decisões delegadas à T1 e é a fonte única dos tokens,
formatos e códigos usados pela implementação e pelos testes.

### Identificador válido de Harness

Um `metadata.name` válido:

- contém de 1 a 63 caracteres;
- usa apenas `[a-z0-9-]` (letras minúsculas ASCII, dígitos e hífen);
- começa e termina com `[a-z0-9]`;
- não contém hífens consecutivos.

Regex equivalente: `^[a-z0-9](?:[-a-z0-9]*[a-z0-9])?$`.

A Spec 006 continua exigindo apenas string não vazia para `metadata.name`; o
identificador acima é uma restrição adicional aplicada somente ao valor gerado
por `darp harness init`, sem alterar a validação da Spec 006.

### Normalização determinística de `metadata.name`

A partir do `filepath.Base` do diretório do projeto, aplicar em ordem:

1. `strings.ToLower` sobre a string inteira;
2. substituir todo caractere que não seja `[a-z0-9]` ou `-` por `-`
   (isso remove espaços, pontuação e Unicode não-ASCII, sem transliteração);
3. colapsar sequências de `-` em um único `-`;
4. remover `-` à esquerda e à direita.

Se o resultado for vazio ou não casar com o identificador válido acima, o
comando falha com `error: could not derive a valid harness name from directory <base>`,
sem criar arquivos e sem entrada interativa.

Exemplos:

| Diretório  | `metadata.name` |
|------------|-----------------|
| `darp-cli` | `darp-cli`      |
| `My Project` | `my-project`  |
| `café`     | `caf`           |
| `v1.2.3`   | `v1-2-3`        |
| `___`      | falha (vazio)   |
| `★ ★ ★`    | falha (vazio)   |

### Mensagens e códigos de saída

Tokens estáveis, adequados a scripts. Saída normal vai para `stdout`; erros de
uso/operacionais vão para `stderr`.

`inspect` — retorna `0` ao concluir a inspeção de `MISSING`, `VALID` ou
`INVALID`; `1` para erro de leitura ou uso inválido.

- `MISSING`: `MISSING .darp/harness.yaml`
- `VALID`: `VALID`, seguido de uma linha por campo, no mínimo `apiVersion`,
  `kind`, `metadata.name` e `metadata.version`, no formato `campo: valor`.
- `INVALID`: `INVALID`, seguido de uma linha por diagnóstico no formato
  `caminho: mensagem` (mesmo formato de `Diagnostic.String()`).
- Erro de leitura: `error: <err>` em `stderr`, código `1`.

`validate` — retorna `0` somente para manifesto válido; `1` para ausente,
inválido ou erro de leitura.

- válido: nenhuma linha de erro e código `0`;
- inválido: uma linha por diagnóstico em `stderr`, código `1`;
- ausente: `error: harness manifest not found: .darp/harness.yaml`, código `1`;
- erro de leitura: `error: <err>`, código `1`.

`init` — retorna `0` após criação; `1` em qualquer falha.

- sucesso: `✔ created .darp/harness.yaml`;
- já existe: `error: .darp/harness.yaml already exists`, código `1`;
- nome inválido: `error: could not derive a valid harness name from directory <base>`, código `1`;
- erro de escrita: `error: <err>`, código `1`.

`doctor` — reusa `Result.ExitCode()` existente (`1` se houver qualquer `FAIL`).
O formato do resumo e das linhas por check não muda.

## Marcos

### M1 — Contrato de saída e nome do Harness

- Fixar o contrato de nomes, mensagens e códigos de saída conforme a seção "Contrato de nomes, mensagens e códigos de saída" acima.
- Confirmar `metadata.version: 0.1.0` como valor inicial.
- Usar `WARNING` para Harness ausente no Doctor, sem tornar o manifesto obrigatório para os demais contratos do projeto.

### M2 — Operações somente leitura

- Implementar `inspect` e `validate` usando o caminho canônico e o validador existente.
- Diferenciar ausência, manifesto inválido e erro de leitura.
- Testar que não há alterações no filesystem.

### M3 — Inicialização segura

- Implementar criação explícita do manifesto mínimo.
- Recusar sobrescrita de qualquer arquivo já existente no caminho canônico.
- Gravar com criação exclusiva atômica: `os.OpenFile(caminho, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)`. O `O_EXCL` garante que uma única invocação vence sob concorrência; se a escrita falhar após a abertura, remover o arquivo recém-criado para não deixar manifesto parcial. Nunca usar `os.Rename` sobre o caminho final (sobrescreveria um arquivo existente).
- Criar o diretório `.darp/` com `os.MkdirAll` (0o755) se estiver ausente.
- Testar diretório `.darp/` ausente, manifesto já existente e falhas de escrita.

### M4 — Integração com Doctor

- Adicionar um check `Harness` que reutilize a validação existente.
- Posicionar o check ao final da lista (após `Assets`), preservando a ordem e a semântica dos oito checks existentes.
- Aplicar os estados especificados para ausente (`WARNING`), válido (`PASS`), inválido (`FAIL`) e erro de leitura (`FAIL`).
- Preservar a natureza somente leitura.

### M5 — Documentação e validação final

- Atualizar ajuda, README e changelog.
- Executar testes unitários e de integração.
- Executar `make verify` e `git diff --check`.
- Revisar arquitetura, documentação, compatibilidade e evidências antes de concluir.

## Estratégia de validação

- Testes unitários para estados retornados por `internal/project/harness`.
- Testes de CLI para subcomandos, argumentos inválidos, mensagens e códigos de saída.
- Testes de inicialização em diretório vazio e em projeto com `.darp/` existente.
- Testes de normalização de nomes, incluindo espaços, pontuação, Unicode e nomes que não possam ser normalizados.
- Teste de criação concorrente: N goroutines invocando a inicialização sobre o mesmo diretório temporário; exatamente uma deve vencer (código 0) e as demais devem falhar com "já existe"; o conteúdo resultante deve ser o manifesto mínimo e nunca sobrescrito.
- Testes de manifesto preexistente comprovando que bytes e metadados não mudam.
- Testes de falha de escrita comprovando ausência de manifesto parcial e preservação de arquivos não relacionados.
- Testes do Doctor para manifesto ausente, válido, inválido e erro de leitura.
- Testes de regressão dos oito checks atuais e de seus códigos de saída.
- Snapshot do filesystem antes/depois de `inspect`, `validate` e `doctor`.
- `make verify`.
- `git diff --check`.

## Recuperação ou rollback

A mudança deve ser reversível removendo o roteamento dos novos subcomandos, o check do Harness e os testes/documentação associados. Não deve exigir migração de arquivos existentes. O manifesto criado por `darp harness init` é um arquivo novo versionável pelo usuário; o CLI não deve removê-lo automaticamente como parte de rollback.

## Premissas

- A Spec 006 é a fonte de verdade do formato e da validação estrutural.
- `internal/project/harness` é a fonte única das regras de validação.
- O Doctor permanece somente leitura.
- A ausência do manifesto não invalida, por si só, um projeto DARP existente.
- Nenhuma operação desta spec precisa de rede, Registry, SDK, Runtime ou agente instalado.
- A inicialização não detecta stack nem configura capabilities automaticamente.
