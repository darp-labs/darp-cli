# Tasks 004 — Scaffold de skills de governança no `darp init`

## Status

Ready for implementation after approval of Spec 004, Plan 004 and ADR 001.

## Related Plan

- [`plan.md`](plan.md)

## Regras para o implementador

- Marque cada checkbox somente após concluir a ação e sua validação.
- Preserve arquivos customizados e registre qualquer decisão de compatibilidade.
- Não introduza detecção de linguagem, framework ou tipologia de projeto.
- Tasks bloqueadas permanecem desmarcadas e devem registrar evidência e motivo.

## Task List

### T1 — Consolidar ativos canônicos

- [x] Criar `internal/project/init/assets/` como fonte canônica embutida para
      lifecycle, quality gates, workflow e as quatro skills. Validar com teste
      de presença e conteúdo mínimo.
- [x] Garantir que os conteúdos são determinísticos e funcionam sem rede.
- [x] Garantir que cada skill contém as seções e limites da Spec 003.
- [x] Remover instruções específicas de Go, Java, Spring, Python, FastAPI,
      providers ou ferramentas obrigatórias.
- [x] Garantir que o workflow não passe a executar skills automaticamente.

### T2 — Expandir a estrutura criada pelo init

- [x] Incluir os diretórios `architecture`, `testing` e `release` na estrutura
      criada pelo serviço de inicialização.
- [x] Criar `SKILL.md` ausente para as quatro skills.
- [x] Criar lifecycle e quality gates completos somente quando os arquivos
      estiverem ausentes.
- [x] Manter a criação de workflow e templates compatível com as specs
      anteriores.

### T3 — Atualizar `darp.yml` com segurança

- [x] Fazer projetos novos registrarem as quatro skills nos caminhos canônicos.
- [x] Implementar a estratégia aditiva do ADR 001, editando somente `skills`.
- [x] Preservar valores existentes, skills adicionais e campos desconhecidos.
- [x] Não sobrescrever `darp.yml` quando o YAML for inválido ou a atualização
      segura não puder ser concluída.
- [x] Cobrir entradas duplicadas, valores customizados, configuração parcial e
      falha na substituição do arquivo.

### T4 — Preservar ativos e idempotência

- [x] Garantir que arquivos existentes nunca sejam substituídos integralmente;
      tratar a atualização aditiva de `darp.yml` como exceção controlada.
- [x] Garantir que `SKILL.md` customizado seja preservado integralmente.
- [x] Garantir que lifecycle e quality gates customizados sejam preservados.
- [x] Atualizar somente os placeholders históricos exatos para os ativos
      canônicos equivalentes.
- [x] Preservar arquivos que tenham qualquer conteúdo diferente do placeholder,
      mesmo quando pareçam incompletos.
- [x] Garantir que a segunda execução não remova arquivos do usuário.
- [x] Garantir que uma execução parcial possa ser reparada na execução seguinte.

### T5 — Criar testes agnósticos de projeto

- [x] Testar projeto novo sem código.
- [x] Testar fixture Java com Spring.
- [x] Testar fixture Python com FastAPI.
- [x] Testar fixture Go.
- [x] Testar monorepo com múltiplas linguagens.
- [x] Testar projeto parcialmente inicializado.
- [x] Testar YAML inválido, campos desconhecidos e skills adicionais.
- [x] Testar YAML semanticamente incompleto e chaves duplicadas em `skills`.
- [x] Testar upgrade dos placeholders históricos de quality gates, documentation
      e workflow.
- [x] Testar preservação de cada placeholder com uma alteração mínima de
      conteúdo.
- [x] Testar falhas de filesystem sem remoção de arquivos existentes.
- [x] Testar idempotência e conteúdo mínimo de todos os contratos.

### T6 — Validar integração e documentação

- [x] Executar `go test ./...`.
- [x] Executar `darp doctor` em cada fixture válida.
- [x] Confirmar que nenhum teste depende de detecção de stack para obter
      sucesso.
- [x] Confirmar que nenhum comando do projeto-alvo é executado pelo `init`.
- [x] Atualizar README e documentação da inicialização com a nova estrutura.
- [x] Sincronizar os contratos versionados da raiz com os ativos canônicos e
      registrar a mudança de bootstrap em `CHANGELOG.md`.
- [x] Atualizar o `darp.yml` deste repositório para registrar as quatro skills
      de governança, mantendo `security-review` independente.
- [x] Confirmar que `.agents/skills/security-review/` não foi alterada.
- [x] Executar `git diff --check`.
- [x] Revisar o diff final contra Spec 004, Plan 004 e todas as tasks.

## Validation Checklist

- [x] Spec 004 aprovada.
- [x] Plan 004 aprovado.
- [x] Todas as tasks T1–T6 concluídas.
- [x] Projetos novos recebem as quatro skills.
- [x] Projetos existentes são reparados sem sobrescrita destrutiva.
- [x] Configuração inválida existente é reportada e permanece intacta.
- [x] Fixtures de múltiplas stacks passam no `darp doctor`.
- [x] Testes existentes passam.
- [x] Nenhum bloqueio pendente.

## Notas e bloqueios

Registrar decisões sobre preservação de YAML, conteúdo canônico e compatibilidade
com projetos já inicializados. Uma necessidade de executar skills, workflows ou
detectar stack deve gerar uma nova especificação, não uma alteração silenciosa
nesta implementação.
