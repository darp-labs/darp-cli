# Plan 006 — Harness Manifest

## Status

Approved — implementação autorizada após aprovação da Spec 006 e ADR 001.

## Related Artifacts

- [`spec.md`](spec.md)
- [`adr/001-harness-manifest-contract.md`](adr/001-harness-manifest-contract.md)
- [`docs/AGENTIC_HARNESS.md`](../../../docs/AGENTIC_HARNESS.md)

## Estratégia de implementação

Formalizar o contrato do Harness como um manifesto independente de `darp.yml`,
armazenado em `.darp/harness.yaml`. A implementação limita-se a biblioteca
interna de leitura e validação estrutural determinística, mantendo providers,
ferramentas, políticas e workflows como dados declarativos. Não será criado
runtime, executor ou subcomando de Harness nesta spec.

O modelo interno deve separar parsing, validação e diagnóstico. Nenhuma camada
deve executar conteúdo referenciado, acessar rede ou descobrir a stack do
projeto. Campos desconhecidos devem ser preservados pelo modelo de leitura. A
preservação de comentários e ordem textual não faz parte do contrato.

## Áreas afetadas

- `internal/project/harness/` para modelo, parsing e validação;
- `internal/cli/` não será alterado nesta spec;
- testes unitários e de integração do contrato;
- `.darp/harness.yaml` como fixture ou manifesto do próprio repositório, se
  aprovado pela implementação;
- documentação técnica e `README.md` apenas para comportamento implementado;
- `CHANGELOG.md` somente se uma interface de usuário for adicionada.

Não alterar o comportamento de `darp init`, `darp doctor`, a descoberta de
assets ou a configuração existente sem uma decisão adicional registrada.

## Marcos

### M1 — Contrato e modelo

Definir tipos para identidade do manifesto, metadata e blocos declarativos,
mantendo extensões desconhecidas sem dependência de provider.

### M2 — Parsing e validação

Implementar parsing YAML, validações obrigatórias, segurança de referências e
diagnósticos estáveis. A validação deve aceitar campos desconhecidos e rejeitar
estruturas conhecidas com tipos inválidos.

### M3 — Isolamento e compatibilidade

Garantir que a leitura seja read-only, sem rede, execução, descoberta de stack
ou alteração de `darp.yml`. Confirmar que o contrato de assets atual continua
sendo a fonte de identidade dos assets locais.

### M4 — Testes e documentação

Adicionar fixtures válidas e inválidas, atualizar documentação somente para
capacidades implementadas e revisar o resultado contra esta spec e o ADR.

## Estratégia de validação

- `GOCACHE=/tmp/darp-cli-go-build-cache go test ./...`;
- testes de parsing e validação para todos os critérios de aceitação;
- testes de tipos incompatíveis, versões não suportadas e YAML inválido;
- testes de campos desconhecidos e preservação semântica;
- testes de caminhos de assets, traversal, absolutos e referências incompletas;
- testes que comprovem ausência de execução, rede e efeitos colaterais;
- teste de leitura do caminho canônico `.darp/harness.yaml` sem escrita;
- `git diff --check`;
- revisão de arquitetura, documentação e compatibilidade conforme os gates do
  projeto.

## Recuperação e compatibilidade

A validação não deve escrever no manifesto nem em qualquer outro arquivo. Uma
falha de parsing ou validação não pode alterar o estado do projeto. Qualquer
alteração posterior do manifesto deverá usar escrita segura e ser especificada
separadamente.

`apiVersion: darp.dev/v1alpha1` é a única versão suportada nesta primeira
implementação. Novas versões devem ser introduzidas por migração ou contrato
explícito, nunca por interpretação silenciosa.

## Assumptions

- `.darp/harness.yaml` é separado de `darp.yml` para manter distintos o
  contrato do projeto e a definição do Harness.
- Os blocos opcionais terão estrutura extensível; seus contratos detalhados
  serão refinados em specs futuras.
- O modelo interno preserva nós YAML desconhecidos, mas não garante round-trip
  textual de comentários e ordem.
- A primeira implementação pode expor apenas biblioteca interna e testes; uma
  subcomanda CLI permanece fora desta spec e exigirá especificação própria.
