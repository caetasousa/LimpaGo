# CLAUDE.md — LimpaGo

## Identidade do projeto

**LimpaGo** é uma plataforma de intermediação de serviços de limpeza.
- Linguagem: **Go 1.22**, módulo `limpaGo`
- Arquitetura: **DDD (Domain-Driven Design) + Clean Architecture**
- Camada HTTP: **Chi** com **Swagger** (swaggo)
- Idioma do código: **português** (nomes de tipos, funções, variáveis, erros)

---

## Estrutura do projeto

```
limpaGo/
├── domain/
│   ├── entity/          Entidades com identidade (Usuario, Limpeza, Solicitacao, Agenda...)
│   ├── valueobject/     Objetos de valor imutáveis (TipoLimpeza, Nota, Endereco, Paginacao...)
│   ├── service/         Serviços de domínio — regras de negócio entre múltiplas entidades
│   ├── repository/      Interfaces de repositório — contratos de persistência (sem implementação)
│   ├── errors/          Erros sentinela do domínio (erros.go)
│   └── testutil/        Mocks in-memory dos repositórios para uso em testes
├── DOCUMENTACAO.md      Documentação completa do domínio e regras de negócio
└── CLAUDE.md            Este arquivo
```

---

## Comandos essenciais

```bash
# Rodar todos os testes com race detector
go test ./... -race -count=1

# Cobertura de testes
go test ./... -cover

# Análise estática
go vet ./...

# Verificar compilação
go build ./...
```

---

## Convenções de código

### Nomes e construtores
- Tudo em **português** — não traduzir para inglês
- Construtores: `NovoXxx()` ou `NovaXxx()` retornam `(*Tipo, error)` ou `*Tipo`
- Métodos de verificação: `EPessoal()`, `EProfissional()`, `EPublicadoPor()`, `EstaPreenchido()`

### Erros
- Erros sentinela em `domain/errors/erros.go`: `var ErrXxx = errors.New("...")`
- Erros de validação de campo: `&entity.ErroValidacao{Campo: "nome_campo", Mensagem: "..."}`
- Propagar com `fmt.Errorf("contexto: %w", err)`
- **Nunca usar `panic` para tratamento de erro**

### Value objects
- São imutáveis; definidos pelo conteúdo, não por ID
- Sempre têm `Validar() error` quando carregam restrições de domínio
- Exemplo: `TipoLimpeza`, `Nota`, `Endereco`, `Paginacao`

### Repositórios
- São **interfaces** em `domain/repository/` — nunca implementar persistência real no domínio
- Injetados nos serviços via construtor: `NovoServicoXxx(repo RepositorioXxx) *ServicoXxx`
- Mocks para testes ficam em `domain/testutil/`

---

## Convenções de teste

Seguir as skills instaladas em `.agents/skills/golang-testing/SKILL.md` e `.agents/skills/golang-pro/SKILL.md`.

### Padrão obrigatório
- **Table-driven** com `t.Run()` subtests
- `t.Parallel()` em todos os testes e subtests independentes
- `t.Helper()` em todas as funções auxiliares de teste
- Mensagens de erro no formato: `got X; want Y`
- Campo `wantErr bool` para consolidar happy/sad path em uma única tabela
- Campo `wantErrIs error` quando precisa verificar sentinela específica com `errors.Is()`

### Organização
- Serviços testados com **black-box** (`package service_test`)
- Mocks in-memory de `domain/testutil/` para todos os repositórios
- Rodar sempre com `-race -count=1`

### Exemplo de estrutura
```go
func TestNovaXxx(t *testing.T) {
    t.Parallel()
    tests := []struct {
        name       string
        input      string
        want       string
        wantErr    bool
        wantErrIs  error
    }{
        {name: "válido", input: "ok", want: "ok"},
        {name: "vazio", input: "", wantErr: true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()
            got, err := NovaXxx(tt.input)
            if tt.wantErr {
                if err == nil {
                    t.Fatal("expected error; got nil")
                }
                if tt.wantErrIs != nil && !errors.Is(err, tt.wantErrIs) {
                    t.Errorf("got err %v; want %v", err, tt.wantErrIs)
                }
                return
            }
            if err != nil {
                t.Fatalf("unexpected error: %v", err)
            }
            if got.Campo != tt.want {
                t.Errorf("got %v; want %v", got.Campo, tt.want)
            }
        })
    }
}
```

---

## Regras de negócio principais

