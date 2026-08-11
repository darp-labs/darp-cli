# Plan 005 — Descoberta e registro de assets existentes

## Status

Draft inicial — não iniciar antes da conclusão da Spec 004 e da resolução das
pendências da Spec 005.

## Related Specifications

- [`spec.md`](spec.md)
- [`../004-init-governance-skills/spec.md`](../004-init-governance-skills/spec.md)
- [`../002-doctor/spec.md`](../002-doctor/spec.md)

## Estratégia preliminar

Separar descoberta, classificação e persistência. A descoberta deve produzir
resultados estruturados sem alterar o projeto; uma etapa posterior decide como
registrá-los no `darp.yml`. O comportamento final, o schema de persistência e as
regras de identidade dependem das decisões abertas na spec.

Nenhum código deve ser implementado enquanto a seção de assets do `darp.yml`,
as famílias suportadas e o contrato do `darp doctor` não estiverem aprovados.

## Dependências e ordem

1. Concluir e aprovar a Spec 004.
2. Resolver as perguntas abertas da Spec 005.
3. Definir o schema de registro e, se necessário, criar um ADR.
4. Definir a API de descoberta e as regras de segurança para caminhos e
   symlinks.
5. Implementar descoberta sem efeitos colaterais.
6. Implementar persistência aditiva no `darp.yml`.
7. Atualizar `darp doctor`, documentação, changelog e testes conforme o escopo
   aprovado.

## Áreas potencialmente afetadas

- `internal/project/init/`;
- novo pacote de descoberta, se a separação de responsabilidades exigir;
- `internal/project/doctor/`, somente se o novo contrato exigir validação;
- testes de init, descoberta e doctor;
- `README.md` e documentação do comando;
- `CHANGELOG.md`;
- contrato e templates de `darp.yml`.

## Áreas protegidas preliminares

- assets encontrados no projeto do usuário;
- skills específicas existentes, incluindo `security-review`;
- execução de workflows e skills;
- código da aplicação e arquivos fora dos padrões aprovados.

## Marcos preliminares

### M1 — Contrato de assets

Definir schema, identidade, famílias suportadas, regras de exclusão e relação
com `skills`.

### M2 — Descoberta segura

Implementar varredura local, determinística e sem execução, com tratamento de
symlinks, limites de diretório e caminhos excluídos.

### M3 — Registro aditivo

Persistir somente mudanças aprovadas, preservando configuração e registros
existentes.

### M4 — Validação

Testar repositórios vazios, múltiplas famílias, monorepos, conflitos,
duplicidades, reexecução e falhas de filesystem.

### M5 — Integração e documentação

Integrar ao comando aprovado, atualizar `doctor` se necessário e documentar
comportamento, limites e impacto de compatibilidade.

## Estratégia de validação preliminar

- testes unitários da descoberta sem filesystem real quando possível;
- testes de integração com fixtures de cada família suportada;
- verificação de que nenhum comando ou rede é acionado;
- comparação do `darp.yml` antes/depois para comprovar preservação;
- segunda execução para comprovar idempotência;
- `darp doctor` antes e depois, conforme o contrato final;
- `go test ./...`;
- `git diff --check`.

## Recuperação preliminar

Falha na descoberta não deve modificar arquivos. Falha na persistência deve
preservar o `darp.yml` original e informar a causa. A estratégia final deverá
reutilizar as garantias de atualização segura decididas na Spec 004.

## Premissas

- A Spec 004 será concluída antes desta feature.
- Assets externos não serão executados nem convertidos automaticamente.
- A descoberta será agnóstica a linguagem, framework e provedor.
- Decisões ainda abertas não devem ser escondidas na implementação.
