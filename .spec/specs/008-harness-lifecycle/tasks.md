# Tasks 008 — Harness Lifecycle

## Status

Approved — artefatos aprovados. Todas as tarefas permanecem desmarcadas até que a implementação e sua validação correspondente sejam concluídas.

## Related Plan

- [plan.md](plan.md)
- [spec.md](spec.md)

## Regras para o implementador

- Não iniciar antes da aprovação da Spec 008, do plano e das tasks.
- Reutilizar o pacote `internal/project/harness`; não duplicar parsing ou validação.
- Preservar o contrato da Spec 006 e o arquivo `darp.yml`.
- Não criar runtime, registry, SDK, adapters, projeções ou instalação de assets.
- `inspect`, `validate` e `doctor` devem permanecer somente leitura.
- `init` não pode sobrescrever arquivos existentes nem deixar manifesto parcial.
- Manter todas as tasks desmarcadas até a validação correspondente passar.
- Registrar comandos, resultados e evidências reais; não declarar testes executados sem executá-los.

## Task List

- [ ] **T1 — Confirmar contrato de nomes, mensagens e códigos de saída.** O contrato está definido no plano (seção "Contrato de nomes, mensagens e códigos de saída"); validar que a normalização atende à Spec 006 e registrar os exemplos-limite nos testes.
- [ ] **T2 — Preparar testes do pacote Harness para estados do filesystem.** Cobrir manifesto ausente, válido, inválido e erros de leitura sem alterar o contrato de parsing.
- [ ] **T3 — Implementar `darp harness inspect`.** Reutilizar o loader/validador existente; exibir estado e metadados; comprovar ausência de efeitos colaterais.
- [ ] **T4 — Implementar `darp harness validate`.** Exibir diagnósticos e aplicar os códigos de saída definidos; não duplicar validação.
- [ ] **T5 — Implementar `darp harness init`.** Criar apenas o manifesto mínimo; recusar sobrescrita com criação exclusiva; usar estratégia segura contra concorrência e arquivo parcial; testar falhas e preservação de arquivos.
- [ ] **T6 — Integrar Harness ao `darp doctor`.** Acrescentar um check reutilizando o pacote Harness e aplicar PASS/WARNING/FAIL conforme a Spec 008.
- [ ] **T7 — Atualizar testes de regressão do Doctor e da CLI.** Confirmar os checks existentes, a saída resumida, os códigos de saída e o comportamento em diretórios sem Harness.
- [ ] **T8 — Atualizar ajuda e documentação.** Documentar os novos comandos e seus limites sem sugerir ativação ou execução de recursos.
- [ ] **T9 — Atualizar `CHANGELOG.md`.** Registrar a funcionalidade em `Unreleased` apenas após sua implementação.
- [ ] **T10 — Executar validação completa.** Executar testes focados, `make verify` e `git diff --check`; registrar resultados e corrigir falhas dentro do escopo.
- [ ] **T11 — Revisão final de arquitetura e documentação.** Confirmar aderência à Constituição, à Spec 006, aos ADRs existentes e à interoperabilidade; não renumerar ADRs existentes.

## Validation Checklist

- [ ] Spec, plano e tasks aprovados antes de iniciar implementação.
- [ ] Testes do pacote Harness passam.
- [ ] Testes de CLI cobrem todos os subcomandos, estados e códigos de saída.
- [ ] Testes do Doctor cobrem manifesto ausente, válido, inválido e erro de leitura.
- [ ] `inspect`, `validate` e `doctor` não modificam arquivos.
- [ ] `init` preserva manifesto existente, é seguro contra concorrência e não deixa arquivo parcial em falha.
- [ ] A normalização de `metadata.name` é determinística, não interativa e testada para casos comuns e limites.
- [ ] `make verify` passa.
- [ ] `git diff --check` passa nos arquivos alterados.
- [ ] README e changelog correspondem ao comportamento implementado.
- [ ] Revisões de arquitetura, documentação, compatibilidade e segurança foram registradas com evidências.

## Notes and Assumptions

- A validação atual do manifesto é estrutural; esta spec não expande o schema nem interpreta semanticamente todos os blocos opcionais.
- O Doctor atualmente possui oito checks. A implementação deverá atualizar os testes que dependem da contagem ou da ordem, sem remover verificações existentes.
- A criação do manifesto é a única operação com escrita de arquivo introduzida nesta spec.
- Nenhuma tarefa deve ser marcada como concluída com base apenas na redação do artefato.