| Fluxo | Descrição |
|---|---|
| Publicar serviço | Profissional cria `Limpeza` com `ValorHora` e `DuracaoEstimada` |
| Preço total | `ValorHora × DuracaoEstimada` — capturado na criação da `Solicitacao` |
| Solicitar serviço | Cliente cria `Solicitacao` → sistema verifica disponibilidade → status `pendente` |
| Verificação de agenda | 2 passos: (1) cobre disponibilidade semanal? (2) sem conflito de bloqueio? |
| Aceitar | Profissional aceita → verifica novamente → cria `BloqueioServico` na agenda |
| Cancelar aceita | Multa de **20%** se cancelar com **< 24h** de antecedência |
| Avaliar | Apenas solicitações `aceita` → cria `Avaliacao` → marca como `concluida` |
| Bloqueio pessoal | Profissional pode bloquear horários pessoais (sem vincular a serviço) |

### Máquina de estados da Solicitação
```
pendente → aceita (profissional) → concluída (via avaliação)
pendente → rejeitada (profissional)
pendente → cancelada (cliente, sem multa)
aceita   → cancelada (cliente, possível multa 20%)
```

---

## O que NÃO fazer

- ❌ Não adicionar dependências externas ao `go.mod`
- ❌ Não implementar HTTP, banco de dados ou infraestrutura na camada `domain/`
- ❌ Não criar arquivos `README.md` ou documentação extra sem ser pedido
- ❌ Não renomear ou traduzir nomes existentes (o projeto é em pt-BR por design)
- ❌ Não usar `panic` para tratamento de erros de domínio
- ❌ Não acessar repositórios diretamente nos serviços que já delegam para `ServicoAgenda`

---

## Skills instaladas

- `.agents/skills/golang-testing/SKILL.md` — padrões de teste Go (table-driven, benchmarks, fuzzing)
- `.agents/skills/golang-pro/SKILL.md` — boas práticas Go (concorrência, interfaces, generics, estrutura de projeto)

---

## Padrão de commits

Cada commit deve descrever **com precisão** o que foi alterado. Nunca agrupar mudanças não relacionadas em um único commit vago.

### Formato obrigatório

```
<tipo>: <resumo claro do que foi feito>

- <detalhe 1 do que mudou e por quê>
- <detalhe 2 do que mudou e por quê>
- ...
```

**Tipos:** `feat`, `fix`, `refactor`, `docs`, `test`, `chore`

### Regras de granularidade

- Se foram alterados **arquivos de naturezas diferentes** (ex: `docker-compose.yml` + testes de integração + documentação), listar cada grupo no corpo do commit
- **Nunca** usar mensagens vagas como "ajustes", "correções", "updates" ou "misc"
- Citar os **arquivos principais** alterados quando a mudança não for óbvia pelo tipo
- Se a mudança envolve portas, URLs, variáveis de ambiente ou configurações — **mencionar os valores** no corpo

### Exemplos corretos

```
feat: adicionar testes de integração para repositórios PostgreSQL

- 9 arquivos de teste em infra/postgres/ com build tag integration
- testutil_test.go com helpers compartilhados (criarBancoTeste, limparTabelas)
- transacao_integracao_test.go valida commit e rollback
```

```
chore: ajustar portas para evitar conflito com outros projetos

- docker-compose.yml: postgres 5432 → 5434, api 8080 → mantida
- .env.exemplo: DATABASE_URL atualizada para porta 5434
- Makefile: run-pg atualizado para porta 5434
- DOCUMENTACAO.md: URLs corrigidas para localhost:8080
```

```
feat: adicionar Dockerfile, Makefile e .gitignore

- Dockerfile: multi-stage build golang:1.22 → distroless (~20MB)
- Makefile: comandos dev, test, test-integration, docker-up/down, swagger
- .gitignore: protege .env, binários e arquivos de IDE
```

---

## Fluxo de trabalho ao terminar alterações

### 1. Atualizar a documentação
Sempre que houver alterações no código, **atualizar obrigatoriamente**:

- **`DOCUMENTACAO.md`** — documento principal do projeto. Descreve o que é o projeto, a filosofia, os fluxos de negócio, as entidades, os serviços e as regras. Atualizar sempre que:
  - Novas entidades, serviços ou regras de negócio forem adicionadas
  - Endpoints da API forem criados, alterados ou removidos
  - Fluxos ou comportamentos existentes forem modificados
  - O nome ou propósito do projeto mudar
  - Novos comandos, variáveis de ambiente ou passos de setup forem adicionados

- **`README.md`** — é exatamente igual ao `DOCUMENTACAO.md`. Após atualizar a documentação, **sempre copiar o conteúdo inteiro** para o README:
  ```bash
  cp DOCUMENTACAO.md README.md
  ```

> **Regra**: `README.md` não tem conteúdo próprio — é sempre uma cópia fiel do `DOCUMENTACAO.md`.

#### Seção obrigatória no final de DOCUMENTACAO.md e README.md

