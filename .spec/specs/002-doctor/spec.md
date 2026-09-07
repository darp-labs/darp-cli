# Spec-002 — `darp doctor`

## Status

Approved and implemented — validated on 2026-08-27.

---

# Objetivo

Implementar o comando:

```bash
darp doctor
```

O comando realiza um diagnóstico completo da integridade de um projeto DARP.

Seu objetivo é verificar se o projeto segue todos os contratos definidos pelo framework.

O comando é **somente leitura**.

Ele nunca deverá:

* modificar arquivos;
* criar arquivos;
* excluir arquivos;
* executar workflows;
* executar skills;
* corrigir problemas automaticamente.

---

# Escopo

O `doctor` valida exclusivamente artefatos pertencentes ao DARP.

Ele **não** valida:

* código da aplicação;
* dependências externas;
* linguagens de programação;
* compiladores;
* ferramentas como Docker, Java, Node, Go etc.

---

# Estrutura esperada do projeto

```text
.
├── README.md
├── darp.yml
├── specs/
└── .darp/
    ├── lifecycle.md
    ├── governance/
    │   └── quality-gates.md
    ├── workflows/
    └── templates/
└── .agents/
    └── skills/
```

---

# Contrato do arquivo darp.yml

O arquivo é obrigatório.

Deve possuir sintaxe YAML válida.

## Schema mínimo

```yaml
version: "1.0"

project:
  name: darp-cli

governance:
  lifecycle: .darp/lifecycle.md

workflows:
  default: implement

skills:
  documentation: .agents/skills/documentation
```

## Campos obrigatórios

| Campo                | Obrigatório |
| -------------------- | ----------- |
| version              | Sim         |
| project              | Sim         |
| project.name         | Sim         |
| governance           | Sim         |
| governance.lifecycle | Sim         |
| workflows            | Sim         |
| skills               | Sim         |

## Regras

* campos obrigatórios não podem ser vazios;
* caminhos devem existir;
* caminhos devem ser relativos ao projeto;
* chaves desconhecidas são permitidas (forward compatibility);
* ausência de qualquer campo obrigatório gera FAIL.

---

# Contrato dos Workflows

Todos os workflows devem estar em:

```text
.darp/workflows/
```

Formato:

```text
<nome>.yaml
```

Cada workflow deve conter no mínimo:

```yaml
name: implement

steps:
  - implement
  - test
  - documentation
```

## Regras

Todo workflow deve possuir:

* name
* steps

As regras são:

* nome único;
* steps não vazio;
* ordem preservada;
* cada step deve possuir uma skill correspondente.

Exemplo inválido:

```yaml
name: implement

steps: []
```

---

# Contrato das Skills

Toda skill de agente reside em:

```text
.agents/skills/<nome>/
```

Estrutura mínima obrigatória:

```text
documentation/

    SKILL.md
```

Recursos auxiliares como `prompts/`, `examples/`, `references/`, `scripts/` e
`templates/` são opcionais. Quando presentes, devem ser diretórios.

## Regras

Uma skill válida deve possuir:

* diretório existente;
* SKILL.md;
* recursos auxiliares opcionais, quando utilizados, devem possuir o tipo correto.

Arquivos adicionais são permitidos.

---

# Contrato dos Templates

Todos os templates ficam em:

```text
.darp/templates/
```

Não há templates obrigatórios na versão 1.0. Arquivos adicionais são opcionais.

---

# Governança

Arquivos obrigatórios:

```text
.darp/lifecycle.md

.darp/governance/quality-gates.md
```

Ambos devem existir.

---

# Compatibilidade de versão

O campo

```yaml
version:
```

representa a versão do contrato do projeto DARP.

Não representa a versão da aplicação.

## Regras

O CLI possui uma versão de contrato suportada.

Exemplo:

CLI

```text
Supports:

1.x
```

Projeto

```yaml
version: "1.0"
```

Resultado:

PASS

Projeto

```yaml
version: "2.0"
```

Resultado:

FAIL

Projeto

```yaml
version: "1.5"
```

Resultado:

PASS

Projeto

```yaml
version: "0.9"
```

Resultado:

WARNING

A política é:

* mesma major → compatível;
* major superior → FAIL;
* major inferior → WARNING.

Essa regra permite evolução gradual do framework.

---

# Checks obrigatórios

O doctor executará obrigatoriamente os seguintes checks.

## Configuration

Valida:

* existência do darp.yml;
* schema;
* sintaxe;
* campos obrigatórios.

---

## Structure

Valida diretórios obrigatórios.

---

## Workflows

Valida:

* estrutura;
* arquivos;
* steps.

---

## Skills

Valida:

* estrutura mínima;
* diretórios;
* SKILL.md;
* prompts;
* examples.

---

## Templates

Valida templates obrigatórios.

---

## Governance

Valida:

* lifecycle.md;
* quality-gates.md.

---

