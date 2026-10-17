# Spec 009 — Sincronização explícita do inventário de assets

## Status

Approved — especificação aprovada formalmente pelo responsável do projeto em
2026-10-14. Implementação autorizada somente conforme esta Spec e o plano e as
tasks aprovados.

## Contexto

A Spec 005 define a descoberta local de assets durante `darp init`. Ela registra
assets externos na seção `assets` do `darp.yml`, mantém entradas já registradas
e determina que assets removidos do workspace permaneçam no manifesto até uma
ação explícita futura. A mesma Spec limita os padrões descobertos às raízes de
famílias externas listadas no contrato; `.agents/skills/` não está entre elas.

A Spec 004 reserva `skills` para caminhos `.agents/skills/<nome>` e define as
quatro skills de governança distribuídas pelo `init`: `documentation`,
`architecture`, `testing` e `release`. O Doctor também valida as entradas de
`skills` contra esse caminho e examina cada diretório de skill diretamente sob
`.agents/skills/`.

Foi decidido pelo responsável do projeto que skills válidas em
`.agents/skills/` devem ser descobertas e registradas em `skills`, mantendo as
quatro skills de governança conhecidas. Também foi decidido que erros de leitura
durante a descoberta devem gerar aviso e ser registrados como ignorados. Essas
decisões alteram e ampliam os contratos de descoberta existentes e precisam ser
formalizadas por esta Spec aprovada antes da implementação.

O ciclo de uso desejado é executar `darp init` uma vez para o bootstrap inicial
e usar `darp sync` nas atualizações posteriores do inventário. Isso não altera a
compatibilidade do `init`: reexecuções existentes continuam aditivas e seguras,
mas não passam a sincronizar nem remover registros.

A Spec 008 — Harness Lifecycle — está aprovada, mas ainda não concluída: seus
critérios e tasks continuam desmarcados, e o código atual não expõe os
subcomandos Harness definidos nela. A operação definida aqui é independente do
manifesto `.darp/harness.yaml` e não deve alterar `darp init`, Harness ou
`darp doctor`.

## Problema

O comportamento aditivo do `darp init` não reconcilia `darp.yml` com o estado
atual do workspace. Um asset removido continua registrado; skills adicionadas
em `.agents/skills/` não são descobertas pela Spec 005 atual. Reexecutar `init`
não resolve essas diferenças, e não existe operação explícita de sincronização.

Além disso, o manifesto legado não contém provenance que distinga entradas
criadas automaticamente de entradas mantidas manualmente. A sincronização não
pode inferir ownership apenas pela presença de uma entrada e removê-la sem uma
regra aprovada.

## Objetivos

1. Definir uma operação explícita que compare o inventário de `darp.yml` com os
   assets reconhecidos no workspace.
2. Incluir skills válidas em `.agents/skills/` na seção `skills`, sem alterar os
   formatos originais das skills.
3. Permitir que o inventário gerenciado seja atualizado após inclusão e remoção
   de assets, sem modificar `darp init` implicitamente.
4. Preservar campos desconhecidos, metadados de assets e configurações que não
   pertencem à descoberta automática.
5. Definir resultados determinísticos para inclusão, manutenção, remoção,
   ignorados, avisos e falhas.
6. Não iniciar, executar, instalar, mover ou reescrever assets.

## Não objetivos

- Alterar retroativamente o comportamento aditivo de `darp init`.
- Implementar execução, instalação, ativação, validação semântica ou segurança
  de conteúdo dos assets.
- Sincronizar `.darp/harness.yaml` ou modificar os comandos Harness da Spec 008.
- Descobrir skills em diretórios aninhados sob `.agents/skills/`; o Doctor e a
  estrutura vigente tratam cada filho direto como uma skill.
- Resolver automaticamente conflitos semânticos ou identidades ambíguas.
- Acessar rede, consultar Registry, executar comandos ou seguir symlinks.
- Remover entradas sem regra explícita de ownership aprovada.

