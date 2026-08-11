# Spec 005 — Descoberta e registro de assets existentes

## Status

Draft — planejamento inicial. Não iniciar a implementação antes da conclusão e
aprovação da [Spec 004](../004-init-governance-skills/spec.md).

## Related Specifications

- [Spec 004 — Scaffold de skills de governança](../004-init-governance-skills/spec.md)
- [Spec 002 — `darp doctor`](../002-doctor/spec.md)

## Contexto

O DARP CLI gerencia assets reutilizáveis para desenvolvimento assistido por IA,
mas um repositório existente pode conter instruções em locais convencionais de
diversas ferramentas, como `.github/`, `.claude/`, `CLAUDE.md`, `.codex/`,
`AGENTS.md`, `.cursor/` e `.gemini/`.

Atualmente, o `darp init` cria apenas a estrutura padrão do DARP. Ele não
descobre nem registra arquivos já presentes no repositório.

Esta demanda é separada da Spec 004: a Spec 004 trata do scaffold determinístico
de contratos DARP; esta spec trata da descoberta de conteúdo preexistente e da
sua representação no `darp.yml`.

## Problema

Ao inicializar ou atualizar um repositório existente, o usuário não recebe uma
visão registrada dos assets que já fazem parte das convenções de agentes e
ferramentas do projeto. Isso dificulta a composição, auditoria e futura gestão
desses assets pelo DARP.

## Objetivos preliminares

1. Detectar, de forma determinística e sem rede, arquivos e diretórios de
   instruções suportados pelo DARP.
2. Registrar os assets descobertos no `darp.yml` sem sobrescrever arquivos ou
   configurações existentes.
3. Preservar compatibilidade com projetos que não possuem esses arquivos.
4. Informar o que foi descoberto, ignorado ou considerado ambíguo.
5. Manter a descoberta independente de provedor, linguagem e framework.

## Não objetivos preliminares

- executar instruções, workflows ou assets descobertos;
- interpretar o conteúdo semântico dos arquivos;
- instalar ferramentas ou acessar a rede;
- mover, copiar, converter ou sobrescrever assets existentes;
- inferir a stack da aplicação;
- definir nesta spec o contrato final de todos os tipos de asset;
- alterar o `darp doctor` antes de o novo contrato ser aprovado;
- registrar automaticamente todo arquivo encontrado no repositório.

## Escopo preliminar

A primeira versão deverá avaliar somente padrões explicitamente suportados. A
lista abaixo é inicial e ainda precisa de aprovação:

| Família | Exemplos iniciais |
| --- | --- |
| GitHub | `.github/` e arquivos de instruções suportados dentro dele |
| Claude | `.claude/`, `CLAUDE.md` |
| Codex | `.codex/`, `AGENTS.md` |
| Cursor | `.cursor/` |
| Gemini | `.gemini/` |

A descoberta deverá produzir uma representação estruturada, com pelo menos o
caminho relativo e a família reconhecida. O formato exato permanece em aberto.

## Fluxos preliminares

### Inicialização com descoberta

1. O usuário executa o comando de inicialização com a opção de descoberta que
   vier a ser aprovada.
2. O CLI examina somente padrões suportados no diretório do projeto.
3. O CLI registra os resultados no `darp.yml` e informa um resumo.
4. Nenhum asset encontrado é executado ou modificado.

### Atualização de projeto registrado

O CLI preserva registros existentes, adiciona somente entradas novas conforme
as regras de identidade que serão definidas e informa conflitos ou arquivos
removidos sem apagá-los silenciosamente.

### Repositório sem assets suportados

O comando termina normalmente, não cria registros artificiais e informa que
nenhum asset suportado foi encontrado.

## Requisitos preliminares

- R1. A descoberta deve ser somente local e determinística.
- R2. Caminhos registrados devem ser relativos e permanecer dentro do projeto.
- R3. A descoberta não deve executar comandos, scripts ou acessar a rede.
- R4. Assets existentes devem ser preservados integralmente.
- R5. A atualização do `darp.yml` deve ser aditiva e segura, seguindo a Spec
  004 quando aplicável.
- R6. O resultado deve distinguir encontrado, já registrado, ignorado e
  ambíguo/conflitante.
- R7. O comportamento deve ser testável em repositórios vazios, monorepos e
  projetos com múltiplas famílias.

## Compatibilidade e arquitetura atual

O `darp.yml` atualmente possui uma seção `skills`, e o `darp doctor` exige que
as entradas dessa seção apontem para `.agents/skills/<nome>`. Portanto, esta
spec não pode registrar diretamente `.github/`, `.claude/`, `AGENTS.md` ou
outros caminhos externos como skills sem alterar esse contrato.

Antes da aprovação, deve ser decidido se os assets serão:

1. registrados em uma nova seção, como `assets:`;
2. representados por tipos específicos dentro de uma seção existente; ou
3. convertidos para skills DARP — hipótese que não deve ser presumida, pois
   pode alterar identidade, conteúdo e comportamento dos arquivos originais.

## Critérios de aceitação preliminares

- [ ] O contrato final de registro de assets está definido e documentado.
- [ ] A lista de padrões suportados e regras de exclusão está definida.
- [ ] A descoberta não executa, copia, move ou sobrescreve assets encontrados.
- [ ] O `darp.yml` preserva configurações e registros existentes.
- [ ] Reexecuções são determinísticas e não geram duplicatas.
- [ ] Conflitos, duplicidades, symlinks e arquivos removidos possuem comportamento
      definido e testado.
- [ ] O `darp doctor` foi atualizado ou explicitamente declarado compatível.
- [ ] Testes cobrem repositório vazio, assets únicos, múltiplas famílias,
      monorepo, caminhos inválidos e falhas de atualização.
- [ ] README, changelog e documentação do `darp init` refletem o comportamento.

## Perguntas e pendências

Resolver após a Spec 004:

1. A descoberta será automática ou opt-in, por exemplo `darp init --discover`?
2. Qual será o nome e o schema da seção no `darp.yml`?
3. O registro será por arquivo, diretório, família ou tipo semântico?
4. Quais arquivos dentro de cada diretório serão reconhecidos?
5. `AGENTS.md` e `CLAUDE.md` serão assets independentes ou aliases de família?
6. Como tratar nomes iguais, duplicatas e conflitos entre famílias?
7. Symlinks serão ignorados, registrados ou aceitos apenas dentro do projeto?
8. Arquivos removidos permanecem registrados, ficam ausentes ou são removidos
   mediante comando explícito?
9. O `doctor` validará existência e tipo dos assets registrados?
10. A descoberta atravessará limites de monorepo ou respeitará diretórios
    configuráveis?
11. Deve existir uma opção de exclusão para diretórios grandes ou gerados?
12. É necessário um ADR para schema e regras de identidade?

## Riscos

- registrar arquivos demais pode gerar ruído;
- diferentes ferramentas podem usar o mesmo nome para conceitos distintos;
- seguir symlinks pode expor arquivos fora do projeto;
- atualizar YAML preservando configuração local continua sendo sensível;
- descoberta automática pode surpreender usuários ao alterar `darp.yml`;
- tratar instruções como skills pode criar uma falsa promessa de execução.
