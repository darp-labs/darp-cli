# Plan 005 — Descoberta e registro de assets existentes

## Status

Aprovado e implementado — contrato resolvido na Spec 005 e no
[ADR 001](adr/001-asset-registry-schema-and-identity.md).

## Related Specifications

- [`spec.md`](spec.md)
- [`../004-init-governance-skills/spec.md`](../004-init-governance-skills/spec.md)
- [`../002-doctor/spec.md`](../002-doctor/spec.md)
- [`adr/001-asset-registry-schema-and-identity.md`](adr/001-asset-registry-schema-and-identity.md)

## Estratégia

Separar descoberta, classificação e persistência. A descoberta deve produzir
resultados estruturados sem alterar o projeto; uma etapa posterior decide como
registrá-los no `darp.yml`.

O schema de persistencia aprovado e `assets`, com entradas contendo `path`,
`family` e `type`. O `darp init` executara descoberta automaticamente e fara
persistencia aditiva apenas quando o `darp.yml` puder ser editado com seguranca.

## Dependências e ordem

1. Concluir e aprovar a Spec 004.
2. Aprovar a Spec 005, Plan 005, Tasks 005 e ADR 001.
3. Definir a API de descoberta e as regras de segurança para caminhos e
   symlinks.
4. Implementar descoberta sem efeitos colaterais.
5. Implementar persistência aditiva no `darp.yml`.
6. Atualizar `darp doctor`, documentação, changelog e testes conforme o escopo
   aprovado.

## Áreas potencialmente afetadas

- `internal/project/init/`;
- novo pacote de descoberta, se a separação de responsabilidades exigir;
- `internal/project/doctor/`, somente se o novo contrato exigir validação;
- testes de init, descoberta e doctor;
- `README.md` e documentação do comando;
- `CHANGELOG.md`;
- contrato e templates de `darp.yml`.

## Áreas protegidas

- assets encontrados no projeto do usuário;
- skills específicas existentes, incluindo `security-review`;
- execução de workflows e skills;
- código da aplicação e arquivos fora dos padrões aprovados.

## Marcos

### M1 — Contrato de assets

Implementar o contrato aprovado: secao `assets`, identidade `(path, family,
type)`, familias suportadas, ausencia de exclusoes configuraveis na v1 e
separacao de `skills`.

### M2 — Descoberta segura

Implementar varredura local, determinística e sem execução, com tratamento de
symlinks ignorados, limites de diretorio e raizes conhecidas.

### M3 — Registro aditivo

Persistir somente mudanças aprovadas, preservando configuração e registros
existentes. Informar `already registered` para entradas ja presentes e reportar
duplicidade/ambiguidade sem escolha automatica.

### M4 — Validação

Testar repositórios vazios, múltiplas famílias, monorepos, conflitos,
duplicidades, reexecução e falhas de filesystem.

### M5 — Integração e documentação

Integrar ao comando aprovado, atualizar `doctor` se necessário e documentar
comportamento, limites e impacto de compatibilidade.

## Estratégia de validação

- testes unitários da descoberta sem filesystem real quando possível;
- testes de integração com fixtures de cada família suportada;
- verificação de que nenhum comando ou rede é acionado;
- comparação do `darp.yml` antes/depois para comprovar preservação;
- segunda execução para comprovar idempotência;
- `darp doctor` antes e depois, conforme o contrato final;
- `go test ./...`;
- `git diff --check`.

## Recuperação

Falha na descoberta não deve modificar arquivos. Falha na persistência deve
preservar o `darp.yml` original e informar a causa. A estratégia final deverá
reutilizar as garantias de atualização segura decididas na Spec 004.

## Premissas

- A Spec 004 será concluída e aprovada antes desta feature.
- Assets externos não serão executados nem convertidos automaticamente.
- A descoberta será agnóstica a linguagem, framework e provedor.
- Assets ambiguos ou conflitantes devem ser informados, nao resolvidos
   automaticamente. A identidade semantica e `(categoria, nome normalizado)`;
   categorias usam o tipo do asset, skills e extensions usam o diretorio pai,
   os demais tipos usam o nome do arquivo sem seu sufixo convencional.

## Estado de validação

Validação concluída. O `doctor` valida todas as combinações
`(path, family, type)`, rejeita symlinks registrados e combinações
incompatíveis, classifica diretório no lugar de arquivo como `FAIL` e trata
asset ausente como `WARNING`. A cobertura de conflitos semânticos e de falhas
de persistência está completa, incluindo ausência de efeitos colaterais e
limpeza do temporário do `darp.yml`.