## Alternativas de interface

### Alternativa A — comando `darp sync` (recomendada)

Cria uma operação própria para sincronização de inventário. Mantém a distinção
entre `init` (bootstrap aditivo), `sync` (reconciliação explícita), `doctor`
(diagnóstico) e os comandos Harness da Spec 008. Também deixa espaço para
futuras opções de simulação e escopo sem acumular responsabilidades em `init`.

### Alternativa B — opção em `darp init`

Acrescenta uma opção como `darp init --sync`. Evita um comando novo, mas associa
reconciliação destrutiva ao bootstrap e torna menos clara a diferença entre
inicializar/reparar contratos DARP e reconciliar inventários de assets.

**Decisão aprovada:** adotar `darp sync`. A invocação explícita reduz efeitos
surpresa, deixa clara a responsabilidade de atualização após o bootstrap e
preserva o contrato atual de `darp init`.

## Contrato funcional

As decisões de produto sobre interface, ownership legado, provenance, skills
de governança e resultados parciais foram aprovadas e são normativas nesta Spec.

### Interface

```text
darp sync [--dry-run]
```

- `darp sync` calcula e aplica a reconciliação permitida pela regra de
  ownership definida nesta Spec.
- `darp sync --dry-run` calcula e apresenta o mesmo relatório sem modificar
  arquivos.
- A ausência de opções preserva o comportamento de todos os comandos
  existentes, em especial `darp init`.
- Argumentos desconhecidos retornam erro de uso e não modificam arquivos.

`--dry-run` é recomendado para que o usuário possa revisar adições e remoções.
Não se propõe exigir `--apply` por padrão: `darp sync` já é uma ação explícita.
Registros que não correspondam aos contratos de discovery não são adotados nem
removidos automaticamente; devem ser preservados e reportados.

### Fontes de descoberta

1. Os padrões externos de arquivo da Spec 005 continuam reconhecidos e
   registrados em `assets` com `path`, `family` e `type`.
2. Skill DARP é reconhecida em `.agents/skills/<nome>/SKILL.md`, onde `<nome>`
   é um único segmento de diretório válido. O diretório deve ser filho direto
   de `.agents/skills/`; `SKILL.md` deve ser arquivo regular.
3. Uma skill encontrada nessa estrutura é registrada como
   `skills.<nome>: .agents/skills/<nome>`.
4. As quatro skills de governança continuam válidas e registradas nos caminhos
   canônicos definidos pela Spec 004. A sincronização não substitui seu
   conteúdo. A presença física determina a entrada espelhada, sem exceção de
   proteção para as quatro skills.
5. Diretórios ausentes são opcionais e não geram erro. Arquivos que não
   correspondem aos padrões documentados não são registrados.
6. Nomes de skill duplicados após normalização, conflitos semânticos, symlinks,
   arquivos não regulares e estruturas incompletas não são registrados
   automaticamente; são reportados conforme o formato final de diagnóstico.

### Planejamento da reconciliação

A sincronização calcula um plano antes de gravar:

- **adicionado:** identidade descoberta que não está no inventário;
- **mantido:** identidade descoberta já registrada; seus campos adicionais e
  valores existentes permanecem intactos;
- **removido:** entrada comprovadamente gerenciada pela descoberta, mas ausente
  do workspace, conforme a regra de ownership desta Spec;
- **ignorado:** caminho ou subárvore cuja leitura falhou, symlink, formato
  incompatível, conflito ou candidato que não pode ser classificado com
  segurança;
- **aviso:** resultado parcial que não impediu uma sincronização segura;
- **falha:** manifesto inválido, erro de escrita ou estado em que a segurança
  da atualização não pode ser garantida.

O resultado deve ser ordenado deterministicamente por tipo de operação e
identidade normalizada. A mesma entrada e o mesmo estado do workspace devem
produzir o mesmo plano.

### Ownership, preservação e gravação

