# CLAUDE.md — Phresh-Go

## Identidade do projeto

**Phresh-Go** é a camada de domínio de uma plataforma de intermediação de serviços de limpeza.
- Linguagem: **Go 1.22**, módulo `phresh-go`
- Arquitetura: **DDD (Domain-Driven Design) + Clean Architecture**
- Dependências externas: **nenhuma** — apenas a biblioteca padrão do Go (`stdlib`)
- Idioma do código: **português** (nomes de tipos, funções, variáveis, erros)

---

## Estrutura do projeto

```
phresh-go/
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
- Métodos de verificação: `EPessoal()`, `EFaxineiro()`, `EPublicadoPor()`, `EstaPreenchido()`

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
| Publicar serviço | Faxineiro cria `Limpeza` com `ValorHora` e `DuracaoEstimada` |
| Preço total | `ValorHora × DuracaoEstimada` — capturado na criação da `Solicitacao` |
| Solicitar serviço | Cliente cria `Solicitacao` → sistema verifica disponibilidade → status `pendente` |
| Verificação de agenda | 2 passos: (1) cobre disponibilidade semanal? (2) sem conflito de bloqueio? |
| Aceitar | Faxineiro aceita → verifica novamente → cria `BloqueioServico` na agenda |
| Cancelar aceita | Multa de **20%** se cancelar com **< 24h** de antecedência |
| Avaliar | Apenas solicitações `aceita` → cria `Avaliacao` → marca como `concluida` |
| Bloqueio pessoal | Faxineiro pode bloquear horários pessoais (sem vincular a serviço) |

### Máquina de estados da Solicitação
```
pendente → aceita (faxineiro) → concluída (via avaliação)
pendente → rejeitada (faxineiro)
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
