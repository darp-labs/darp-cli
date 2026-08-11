# ADR 001 — Fonte dos ativos e atualização segura do `darp.yml`

## Status

Proposed — requires approval with Spec 004 and Plan 004.

## Contexto

O `darp init` precisa distribuir contratos de governança sem depender do
diretório do repositório-fonte e precisa adicionar as quatro skills a um
`darp.yml` existente sem perder configuração local. Uma escrita direta não
permite preservar o arquivo original quando a atualização falha.

## Decisão

1. Os ativos distribuídos serão arquivos versionados em
   `internal/project/init/assets/`, incorporados ao binário com `go:embed`.
   Esse diretório será a fonte canônica do conteúdo usado pelo `init`:
   lifecycle, quality gates, workflow e os quatro `SKILL.md`. Os contratos
   equivalentes mantidos na raiz devem ser sincronizados quando o conteúdo
   canônico mudar e terão testes de conteúdo mínimo.
2. O `darp.yml` será analisado com `gopkg.in/yaml.v3` antes de qualquer escrita.
   A atualização só gerencia as quatro chaves sob `skills`.
3. Chaves existentes, valores customizados, skills adicionais, campos
   desconhecidos e comentários que possam ser preservados pelo formato serão
   mantidos. Chaves duplicadas no mapa `skills`, estrutura não suportada para
   edição segura ou YAML inválido são erro: o arquivo original e os demais
   ativos não serão alterados.
4. A alteração será escrita por arquivo temporário no mesmo diretório, seguida
   de substituição controlada, com a abstração de filesystem estendida para
   leitura e renomeação. Se a substituição falhar, o original será mantido.
5. O `init` não corrige campos obrigatórios ausentes ou inválidos além das
   quatro entradas de skills. Nesses casos ele reporta o problema e não promete
   que `darp doctor` passará.
6. Placeholders históricos exatos gerados por versões anteriores do `init`
   podem ser atualizados para os ativos canônicos equivalentes. A comparação é
   byte a byte e limitada à lista da Spec 004; qualquer divergência é tratada
   como conteúdo customizado e preservada.
7. O `darp.yml` versionado no repositório do DARP CLI será atualizado para
   registrar as quatro skills de governança. A skill `security-review` continua
   independente e não será adicionada por esta spec.

## Consequências

- O conteúdo embarcado é determinístico e funciona sem rede.
- A preservação de semântica YAML é obrigatória; preservação byte a byte de
  comentários e formatação é garantida apenas quando o formato permite edição
  segura.
- A atualização de configuração passa a exigir leitura, arquivo temporário e
  renomeação testáveis no filesystem de teste.
- Um projeto novo continua recebendo um `darp.yml` completo e válido; um
  projeto existente com configuração inválida não é silenciosamente reparado.
- Projetos antigos podem receber upgrade dos placeholders oficiais sem que
  arquivos customizados sejam sobrescritos.

## Alternativas rejeitadas

- Reescrever todo o YAML com `yaml.Marshal`, pois isso pode remover comentários,
  formatação e campos que não façam parte do modelo conhecido.
- Ler contratos do diretório de trabalho, pois isso faria o binário depender do
  repositório-fonte ou da rede.