- Atualizar somente as partes de `skills` e `assets` definidas como propriedade
  da sincronização.
- Preservar campos de topo desconhecidos, configurações manuais, valores
  existentes de skills mantidas e chaves adicionais em registros de assets.
- Não substituir valor manual existente por caminho calculado sem uma regra de
  conflito explicitamente aprovada.
- Não remover dados sem ownership conhecido pela regra de sincronização.
- Guardar o ownership gerenciado no próprio `darp.yml`, em metadados namespaced
  sob `darp.sync`. O ledger deve identificar skills por nome e caminho e assets
  por `(path, family, type)`; o ledger registra ownership e não substitui os
  inventários `skills` e `assets`.
- Na primeira sincronização, adotar como gerenciadas as entradas legadas que
  correspondem aos contratos de discovery: skills com nome/caminho canônico
  `.agents/skills/<nome>` e assets cuja combinação de caminho, família e tipo
  corresponde a uma regra da Spec 005. A adoção inclui entradas cujo arquivo ou
  diretório já está ausente; elas serão removidas no mesmo plano se o escopo
  correspondente pôde ser lido com sucesso.
- Preservar entradas que não correspondam a uma regra reconhecida e não constem
  no ledger. Nunca inferir que uma entrada arbitrária é gerenciada apenas por
  seu nome.
- Ao remover uma entrada gerenciada, reportar sua identidade. Entradas adicionais
  em um registro de asset mantido são preservadas; quando o registro inteiro é
  removido, seus campos adicionais saem com ele e devem aparecer na prévia.
- Validar todo o plano antes da escrita. Manifesto inválido, YAML duplicado,
  estrutura não editável ou conflito que afete a segurança deve deixar
  `darp.yml` byte a byte inalterado.
- Em escrita aplicada, usar substituição atômica ou mecanismo equivalente que
  preserve o original em falha. Nunca escrever parcialmente o manifesto.
- Não escrever nos arquivos originais dos assets.
- Segunda execução sem mudanças não reformatará nem alterará o manifesto.

### Erros de leitura e sincronização parcial

- Erro de leitura de uma raiz ou candidato deve gerar aviso e registro como
  ignorado, conforme decisão já tomada pelo responsável do projeto.
- Uma raiz/subárvore que não pôde ser lida é estado desconhecido, não prova de
  remoção. Entradas sob esse escopo não podem ser removidas nessa execução.
- Outras raízes legíveis podem ser sincronizadas se o manifesto e a gravação
  forem seguros; o relatório deve identificar claramente o resultado parcial.
- Erro ao ler ou validar `darp.yml` impede qualquer escrita.
- Erro ao persistir deve preservar o arquivo original e retornar falha.

### Códigos de saída

- `0`: sincronização concluída, com ou sem avisos por itens ignorados;
- `1`: erro de uso, manifesto inválido, conflito que impede atualização segura
  ou erro operacional de escrita.

Isso segue a convenção atual: o Doctor retorna sucesso quando há warnings sem
`FAIL`.

## Exemplos

Antes:

```yaml
skills:
  documentation: .agents/skills/documentation
  architecture: .agents/skills/architecture
  testing: .agents/skills/testing
  release: .agents/skills/release
assets:
  - path: .github/skills/old/SKILL.md
    family: github
    type: skill
    owner: platform-team
```

Workspace atual contém as quatro skills de governança, duas skills válidas em
`.agents/skills/` e não contém mais `.github/skills/old/SKILL.md`. Após sync, o
resultado esperado é registrar as duas skills novas e retirar o asset GitHub
porque a entrada legada corresponde a um padrão conhecido de discovery e o
workspace foi lido sem erro. O campo adicional `owner` é preservado enquanto a
entrada for mantida; ao remover o registro inteiro, seus campos extras também
são removidos e identificados na prévia.

Exemplo de prévia:

```text
Would add skill: custom-review -> .agents/skills/custom-review
Would remove asset: .github/skills/old/SKILL.md (github/skill)
Kept 4 governance skills
Warnings: 0
```

Formato de metadados de ownership no `darp.yml` após sincronização:

```yaml
darp:
  sync:
    managed:
      skills:
        - name: custom-review
          path: .agents/skills/custom-review
      assets: []
```

O ledger deve conter apenas identidades atualmente gerenciadas e não deve
funcionar como um registro de identidades removidas. Leitores atuais que ignoram
campos desconhecidos devem continuar funcionando.

As mensagens finais da CLI serão detalhadas no plano, sem alterar as categorias
de resultado e os códigos definidos nesta Spec.

## Requisitos

- **R1 — Explícito:** apenas `darp sync` inicia sincronização. `darp init`,
  `darp doctor` e comandos Harness não adquirem efeitos de sincronização.
- **R2 — Descoberta:** aplicar padrões da Spec 005 e incluir
  `.agents/skills/<nome>/SKILL.md` conforme o contrato desta Spec.
- **R3 — Inclusão:** adicionar entradas novas ao campo correspondente sem
  duplicar identidades.
- **R4 — Remoção segura:** remover somente entradas comprovadamente gerenciadas e
  ausentes; na migração inicial, adotar as entradas legadas que correspondam
  exatamente aos contratos de discovery; preservar e reportar entradas fora
  desses contratos.
- **R5 — Preservação:** preservar configuração manual, campos desconhecidos,
  metadados de entradas mantidas e arquivos físicos dos assets.
- **R6 — Idempotência:** uma segunda execução sem mudança do workspace não
  altera bytes do manifesto.
- **R7 — Determinismo:** descoberta, plano, mensagens e ordenação são estáveis
  para o mesmo estado observável.
- **R8 — Leitura parcial:** erros de leitura geram aviso/ignorado e impedem
  remoções no escopo não lido.
- **R9 — Falha segura:** manifesto inválido, estrutura ambígua ou erro de escrita
  não causa perda nem gravação parcial.
- **R10 — Prévia:** `--dry-run` mostra o plano sem alterar o filesystem.
- **R11 — Relatório:** informar contagem e identidade de itens adicionados,
  removidos, mantidos, ignorados e avisos.
- **R12 — Compatibilidade:** manter os contratos de `darp init`, `darp doctor`,
  `darp harness inspect`, `darp harness validate` e `darp harness init` da Spec
  008; Harness permanece separado em `.darp/harness.yaml`.

## Fluxos

### Sincronização com alterações

1. Usuário executa `darp sync` ou `darp sync --dry-run`.
2. CLI lê e valida `darp.yml`, descobre assets e coleta erros/ignorados.
3. CLI constrói e apresenta plano determinístico de alterações.
4. No modo aplicado, CLI grava somente se o plano for seguro pela regra de
   ownership aprovada.
5. CLI informa adições, remoções, itens mantidos, ignorados e avisos.

### Manifesto inválido ou leitura parcial

1. Manifesto inválido impede escrita e mantém seus bytes originais.
2. Erro de leitura de asset gera aviso e item ignorado.
3. Entradas do escopo não lido não são consideradas removidas.
4. A CLI retorna código conforme a classificação aprovada de sucesso parcial.

### Sem mudanças

1. A CLI descobre os mesmos assets já representados.
2. Nenhum byte do manifesto é regravado.
3. O relatório informa que não houve alterações.

## Critérios de aceitação

- [ ] Um `darp sync` registra duas novas skills válidas em `.agents/skills/`.
- [ ] Um asset antes descoberto e removido do workspace deixa o inventário
      gerenciado, sem apagar registro manual ou de ownership ambíguo.
- [ ] Skills das quatro famílias existentes e `.agents/skills/` são classificadas
      no destino correto (`assets` ou `skills`).
