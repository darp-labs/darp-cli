# Contribuindo com o DARP CLI

Obrigado por contribuir com o DARP CLI. O projeto busca tornar assets reutilizáveis de IA versionáveis, compartilháveis e instaláveis como pacotes.

## Antes de começar

Leia estes documentos antes de propor uma mudança:

- [README](README.md)
- [AGENTS.md](AGENTS.md)
- [Código de Conduta](CODE_OF_CONDUCT.md)

Especificações, planos, ADRs e contexto interno são mantidos no repositório privado `darp-cli-specs`. Para uma contribuição funcional, peça a uma pessoa mantenedora a confirmação do escopo aprovado; não inclua cópias de documentos privados neste repositório público.

Para uma mudança relevante, procure uma issue ou abra uma nova discussão antes de implementar. O projeto segue Specification-Driven Development; funcionalidades não devem ser implementadas sem uma especificação aprovada.

## Tipos de contribuição

São bem-vindas contribuições de:

- código Go e correções de bugs;
- testes automatizados;
- documentação e exemplos;
- skills, prompts e outros assets de IA;
- melhorias de governança e ferramentas de desenvolvimento;
- relatos de bugs e propostas de funcionalidades.

## Configuração local

Requisitos:

- Git;
- Go compatível com a versão declarada em [go.mod](go.mod).

Clone o repositório e execute as validações disponíveis:

```bash
git clone https://github.com/darpbr/darp-cli.git
cd darp-cli
make verify
```

Comandos úteis:

| Comando | Objetivo |
| --- | --- |
| `make build` | Compila o CLI para a plataforma atual. |
| `make test` | Executa os testes Go. |
| `make lint` | Executa `golangci-lint` ou, como fallback, `go vet`. |
| `make fmt` | Formata os arquivos Go. |
| `make verify` | Executa formatação, testes e lint. |
| `make run` | Executa o CLI localmente. |

## Fluxo de desenvolvimento

1. Crie uma branch descritiva a partir de `main`.
2. Para mudanças funcionais, confirme com uma pessoa mantenedora a especificação aprovada no repositório privado `darp-cli-specs`; contribuições sem acesso devem seguir um escopo público aprovado.
3. Mantenha a alteração pequena, focada e compatível com a arquitetura existente.
4. Adicione ou atualize testes quando o comportamento mudar.
5. Atualize a documentação afetada.
6. Execute `make verify` antes de abrir o pull request.
7. Abra um pull request usando o template do repositório.

O fluxo do projeto é:

`Vision → Constitution → Specification → Plan → Tasks → Implementation → Review`

## Padrões de implementação

- Prefira Go idiomático e a biblioteca padrão quando ela for suficiente.
- Mantenha funções focadas, pacotes coesos e baixo acoplamento.
- Preserve comportamento existente e evite dependências desnecessárias.
- Não adicione segredos, credenciais, tokens ou dados pessoais ao repositório.
- Não documente funcionalidades que ainda não existem.
- Explique no pull request decisões que afetem arquitetura, compatibilidade ou segurança.

## Contribuições assistidas por IA

Ferramentas de IA podem ser usadas, mas a pessoa autora permanece responsável por todo o conteúdo enviado. Revise criticamente o resultado, confirme as fontes, execute os testes e certifique-se de compreender as alterações e seus impactos.

## Issues e vulnerabilidades

Use issues para bugs reproduzíveis, dúvidas de projeto e propostas de melhoria. Pesquise issues existentes antes de abrir uma nova e inclua contexto, passos para reprodução e resultado esperado quando aplicável.

Não publique vulnerabilidades, credenciais ou informações sensíveis em issues ou pull requests públicos. Para uma questão de segurança ou de privacidade, contate os mantenedores pelo perfil [@darpbr](https://github.com/darpbr), solicitando um canal privado.

## Pull requests

Pull requests devem:

- descrever o problema e a solução;
- relacionar a issue ou especificação correspondente;
- informar testes e comandos executados;
- registrar riscos, limitações e possíveis breaking changes;
- conter somente alterações relacionadas ao objetivo;
- estar prontos para revisão ou marcados como Draft quando ainda estiverem em andamento.

Os mantenedores podem solicitar ajustes de código, testes, documentação, arquitetura ou escopo. Responda aos comentários no próprio pull request para manter o histórico público e rastreável.

## Revisão e merge

Um pull request pode ser integrado quando o escopo estiver claro, os checks aplicáveis passarem, as revisões forem tratadas e a documentação estiver alinhada. A aprovação final é responsabilidade dos mantenedores definidos em [CODEOWNERS](.github/CODEOWNERS).

## Código de Conduta

A participação no projeto está sujeita ao [Código de Conduta](CODE_OF_CONDUCT.md). Questões de comportamento devem ser tratadas de forma respeitosa e, quando necessário, reportadas em privado aos mantenedores.