## Version Compatibility

Valida compatibilidade entre contrato do projeto e contrato suportado pelo CLI.

---

# Estados possíveis

Todo check deve produzir exatamente um estado.

| Estado  | Significado                       |
| ------- | --------------------------------- |
| PASS    | Verificação concluída com sucesso |
| WARNING | Problema não crítico              |
| FAIL    | Problema crítico                  |

---

# Código de saída

| Código | Significado            |
| ------ | ---------------------- |
| 0      | Apenas PASS ou WARNING |
| 1      | Pelo menos um FAIL     |

Warnings nunca alteram o código de saída.

---

# Exemplos

## Projeto saudável

```text
Running DARP Doctor...

✔ Configuration
✔ Structure
✔ Workflows
✔ Skills
✔ Templates
✔ Governance
✔ Version Compatibility

Summary

Passed: 7
Warnings: 0
Errors: 0

Project healthy.
```

---

## Projeto com warning

```text
Running DARP Doctor...

✔ Configuration
✔ Structure
✔ Workflows
✔ Skills
✔ Templates
✔ Governance
⚠ Project contract version is older than CLI recommendation

Summary

Passed: 6
Warnings: 1
Errors: 0

Project healthy with warnings.
```

Exit Code

```text
0
```

---

## Projeto inválido

```text
Running DARP Doctor...

✔ Configuration
✖ Missing .darp/workflows
✔ Skills
✖ Workflow "implement" not found
✔ Templates
✔ Governance
✖ Unsupported contract version 2.0

Summary

Passed: 4
Warnings: 0
Errors: 3

Project has errors.
```

Exit Code

```text
1
```

---

# Requisitos Não Funcionais

* Arquitetura baseada em checks.
* Cada check deve ser independente.
* Fácil inclusão de novos checks.
* Nenhum check deve depender diretamente de outro.
* Tempo de execução inferior a dois segundos para projetos típicos.

---

# Critérios de Aceitação

A implementação será considerada concluída quando:

* todos os contratos definidos nesta especificação forem validados;
* todos os checks obrigatórios forem implementados;
* warnings e failures forem tratados corretamente;
* o resumo final corresponder ao estado do projeto;
* o código de saída estiver de acordo com esta especificação;
* nenhum arquivo do projeto for modificado durante a execução.

---

# Evolução

Novos contratos poderão ser adicionados em versões futuras, como:

* plugins;
* policies;
* documentação;
* CHANGELOG;
* arquitetura;
* segurança;
* cobertura de testes.

A inclusão de novos contratos não deverá exigir alterações na arquitetura do `doctor`, apenas a implementação de novos checks.

# Alinhamento com o Projeto

A implementação desta especificação introduz o primeiro contrato formal do arquivo `darp.yml`.

Para evitar divergências entre o CLI, a documentação e os testes, este contrato deverá ser adotado como referência oficial do projeto DARP.

---

# Contrato Oficial do `darp.yml`

A partir desta especificação, todo projeto DARP deverá possuir um arquivo `darp.yml` compatível com o seguinte schema mínimo:

```yaml
version: "1.0"

project:
  name: "<project-name>"

governance:
  lifecycle: .darp/lifecycle.md

workflows:
  default: implement

skills:
  documentation: .agents/skills/documentation
```

Este passa a ser o contrato oficial utilizado pelo CLI.

Novos campos poderão ser adicionados em versões futuras sem quebrar a compatibilidade.

---

# Impactos no `darp init`

A implementação do comando:

```bash
darp init
```

deverá ser atualizada para gerar automaticamente um `darp.yml` compatível com este contrato.

O arquivo gerado pelo `init` passa a ser considerado a referência para novos projetos.

---

# Impactos nos Testes

Todos os testes automatizados que criam projetos DARP deverão utilizar este mesmo contrato.

Não deverão existir múltiplos formatos válidos de `darp.yml` dentro do repositório.

Sempre que um projeto de teste for criado, ele deverá representar exatamente a estrutura esperada pelo `darp doctor`.

---

# Impactos na Documentação

Toda documentação que descreva a estrutura de um projeto DARP deverá ser atualizada para refletir este contrato.

No mínimo, deverão ser revisados:

* `README.md`;
* documentação do comando `darp init`;
* exemplos de criação de projetos;
* exemplos presentes em `specs/`;
* demais documentos que descrevam a estrutura do projeto.

Após esta atualização, o schema definido nesta especificação passa a ser a única referência oficial para o formato do arquivo `darp.yml`.

---

# Critério de Aceitação Adicional

A implementação desta especificação somente será considerada concluída quando:

* o `darp init` gerar um `darp.yml` compatível com este contrato;
* todos os testes utilizarem o mesmo formato de `darp.yml`;
* não existir documentação descrevendo um formato diferente;
* o `darp doctor` validar exatamente este contrato, sem regras implícitas ou formatos alternativos.
