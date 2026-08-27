# Tasks 005 — Descoberta e registro de assets existentes

## Status

Concluída.

## Related Plan

- [`plan.md`](plan.md)

## Regras para o implementador

- Não iniciar antes da conclusão da Spec 004.
- Seguir a Spec 005, Plan 005 e ADR 001.
- Toda descoberta deve ser somente leitura até a etapa aprovada de persistência.
- Não executar, copiar, mover ou sobrescrever assets do usuário.
- Tasks bloqueadas permanecem desmarcadas e devem registrar evidência e motivo.

## Task List

### T1 — Fechar o contrato

- [x] Validar que a Spec 005, Plan 005 e ADR 001 foram aprovados.
- [x] Confirmar que a Spec 004 foi concluída e aprovada.
- [x] Mapear o schema `assets` para o modelo interno sem alterar `skills`.
- [x] Registrar os estados de resultado esperados: encontrado, adicionado, já
      registrado, ignorado, ambíguo/conflitante e ausente.

### T2 — Definir integração arquitetural

- [x] Definir a separação entre descoberta, classificação e persistência.
- [x] Definir a integração do `darp init` com descoberta automática.
- [x] Definir a integração do `darp doctor` para validar assets registrados.
- [x] Confirmar que o contrato não conflita com `skills` e seus caminhos atuais.
- [x] Registrar limites de segurança para caminhos, symlinks ignorados e raizes
      conhecidas.

### T3 — Implementar descoberta

- [x] Criar fixtures de repositório vazio, monorepo e famílias suportadas.
- [x] Implementar descoberta determinística e sem rede.
- [x] Garantir que a descoberta não execute arquivos nem comandos.
- [x] Reconhecer apenas os padrões por família definidos na Spec 005.
- [x] Ignorar e reportar symlinks sem segui-los.
- [x] Testar encontrados, ignorados, duplicados, ambíguos, conflitos semânticos
      e inválidos.

### T4 — Implementar registro

- [x] Atualizar `darp.yml` somente conforme o schema aprovado.
- [x] Preservar campos desconhecidos, registros existentes e assets físicos.
- [x] Adicionar apenas assets novos e informar assets já registrados.
- [x] Não persistir automaticamente candidatos ambíguos ou conflitantes.
- [x] Garantir idempotência e manter registrados assets fisicamente removidos.
- [x] Testar falhas de leitura, validação, escrita e substituição, incluindo a
      ausência de efeitos colaterais.

### T5 — Integrar e documentar

- [x] Integrar a descoberta automática ao `darp init`.
- [x] Atualizar `darp doctor` para warnings de assets ausentes.
- [x] Atualizar README e documentação do comando.
- [x] Atualizar `CHANGELOG.md` com o comportamento voltado ao usuário.
- [x] Executar `go test ./...` e `git diff --check`.
- [x] Revisar o diff contra Spec 005, Plan 005 e eventuais ADRs.

## Validation Checklist

- [x] Spec 004 concluída e aprovada.
- [x] Spec 005 e Plan 005 aprovados.
- [x] ADR 001 aprovado.
- [x] Contrato `assets` implementado com `path`, `family` e `type`.
- [x] Descoberta sem execução, rede ou alteração dos assets.
- [x] Registro seguro e idempotente.
- [x] Conflitos semânticos, classificação do `doctor` e limites de segurança
      totalmente testados.
- [x] Documentação e changelog atualizados.
- [x] Nenhum bloqueio pendente.

## Notas e bloqueios

Decisões registradas na Spec 005 e no ADR 001:

- descoberta automática no `darp init`;
- schema `assets` com `path`, `family` e `type`;
- registro somente por asset reconhecido;
- duplicatas e ambiguidades informadas sem escolha automática;
- symlinks ignorados;
- assets ausentes reportados como warning pelo `darp doctor`;
- combinações incompatíveis e symlinks registrados rejeitados pelo `darp doctor`;
- conflitos semânticos agrupados por categoria e nome normalizado;
- sem exclusões configuráveis na v1.

## Evidências

- Descoberta: encontrados, múltiplas famílias, ambíguos, ignorados (symlinks e
  raízes de família), conflitos semânticos, ausência de travessia de monorepo e
  determinismo em `internal/project/discovery/discovery_test.go`.
- Registro: idempotência, não duplicação em reexecução, assets ambíguos e
  conflitantes não persistidos, e rejeição de `assets` malformado antes do
  scaffold em `internal/project/init/service_test.go`.
- Persistência: preservação do `darp.yml` em falha de substituição, limpeza do
  temporário em falha de escrita e de rename, e rollback integral após falha
  tardia, comprovando ausência de efeitos colaterais.
- Doctor: asset válido passa, ausente gera `WARNING`, diretório no lugar de
  arquivo gera `FAIL`, tipo desconhecido e combinação incompatível geram `FAIL`,
  symlink registrado e travessia de caminho geram `FAIL` em
  `internal/project/doctor/service_test.go`.
