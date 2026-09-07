# Tasks 006 — Harness Manifest

## Status

Completed — implementação e validação concluídas.

## Related Plan

- [`plan.md`](plan.md)

## Regras para o implementador

- Não iniciar antes da aprovação da Spec 006, Plan 006 e ADR 001.
- Implementar somente o contrato aprovado; execução de agentes, modelos,
  ferramentas, workflows e MCP está fora do escopo.
- Não criar subcomando ou runtime de Harness nesta spec.
- Preservar `darp.yml`, assets físicos e arquivos desconhecidos.
- Tasks bloqueadas permanecem desmarcadas e devem registrar evidência e motivo.

## Task List

### T1 — Modelo e contrato

- [x] Criar o modelo interno para `apiVersion`, `kind`, `metadata` e blocos
      declarativos opcionais.
- [x] Definir representação para campos desconhecidos sem perder informação
      relevante ao ler ou reemitir o manifesto; comentários e ordem textual não
      fazem parte da garantia.
- [x] Garantir que o modelo não dependa de provider, agente, runtime ou stack.

### T2 — Parsing e validação

- [x] Implementar leitura de `.darp/harness.yaml` com YAML válido e diagnóstico
      determinístico.
- [x] Validar raiz, `apiVersion`, `kind`, `metadata.name` e
      `metadata.version`.
- [x] Validar tipos dos campos conhecidos e rejeitar estruturas malformadas.
- [x] Validar referências de assets contra os limites de caminho do projeto.
- [x] Aceitar e preservar campos desconhecidos.

### T3 — Isolamento e segurança

- [x] Garantir que a validação seja read-only e não execute arquivos, comandos,
      workflows, MCP servers ou modelos.
- [x] Garantir ausência de acesso à rede e de detecção automática de stack.
- [x] Garantir que falhas não alterem `.darp/harness.yaml`, `darp.yml` ou
      qualquer asset.
- [x] Garantir que a biblioteca não execute referências declaradas nem crie
      integração com runtime.

### T4 — Testes

- [x] Testar manifesto mínimo válido e blocos opcionais.
- [x] Testar YAML inválido, raiz inválida, campos obrigatórios ausentes e tipos
      incompatíveis.
- [x] Testar `apiVersion` não suportada e `kind` incorreto.
- [x] Testar providers múltiplos e ausência de dependência de fornecedor.
- [x] Testar assets ainda não instalados, absolutos, com traversal e
      incompletos; a validação não exige disponibilidade local.
- [x] Testar campos desconhecidos, determinismo e ausência de efeitos colaterais.
- [x] Testar regressão dos contratos atuais com `go test ./...`.

### T5 — Documentação e revisão

- [x] Documentar somente as capacidades efetivamente implementadas.
- [x] Atualizar `CHANGELOG.md` se uma interface de usuário for adicionada.
- [x] Executar `git diff --check`.
- [x] Completar revisão de arquitetura, documentação, testes e quality gates.
- [x] Revisar o diff final contra Spec 006, Plan 006 e ADR 001.

## Validation Checklist

- [x] Spec 006, Plan 006 e ADR 001 estão aprovados.
- [x] O manifesto canônico é `.darp/harness.yaml`.
- [x] A validação é determinística, read-only e sem rede ou execução.
- [x] Campos desconhecidos são preservados.
- [x] O contrato de `darp.yml` e a descoberta de assets permanecem compatíveis.
- [x] Todos os testes passam com evidência registrada.
- [x] Nenhum bloqueio permanece sem registro.

## Notas e assumptions

- Provider Profile, Tool Policy e Evaluation Contract devem permanecer
  extensões ou specs separadas.
- Não marcar tasks como concluídas antes da validação correspondente.

## Evidência de validação

- `GOCACHE=/tmp/darp-cli-go-build-cache go test ./...` — PASS.
- `GOCACHE=/tmp/darp-cli-go-build-cache go vet ./...` — PASS.
- `GOCACHE=/tmp/darp-cli-go-build-cache make build` — PASS.
- `GOCACHE=/tmp/darp-cli-go-build-cache make lint` — PASS (fallback para
  `go vet`, sem issues).
- `git diff --check` — PASS.

## Limite aprovado

A implementação desta lista entrega somente biblioteca interna de parsing e
validação estrutural. Execução, orquestração, permissões, segredos,
observabilidade e avaliação operacional ficam para um runtime futuro e não
serão implementados nesta spec.