- [ ] Duas skills com identidade conflitante não são escolhidas automaticamente.
- [ ] Raízes opcionais ausentes não geram falha.
- [ ] Estruturas incompletas e nomes inválidos são ignorados e reportados.
- [ ] Erros de leitura geram aviso e impedem remoções no escopo afetado.
- [ ] Campos extras de topo e campos adicionais de registros mantidos são
      preservados.
- [ ] Manifesto inválido ou YAML duplicado não é alterado.
- [ ] `--dry-run` não modifica arquivos e exibe o mesmo plano da execução
      aplicada.
- [ ] Segunda sincronização sem mudanças mantém o arquivo byte a byte idêntico.
- [ ] Falha de persistência não deixa manifesto parcial nem modifica assets.
- [ ] `darp init`, `darp doctor` e os comandos da Spec 008 preservam seus
      comportamentos aprovados.
- [ ] Testes unitários cobrem parser, propriedade/ownership, plano de diff,
      ordenação, falhas e idempotência; testes de CLI cobrem saída e códigos.
- [ ] README, CHANGELOG e roadmap são atualizados conforme a prática vigente.

## Compatibilidade e riscos

- Como `darp.yml` legado não registra provenance, a primeira sync adota
  registros que correspondam exatamente aos contratos reconhecidos. Isso pode
  remover uma entrada manual que use a mesma identidade; `--dry-run` apresenta
  essa remoção antes da escrita.
- `skills` é um mapa nome→caminho. Entradas fora do formato canônico de skill
  ficam fora do ledger e são preservadas.
- `assets` é uma lista de mapas com identidade `(path, family, type)`. Campos
  adicionais de entradas mantidas são preservados; campos de um registro
  removido acompanham a remoção e aparecem na prévia.
- Uma raiz inacessível não prova que seus assets foram removidos. A sync não
  remove entradas no escopo não lido.
- O ledger deve coexistir com campos desconhecidos sob `darp` e não pode
  invalidar leitores atuais que ignorem essa extensão.

## Dependências

- Spec 005 — Descoberta e registro de assets existentes e ADR 001.
- Spec 004 — Scaffold de skills de governança.
- Spec 002 — Contrato do Doctor para `skills` e assets.
- Spec 008 — Harness Lifecycle, para compatibilidade de CLI; a Spec 008 está
  aprovada, mas sua implementação ainda não está concluída no estado atual do
  repositório.
- `internal/project/discovery`, `internal/project/asset`,
  `internal/project/init` e `internal/project/doctor`.

## Decisões de produto incorporadas

As decisões abaixo foram aprovadas pelo responsável do projeto e fazem parte
do contrato desta Spec:

1. **Interface:** comando explícito `darp sync`, com `--dry-run` opcional. O
   `darp init` destina-se ao bootstrap inicial; sincronizações posteriores são
   explícitas.
2. **Ownership legado:** adotar como gerenciadas as entradas legadas que
   correspondam aos contratos de discovery, inclusive entradas ausentes que
   possam ser classificadas pelo caminho, família e tipo. Entradas fora desses
   contratos permanecem preservadas.
3. **Provenance:** guardar ownership em metadados no próprio `darp.yml`, sob o
   namespace `darp.sync.managed`.
4. **Skills de governança:** não são protegidas contra remoção; `sync` espelha o
   filesystem também para as quatro skills conhecidas.
5. **Leitura parcial:** erros de leitura geram warning/ignorado; as entradas do
   escopo não lido não são removidas. Se as demais alterações seguras forem
   aplicadas, o código de saída é `0`.
6. **Roadmap:** não atualizar como parte desta feature, conforme decisão
   aprovada pelo responsável do projeto.

## Próximos passos SDD

O próximo passo é revisar e aprovar o plano e as tasks em
`.spec/specs/009-asset-sync/` antes da implementação. Não é necessário ADR neste
momento; um ADR será criado somente se a implementação revelar uma decisão
arquitetural separada que não esteja resolvida por esta Spec.
