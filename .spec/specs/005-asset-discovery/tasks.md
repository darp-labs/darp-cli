# Tasks 005 — Descoberta e registro de assets existentes

## Status

Draft inicial — manter todas as tasks desmarcadas até a aprovação da Spec 005.

## Related Plan

- [`plan.md`](plan.md)

## Regras para o implementador

- Não iniciar antes da conclusão da Spec 004.
- Não implementar decisões que ainda estejam listadas como pendências.
- Toda descoberta deve ser somente leitura até a etapa aprovada de persistência.
- Não executar, copiar, mover ou sobrescrever assets do usuário.
- Tasks bloqueadas permanecem desmarcadas e devem registrar evidência e motivo.

## Task List

### T1 — Fechar o contrato

- [ ] Definir se a descoberta será automática ou opt-in.
- [ ] Definir o schema de registro no `darp.yml`.
- [ ] Definir a identidade de um asset e as regras de duplicidade.
- [ ] Definir famílias, padrões, arquivos reconhecidos e exclusões.
- [ ] Decidir o tratamento de symlinks, arquivos removidos e monorepos.
- [ ] Decidir se será necessário um ADR.

### T2 — Definir integração arquitetural

- [ ] Definir a separação entre descoberta, classificação e persistência.
- [ ] Definir se `darp doctor` validará assets registrados.
- [ ] Confirmar que o contrato não conflita com `skills` e seus caminhos atuais.
- [ ] Registrar limites de segurança para caminhos, symlinks e diretórios
      excluídos.

### T3 — Implementar descoberta

- [ ] Criar fixtures de repositório vazio, monorepo e famílias suportadas.
- [ ] Implementar descoberta determinística e sem rede.
- [ ] Garantir que a descoberta não execute arquivos nem comandos.
- [ ] Testar encontrados, ignorados, duplicados, ambíguos e inválidos.

### T4 — Implementar registro

- [ ] Atualizar `darp.yml` somente conforme o schema aprovado.
- [ ] Preservar campos desconhecidos, registros existentes e assets físicos.
- [ ] Garantir idempotência e comportamento definido para remoções.
- [ ] Testar falhas de leitura, validação, escrita e substituição.

### T5 — Integrar e documentar

- [ ] Integrar a descoberta ao fluxo aprovado do `darp init`.
- [ ] Atualizar `darp doctor`, se exigido pelo contrato aprovado.
- [ ] Atualizar README e documentação do comando.
- [ ] Atualizar `CHANGELOG.md` com o comportamento voltado ao usuário.
- [ ] Executar `go test ./...` e `git diff --check`.
- [ ] Revisar o diff contra Spec 005, Plan 005 e eventuais ADRs.

## Validation Checklist

- [ ] Spec 004 concluída e aprovada.
- [ ] Spec 005 e Plan 005 aprovados.
- [ ] Todas as pendências de contrato resolvidas.
- [ ] Descoberta sem execução, rede ou alteração dos assets.
- [ ] Registro seguro e idempotente.
- [ ] Conflitos e limites de segurança testados.
- [ ] Documentação e changelog atualizados.
- [ ] Nenhum bloqueio pendente.

## Notas e bloqueios

Registrar aqui as decisões tomadas após a conclusão da Spec 004, especialmente
o schema do `darp.yml`, o modo de acionamento, as famílias suportadas e a
compatibilidade com `darp doctor`.