O documento deve sempre terminar com uma seção **"Como rodar o projeto"** contendo os comandos completos e em ordem para subir o projeto do zero. Essa seção deve ser mantida atualizada sempre que o setup mudar. Exemplo de estrutura:

```markdown
## Como rodar o projeto

### Pré-requisitos
- Go 1.22+
- Docker e Docker Compose

### 1. Clonar e instalar dependências
...

### 2. Configurar variáveis de ambiente
...

### 3. Subir a infraestrutura
...

### 4. Rodar a API
...
```

### 2. Executar os testes relevantes

**Obrigatório antes de qualquer commit.** Executar os testes correspondentes à área alterada:

- **Tarefa apenas de frontend** (sem alteração em código Go/backend): executar apenas os testes do frontend (`npm test` dentro de `frontend/`). **NÃO** é necessário rodar `go test ./...` nem `make test-integration`.
- **Tarefa de backend** (ou que altera código Go): rodar a suíte completa do backend:
  ```bash
  go test ./... -race -count=1
  ```
  Se houver testes de integração relevantes à mudança, rodar também:
  ```bash
  make test-integration
  ```
- **Tarefa que altera frontend E backend**: rodar os testes de ambos.

> **Regra absoluta**: não fazer commit com testes falhando. Se um teste quebrou, corrigir antes de prosseguir.

### 3. Verificar consistência obrigatória antes de finalizar

**A documentação NUNCA pode estar em desacordo com o código.** Antes de concluir qualquer tarefa, verificar obrigatoriamente:

| O que conferir | Como verificar |
|----------------|----------------|
| Portas (HTTP, banco) | Comparar `main.go`, `docker-compose.yml`, `Makefile` e `DOCUMENTACAO.md` |
| `DATABASE_URL` e variáveis de ambiente | Comparar `.env.exemplo`, `docker-compose.yml` e `DOCUMENTACAO.md` |
| Comandos do Makefile | Executar `make -n <comando>` e confirmar que a doc descreve o mesmo comportamento |
| Endpoints da API | Conferir se os handlers em `api/handler/` batem com o que está documentado |
| Nomes de entidades, campos e serviços | Confirmar que a doc usa os mesmos nomes do código (pt-BR) |

> **Regra absoluta**: se qualquer valor no código (porta, URL, comando, nome) for diferente do que está na documentação, corrigir a documentação **antes** de encerrar a tarefa. Nunca deixar divergência para depois.

### 4. Pedir aprovação, fazer commit e perguntar sobre git push

> **⚠️ REGRA CRÍTICA — NUNCA PULAR ESTE PASSO ⚠️**
> Este passo é **OBRIGATÓRIO** ao final de QUALQUER tarefa que altere código, sem exceção.
> O Claude NÃO pode fazer commit sem antes ter a aprovação explícita do usuário.

**Sequência obrigatória:**

1. **Perguntar se o usuário aprova as alterações** antes de qualquer commit. Exibir um resumo claro do que foi feito e apresentar as opções:
   ```
   Resumo das alterações:
   - <o que foi alterado 1>
   - <o que foi alterado 2>

   O que deseja fazer?
   1. ✅ Aprovar — fazer commit e seguir para push
   2. ❌ Não aprovar — reverter e ajustar
   3. 🗑️ Descartar — apagar tudo que foi feito desde o último commit
   ```

2. **Se o usuário escolher "Não aprovar" (opção 2)**: reverter as alterações feitas na tarefa. Não fazer commit. Perguntar o que o usuário gostaria de diferente.

3. **Se o usuário escolher "Descartar" (opção 3)**: executar `git checkout -- .` e `git clean -fd` para descartar **todas** as alterações não commitadas, voltando ao estado exato do último commit. Confirmar ao usuário que tudo foi descartado. Encerrar a tarefa sem commit.

4. **Se o usuário aprovar (opção 1)**: fazer o commit seguindo o padrão de commits deste documento.

5. **Executar o comando** para listar todos os commits pendentes:
   ```bash
   git log origin/master..HEAD --pretty=format:"• %h %s (%cr)"
   ```

6. **Exibir a mensagem final** no chat, SEMPRE neste formato exato:
   ```
   Esses são os commits que serão enviados:

   • a1b2c3d feat: adicionar testes de integração (há 2 minutos)
   • d4e5f6g fix: corrigir porta do docker-compose (há 5 minutos)

   Deseja que eu faça `git push` para o GitHub agora?
   ```

> **Importante**: o commit só acontece **após aprovação** do usuário.
> **Importante**: a mensagem de commits + pergunta sobre push é a **última coisa** que o Claude escreve na resposta. Nada vem depois dela.

Se a resposta for sim, executar:
```bash
git push origin master
```

Remote configurado: `https://github.com/caetasousa/LimpaGo.git`
