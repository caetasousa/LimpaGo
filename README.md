# LimpaGo — Plataforma de Serviços de Limpeza

## O que é este projeto?

O **LimpaGo** é uma plataforma de intermediação de serviços de limpeza, escrito em Go. Ele modela toda a lógica de negócio para conectar **faxineiros** (profissionais de limpeza) que publicam seus serviços com **clientes** que os contratam, expondo uma **API REST** documentada com Swagger.

O projeto segue os princípios de **Domain-Driven Design (DDD)** e **Arquitetura Limpa**, com separação clara entre domínio e infraestrutura. Toda a persistência é abstraída por interfaces de repositório. A camada HTTP usa o framework **Chi** e a documentação é gerada automaticamente pelo **swaggo**.

---

## Filosofia do modelo

Diferente de marketplaces baseados em lances — onde profissionais competem por preço, o que pode levar à desvalorização da mão de obra — este modelo **valoriza o profissional**:

1. O **faxineiro publica seus serviços** com o valor por hora que considera justo para cada tipo de limpeza
2. O **cliente navega o catálogo** de serviços disponíveis, compara preços e escolhe o que mais lhe interessa
3. O **cliente solicita** um serviço informando data, horário e endereço
4. O **faxineiro decide** se aceita ou rejeita cada solicitação, mantendo total controle sobre sua agenda
5. A **plataforma** cuida da intermediação: verificação de disponibilidade, controle de agenda, política de cancelamento e avaliações

Isso coloca o profissional no controle do seu trabalho, do seu preço e do seu tempo. O cliente se beneficia da transparência (preço visível antes de solicitar) e da qualidade (avaliações e verificação de documentos).

---

## Como funciona?

### Fluxo principal detalhado

#### 1. Registro do Usuário

O usuário se cadastra fornecendo **email** e **nome de usuário**. O sistema valida:
- O email não pode estar em uso por outro usuário
- O nome de usuário deve ter pelo menos 3 caracteres e conter apenas letras, números, underscores ou hífens
- O email é normalizado (lowercase, sem espaços)

Ao registrar, um **perfil base** é criado automaticamente com os dados compartilhados (nome completo, telefone, foto — inicialmente vazios). O usuário pode atualizá-los depois.

#### 2. Criação de Perfis Específicos

Um mesmo usuário pode ter **dois papéis simultâneos**: faxineiro e cliente. Para isso, ele cria perfis específicos:

- **Perfil Faxineiro** — para oferecer serviços. Inclui descrição profissional, anos de experiência, especialidades, cidades atendidas e documentação (RG, CPF, foto do documento). A plataforma pode marcar o profissional como verificado.
- **Perfil Cliente** — para contratar serviços. Inclui endereço completo do imóvel, tipo de imóvel (apartamento, casa ou comercial), número de quartos e banheiros, tamanho em m², observações (ex: "tem animais de estimação") e opcionalmente um faxineiro preferido.

Cada tipo de perfil só pode ser criado uma vez por usuário. Se já existe, o sistema retorna erro.

#### 3. Publicação de Serviço (Limpeza)

O faxineiro publica um serviço de limpeza informando:

| Campo | Obrigatório | Descrição |
|---|---|---|
| Nome | Sim | Nome do serviço (ex: "Limpeza Residencial Completa") |
| Descrição | Não | Detalhes sobre o que inclui |
| Valor por hora | Sim | Quanto cobra por hora (deve ser > 0) |
| Duração estimada | Sim | Quantas horas o serviço normalmente leva (deve ser > 0) |
| Tipo de limpeza | Sim | Um dos 7 tipos aceitos pelo sistema |

O **preço total** é calculado automaticamente: `valor por hora × duração estimada`. O faxineiro pode atualizar qualquer campo depois, e também pode deletar um serviço que publicou. Apenas o faxineiro que publicou pode modificar ou deletar o serviço — essa verificação é feita em todas as operações.

#### 4. Configuração da Agenda

O faxineiro configura sua **disponibilidade semanal** definindo blocos de horário para cada dia da semana. Cada bloco tem:

- **Dia da semana** — domingo (0) a sábado (6)
- **Hora de início** — entre 0 e 23
- **Hora de fim** — entre 1 e 24 (deve ser maior que hora de início)

Exemplo: o faxineiro pode definir que trabalha de segunda a sexta das 8h às 12h e das 14h às 18h, e sábado das 8h às 12h. Cada bloco é independente e pode ser adicionado ou removido individualmente.

Além da disponibilidade, o faxineiro pode criar **bloqueios pessoais** para horários em que não poderá trabalhar (consultas médicas, compromissos pessoais, folgas). Esses bloqueios impedem que clientes solicitem serviços nesses horários, da mesma forma que um serviço já agendado impediria.

#### 5. Catálogo de Serviços

Clientes navegam o catálogo de todos os serviços publicados. O catálogo é **paginado** (padrão: 20 itens por página, máximo: 100) para lidar com grandes volumes de serviços.

#### 6. Solicitação de Serviço

O cliente escolhe um serviço e cria uma solicitação informando a **data e horário desejados**. O sistema faz as seguintes verificações antes de criar:

1. **Faxineiro não pode solicitar o próprio serviço** — se o clienteID é o mesmo que publicou a limpeza, a solicitação é rejeitada
2. **Data não pode ser no passado** — o horário solicitado deve ser futuro
3. **Sem duplicatas** — o cliente não pode ter outra solicitação aberta para o mesmo serviço
4. **Verificação de disponibilidade** — o sistema calcula o período completo (data de início + duração estimada) e verifica:
   - O faxineiro tem um bloco de disponibilidade que cobre esse período nesse dia da semana?
   - Não há nenhum bloqueio (serviço ou pessoal) que conflita com esse período?

O **preço total** é fixado no momento da solicitação (captura o valor vigente) e a solicitação nasce no estado **pendente**.

O cliente também pode definir o **endereço** onde o serviço será realizado — pode digitar manualmente ou copiar do endereço salvo no seu perfil de cliente.

#### 7. Aceitação ou Rejeição pelo Faxineiro

O faxineiro visualiza as solicitações pendentes para seus serviços e pode:

- **Aceitar** — o sistema faz uma **segunda verificação de disponibilidade** (o horário pode ter sido ocupado desde a criação da solicitação). Se ainda estiver disponível, a solicitação muda para **aceita** e um **bloqueio de serviço** é criado automaticamente na agenda, impedindo conflitos futuros.
- **Rejeitar** — a solicitação muda para **rejeitada**. Apenas solicitações pendentes podem ser rejeitadas. Nenhum bloqueio é criado.

Apenas o faxineiro que publicou o serviço pode aceitar ou rejeitar — essa propriedade é verificada antes de qualquer ação.

#### 8. Cancelamento

O cliente pode cancelar solicitações nos estados **pendente** ou **aceita**. A política de multa funciona assim:

| Situação | Consequência |
|---|---|
| Cancelar solicitação **pendente** | Sem custo, sem impacto na agenda |
| Cancelar solicitação **aceita** com **24h+ de antecedência** | Sem custo, bloqueio liberado na agenda |
| Cancelar solicitação **aceita** com **menos de 24h** antes do serviço | **Multa de 20%** do preço total, bloqueio liberado na agenda |

A multa é calculada automaticamente: `preço total × 0.20`. Se a solicitação estava aceita, o bloqueio associado na agenda do faxineiro é removido, liberando o horário.

Apenas o cliente que criou a solicitação pode cancelá-la. Solicitações já rejeitadas, canceladas ou concluídas não podem ser canceladas novamente.

#### 9. Avaliação e Conclusão

Após o serviço ser realizado, o cliente avalia o faxineiro com:

- **Nota** — valor inteiro entre 0 e 5
- **Comentário** — texto opcional descrevendo a experiência

Regras:
- Apenas solicitações no estado **aceita** podem ser avaliadas
- Cada solicitação pode ser avaliada **uma única vez**
- Ao criar a avaliação, a solicitação é automaticamente marcada como **concluída**

O sistema mantém um **agregado de avaliação** por faxineiro, com a média das notas e o total de avaliações recebidas. Isso serve como reputação pública do profissional.

#### 10. Feed de Atividades

Um feed paginado mostra os eventos recentes da plataforma: serviços publicados e atualizados por faxineiros. Cada item do feed contém:

- O serviço (Limpeza) associado
- O tipo de evento (criação ou atualização)
- A data do evento
- Um número de linha para paginação por cursor

O feed utiliza paginação com tamanho configurável (padrão: 20, máximo: 100) e indica se há mais páginas.

---

### Máquina de estados da Solicitação

A solicitação segue um ciclo de vida bem definido com transições controladas:

```
                    ┌─────────────────────────────────────────────┐
                    │                                             │
                    ▼                                             │
              ┌──────────┐   aceitar    ┌─────────┐   avaliar   ┌──────────┐
  criar ────▶ │ pendente │ ──────────▶  │ aceita  │ ─────────▶  │concluída │
              └──────────┘              └─────────┘             └──────────┘
                    │                        │
                    │ rejeitar               │ cancelar
                    ▼                        ▼
              ┌──────────┐              ┌──────────┐
              │rejeitada │              │cancelada │
              └──────────┘              └──────────┘
                    ▲                        ▲
                    │                        │
                    │      cancelar          │
                    └── (do pendente) ───────┘
```

| Transição | Quem executa | Efeitos colaterais |
|---|---|---|
| pendente → aceita | Faxineiro | Verifica disponibilidade novamente, cria bloqueio na agenda |
| pendente → rejeitada | Faxineiro | Nenhum |
| pendente → cancelada | Cliente | Nenhum |
| aceita → concluída | Sistema (via avaliação) | Nenhum |
| aceita → cancelada | Cliente | Libera bloqueio na agenda, possível multa de 20% |

---

### Sistema de Agenda

A agenda do faxineiro é composta por dois mecanismos complementares:

#### Disponibilidade (recorrência semanal)

Define **quando** o faxineiro pode trabalhar. São blocos de horário vinculados a dias da semana, que se repetem toda semana:

- Cada bloco tem dia da semana, hora início e hora fim
- Exemplo: "Segunda das 8h às 12h", "Segunda das 14h às 18h", "Terça das 8h às 18h"
- O faxineiro pode ter múltiplos blocos por dia
- Blocos podem ser adicionados e removidos livremente

Quando um cliente solicita um serviço, o sistema verifica se o período inteiro (início + duração) cai dentro de algum bloco de disponibilidade daquele dia da semana.

#### Bloqueios (datas específicas)

Representam horários **ocupados** em datas específicas. Existem dois tipos:

| Tipo | Criação | Remoção | SolicitacaoID |
|---|---|---|---|
| **Bloqueio de serviço** | Automática — ao aceitar solicitação | Automática — ao cancelar solicitação | Preenchido |
| **Bloqueio pessoal** | Manual — pelo faxineiro | Manual — pelo faxineiro | Nulo (nil) |

Ambos os tipos impedem que novos serviços sejam agendados no período bloqueado. A diferença é que bloqueios de serviço são gerenciados automaticamente pelo ciclo de vida da solicitação, enquanto bloqueios pessoais são controlados diretamente pelo faxineiro.

O bloqueio pessoal não exige motivo — o faxineiro simplesmente marca que não estará disponível naquele período.

Validações comuns a ambos os tipos:
- A data de fim deve ser posterior à data de início
- Não é possível criar bloqueios no passado

#### Fluxo de verificação completo

Quando o sistema precisa verificar se um horário está livre (na criação e na aceitação da solicitação):

1. Calcula o período completo: `data solicitada` até `data solicitada + duração estimada`
2. Busca os blocos de disponibilidade do faxineiro para aquele dia da semana
3. Verifica se existe pelo menos um bloco que **contém** o período inteiro (hora início ≤ hora solicitada E hora fim ≥ hora término)
4. Busca todos os bloqueios (serviço e pessoal) que se sobrepõem ao período
5. Se encontrar qualquer bloqueio, rejeita por conflito de agenda

---

## Tipos de serviço

A plataforma suporta 7 tipos de limpeza, cada um adequado a uma necessidade diferente:

| Tipo | Constante | Residencial? | Descrição |
|---|---|---|---|
| **Limpeza Padrão** | `limpeza_padrao` | Sim | Limpeza de rotina para manter a qualidade e higiene do ambiente |
| **Limpeza Pesada** | `limpeza_pesada` | Sim | Limpeza profunda com maior atenção aos detalhes, ideal para limpezas periódicas |
| **Limpeza Express** | `limpeza_express` | Sim | Serviço rápido com tarefas padronizadas: louça, cama, pano, sanitários, lixo |
| **Limpeza Pré-Mudança** | `limpeza_pre_mudanca` | Sim | Preparar o imóvel antes de uma mudança, deixando-o pronto para o novo morador |
| **Limpeza Pós-Obra** | `limpeza_pos_obra` | Sim | Para ambientes que passaram por reformas, removendo poeira e resíduos de construção |
| **Limpeza Comercial** | `limpeza_comercial` | Não | Para escritórios, consultórios, lojas e outros ambientes comerciais |
| **Passadoria** | `passadoria` | Sim | Serviço de passar roupas |

O tipo é validado na criação e atualização do serviço. O método `EResidencial()` permite filtrar ou categorizar serviços para residências vs. ambientes comerciais.

---

## Precificação

### Modelo de preço

O faxineiro define para **cada serviço** que publica:

| Campo | Tipo | Exemplo | Descrição |
|---|---|---|---|
| **Valor por hora** | `float64` | R$ 50,00/h | Quanto cobra por hora para este tipo de serviço |
| **Duração estimada** | `float64` (horas) | 3.0 | Quanto tempo o serviço normalmente leva |
| **Preço total** | calculado | R$ 150,00 | `ValorHora × DuracaoEstimada` |

O preço total é calculado automaticamente e visível para o cliente antes de solicitar. Quando o cliente cria a solicitação, o preço total é **capturado** naquele momento — se o faxineiro alterar o preço depois, solicitações já criadas mantêm o valor original.

### Política de cancelamento

| Situação | Multa | Cálculo |
|---|---|---|
| Solicitação **pendente** | Nenhuma | — |
| Solicitação **aceita**, cancelada com **24h+ de antecedência** | Nenhuma | — |
| Solicitação **aceita**, cancelada com **menos de 24h** | **20%** do preço total | `PrecoTotal × 0.20` |

A multa protege o profissional contra cancelamentos de última hora, onde ele provavelmente já reservou aquele horário e recusou outros clientes.

---

## Perfis e papéis dos usuários

Um mesmo usuário pode atuar como **faxineiro** e/ou como **cliente**. O sistema possui três níveis de perfil:

### Perfil Base (`Perfil`)
Criado **automaticamente** no registro. Contém dados pessoais compartilhados entre os dois papéis:

| Campo | Tipo | Descrição |
|---|---|---|
| NomeCompleto | `string` | Nome real do usuário |
| Telefone | `string` | Número para contato |
| Imagem | `string` | URL da foto de perfil |
| Email | `string` | Copiado do cadastro (desnormalizado) |
| NomeUsuario | `string` | Copiado do cadastro (desnormalizado) |

### Perfil Faxineiro (`PerfilFaxineiro`)
Criado manualmente quando o usuário quer **oferecer serviços**:

| Campo | Tipo | Descrição |
|---|---|---|
| Descricao | `string` | Apresentação profissional / bio de trabalho |
| AnosExperiencia | `int` | Tempo de experiência na área |
| Especialidades | `[]string` | Lista de tipos de serviço que domina (ex: `["limpeza_padrao", "limpeza_pesada"]`) |
| CidadesAtendidas | `[]string` | Cidades onde aceita trabalhar (ex: `["São Paulo", "Guarulhos"]`) |
| DocumentoRG | `string` | Número do RG |
| DocumentoCPF | `string` | Número do CPF |
| FotoDocumento | `string` | URL da foto do documento para verificação |
| Verificado | `bool` | Se passou pelo processo de verificação da plataforma |

### Perfil Cliente (`PerfilCliente`)
Criado manualmente quando o usuário quer **contratar serviços**:

| Campo | Tipo | Descrição |
|---|---|---|
| Endereco | `Endereco` (value object) | Endereço completo do imóvel (rua, complemento, bairro, cidade, estado, CEP) |
| TipoImovel | `TipoImovel` | `apartamento`, `casa` ou `comercial` |
| Quartos | `int` | Número de quartos (ajuda a estimar duração) |
| Banheiros | `int` | Número de banheiros (ajuda a estimar duração) |
| TamanhoImovelM2 | `float64` | Tamanho em metros quadrados |
| Observacoes | `string` | Ex: "tem animais de estimação", "portaria 24h" |
| FaxineiroPreferidoID | `*int` | ID do faxineiro preferido (opcional) |

### Ações por papel

| Ação | Perfil necessário | Serviço responsável |
|---|---|---|
| Registrar e criar perfil base | Nenhum (qualquer usuário) | `ServicoUsuario.Registrar` |
| Atualizar dados pessoais | Perfil base | `ServicoUsuario.AtualizarPerfil` |
| Publicar serviço de limpeza | Faxineiro | `ServicoLimpeza.Criar` |
| Configurar agenda de disponibilidade | Faxineiro | `ServicoAgenda.AdicionarDisponibilidade` |
| Bloquear horário pessoal | Faxineiro | `ServicoAgenda.CriarBloqueioPessoal` |
| Aceitar/Rejeitar solicitação | Faxineiro | `ServicoSolicitacao.Aceitar/Rejeitar` |
| Navegar catálogo | Cliente | `ServicoLimpeza.ListarCatalogo` |
| Solicitar serviço | Cliente | `ServicoSolicitacao.CriarSolicitacao` |
| Cancelar solicitação | Cliente | `ServicoSolicitacao.CancelarSolicitacao` |
| Avaliar após o serviço | Cliente | `ServicoAvaliacao.CriarAvaliacao` |

---

## Detalhamento do domínio

```
limpaGo/
└── domain/                               Camada de domínio
    │
    ├── entity/                           Entidades — objetos com identidade própria
    │   ├── usuario.go                      Usuario (ID, Email, NomeUsuario, EFaxineiro(), ECliente())
    │   ├── perfil.go                       Perfil + PerfilFaxineiro + PerfilCliente
    │   ├── limpeza.go                      Limpeza (ValorHora, DuracaoEstimada, PrecoTotal())
    │   ├── agenda.go                       Disponibilidade + Bloqueio (serviço e pessoal)
    │   ├── solicitacao.go                  Solicitacao (ciclo de vida, endereço, multa)
    │   ├── avaliacao.go                    Avaliacao (nota + comentário) + AgregadoAvaliacao
    │   ├── feed.go                         ItemFeed + PaginaFeed
    │   └── erro_validacao.go               ErroValidacao (campo + mensagem)
    │
    ├── valueobject/                      Objetos de valor — imutáveis, sem identidade
    │   ├── tipo_limpeza.go                 TipoLimpeza (7 tipos + validação + EResidencial())
    │   ├── tipo_imovel.go                  TipoImovel (apartamento, casa, comercial)
    │   ├── status_solicitacao.go           StatusSolicitacao (5 estados + regras de transição)
    │   ├── nota.go                         Nota (int, 0-5)
    │   ├── endereco.go                     Endereco (rua, complemento, bairro, cidade, estado, CEP)
    │   ├── paginacao.go                    Paginacao (pagina, tamanho com validação)
    │   └── tipo_evento_feed.go             TipoEventoFeed (criacao, atualizacao)
    │
    ├── repository/                       Interfaces de repositório — contratos de persistência
    │   ├── repositorio_usuario.go          BuscarPorEmail, BuscarPorNomeUsuario, Salvar
    │   ├── repositorio_perfil.go           CRUD para Perfil + PerfilFaxineiro + PerfilCliente
    │   ├── repositorio_limpeza.go          CRUD + ListarPorFaxineiro + ListarTodas (catálogo)
    │   ├── repositorio_agenda.go           Disponibilidade (listar, salvar, deletar) + Bloqueios (listar, buscar, salvar, deletar)
    │   ├── repositorio_solicitacao.go      BuscarPorClienteELimpeza, ListarPorLimpeza, ListarPorCliente
    │   ├── repositorio_avaliacao.go        BuscarPorClienteELimpeza, ListarPorFaxineiro, AgregadoPorFaxineiro
    │   └── repositorio_feed.go             BuscarPaginaFeed
    │
    ├── service/                          Serviços de domínio — regras de negócio entre entidades
    │   ├── servico_usuario.go              Registro + CRUD dos 3 perfis
    │   ├── servico_limpeza.go              Publicar, atualizar, deletar, buscar, catálogo
    │   ├── servico_agenda.go               Disponibilidade + bloqueios (serviço e pessoal) + verificação
    │   ├── servico_solicitacao.go          Criar, aceitar, rejeitar, cancelar + integração com agenda
    │   ├── servico_avaliacao.go            Avaliar + estatísticas de reputação
    │   └── servico_feed.go                 Feed paginado de atividades
    │
    └── errors/                           Erros de domínio — sentinelas de negócio
        └── erros.go                        Todos os erros organizados por contexto
```

---

## Conceitos de arquitetura

### Entidades

Objetos com **identidade própria** (ID). Dois objetos com os mesmos dados mas IDs diferentes são entidades distintas. As entidades contêm suas próprias regras de validação e comportamento:

| Entidade | Responsabilidade | Regras principais |
|---|---|---|
| `Usuario` | Cadastro e autenticação | Validação de email e nome de usuário, referências opcionais aos perfis |
| `Perfil` | Dados pessoais compartilhados | Criado automaticamente no registro |
| `PerfilFaxineiro` | Dados profissionais | Documentação, verificação, especialidades |
| `PerfilCliente` | Dados do imóvel | Endereço, tipo, quartos, banheiros, tamanho |
| `Limpeza` | Serviço publicado pelo faxineiro | Validação de preço, duração, tipo; cálculo de preço total |
| `Solicitacao` | Pedido do cliente | Máquina de estados, cálculo de multa, endereço do serviço |
| `Avaliacao` | Nota do cliente ao faxineiro | Nota 0-5 + comentário, uma por solicitação |
| `Disponibilidade` | Bloco semanal de horário livre | Validação de hora início/fim |
| `Bloqueio` | Horário ocupado (serviço ou pessoal) | Validação de período, distinção por SolicitacaoID |

### Objetos de valor

Tipos **imutáveis** sem identidade, definidos pelo seu conteúdo. Dois objetos de valor com os mesmos dados são considerados iguais:

| Objeto de valor | Valores possíveis | Validação |
|---|---|---|
| `TipoLimpeza` | 7 constantes (`limpeza_padrao`, ..., `passadoria`) | `Validar()` rejeita valores fora do enum |
| `TipoImovel` | `apartamento`, `casa`, `comercial` | `Validar()` rejeita valores fora do enum |
| `StatusSolicitacao` | `pendente`, `aceita`, `rejeitada`, `cancelada`, `concluida` | Métodos de transição controlam quais mudanças são válidas |
| `Nota` | inteiro entre 0 e 5 | `NovaNota()` rejeita valores fora do intervalo |
| `Endereco` | rua, complemento, bairro, cidade, estado, CEP | `EstaPreenchido()` verifica campos mínimos |
| `Paginacao` | página (mín 1), tamanho (1-100, padrão 20) | `NovaPaginacao()` corrige valores inválidos automaticamente |
| `TipoEventoFeed` | `criacao`, `atualizacao` | `Validar()` rejeita valores fora do enum |

### Repositórios

**Interfaces** que definem o contrato de acesso a dados. A camada de domínio declara o que precisa, sem saber como os dados são armazenados. A implementação concreta (PostgreSQL, MongoDB, memória, etc.) fica na camada de infraestrutura.

Cada repositório é injetado nos serviços via construtor, seguindo o princípio de **inversão de dependência**.

### Serviços de domínio

Contêm regras de negócio que **envolvem múltiplas entidades** ou **coordenam operações complexas**. Cada serviço recebe seus repositórios necessários por injeção de dependência:

| Serviço | Dependências | Responsabilidade |
|---|---|---|
| `ServicoUsuario` | `RepositorioUsuario`, `RepositorioPerfil` | Registro + gerenciamento dos 3 tipos de perfil |
| `ServicoLimpeza` | `RepositorioLimpeza` | CRUD de serviços + catálogo paginado |
| `ServicoAgenda` | `RepositorioAgenda` | Disponibilidade + bloqueios + verificação de conflitos |
| `ServicoSolicitacao` | `RepositorioSolicitacao`, `RepositorioLimpeza`, `ServicoAgenda` | Ciclo de vida completo da solicitação (cria, aceita, rejeita, cancela) |
| `ServicoAvaliacao` | `RepositorioAvaliacao`, `RepositorioSolicitacao`, `RepositorioLimpeza` | Avaliação + reputação |
| `ServicoFeed` | `RepositorioFeed` | Feed de atividades paginado |

Note que `ServicoSolicitacao` depende de `ServicoAgenda` (não do repositório diretamente) — isso evita duplicação de lógica de verificação de disponibilidade.

### Erros de domínio

Erros sentinela (`var Err... = errors.New(...)`) que representam **violações de regras de negócio**. Organizados por contexto:

| Contexto | Exemplos de erros |
|---|---|
| Usuário | Email já utilizado, nome de usuário já utilizado |
| Perfil | Perfil já existe, perfil não encontrado |
| Limpeza | Não é o faxineiro da limpeza |
| Solicitação | Duplicada, não pode ser cancelada/rejeitada, faxineiro não pode solicitar próprio serviço |
| Agenda | Horário indisponível, conflito de agenda, agendamento no passado |
| Avaliação | Duplicada, solicitação não aceita |

Além dos erros sentinela, existe `ErroValidacao` — um tipo de erro estruturado com **campo** e **mensagem**, usado para validações de dados de entrada nas entidades.

---

## Regras de negócio consolidadas

### Precificação e valores
- O faxineiro define o valor por hora para cada serviço que publica
- O preço total é calculado automaticamente: `valor/hora × duração estimada`
- O preço é capturado na criação da solicitação (imutável depois)
- Multa de 20% por cancelamento tardio (< 24h antes do serviço aceito)

### Agenda e disponibilidade
- O faxineiro define blocos semanais de disponibilidade (dia da semana + hora início/fim)
- Bloqueios de serviço são criados/removidos automaticamente pelo ciclo de vida da solicitação
- Bloqueios pessoais são gerenciados diretamente pelo faxineiro, sem necessidade de motivo
- A disponibilidade é verificada **duas vezes**: na criação e na aceitação da solicitação
- Ambos os tipos de bloqueio impedem novas solicitações no mesmo período
- Não é possível criar bloqueios no passado

### Solicitação e ciclo de vida
- Faxineiro não pode solicitar o próprio serviço
- Cada cliente pode ter apenas uma solicitação **ativa** (pendente ou aceita) por serviço — pode solicitar novamente após conclusão ou cancelamento
- Apenas solicitações pendentes podem ser aceitas ou rejeitadas
- Apenas o faxineiro dono do serviço pode aceitar ou rejeitar
- Apenas o cliente que criou pode cancelar
- Cada solicitação tem seu próprio endereço via value object `Endereco` (pode copiar do perfil do cliente)
- As transições de estado são controladas e validadas
- Na aceitação, o bloqueio é criado antes de salvar a solicitação — se o bloqueio falhar, a solicitação permanece pendente (consistência transacional)

### Avaliação
- Apenas solicitações aceitas podem ser avaliadas
- Uma avaliação por solicitação (sem duplicatas)
- Nota inteira entre 0 e 5, com comentário opcional
- Criar avaliação marca a solicitação como concluída automaticamente
- Estatísticas agregadas por faxineiro (média + total)

### Perfis e permissões
- Perfil base é criado automaticamente no registro
- Cada tipo de perfil específico (faxineiro/cliente) só pode ser criado uma vez
- Um mesmo usuário pode ter ambos os perfis
- Ações são restritas por papel (faxineiro publica, cliente solicita)

---

## API REST

A API é servida em `/api/v1` com **31 endpoints** organizados por recurso:

| Grupo | Endpoints | Auth |
|---|---|---|
| Usuários e perfis | 9 | Maioria autenticada |
| Limpezas (catálogo) | 6 | Público para leitura |
| Solicitações | 6 | Autenticado |
| Agenda | 6 | Autenticado (faxineiro) |
| Avaliações | 3 | Público para leitura |
| Feed | 1 | Público |

A autenticação é feita via header `X-User-ID` (placeholder — preparado para JWT).

O **Swagger UI** fica disponível em `http://localhost:8080/swagger/index.html`.

---

## Estrutura do projeto

```
limpaGo/
├── go.mod                                Módulo Go (limpaGo, Go 1.22)
├── DOCUMENTACAO.md                       Esta documentação
│
├── api/                                  Camada HTTP
│   ├── dto/                              DTOs de request/response (JSON)
│   ├── handler/                          Handlers HTTP (1 por serviço de domínio)
│   ├── middleware/                       Auth, logger, CORS, recovery
│   ├── router/                           Registro de rotas com Chi
│   └── server/                           http.Server com timeouts
│
├── cmd/api/                              Entrypoint — composição e inicialização
│
├── docs/                                 Swagger gerado automaticamente (swag init)
│
└── domain/                               Camada de domínio
    ├── entity/                           Entidades com identidade
    ├── valueobject/                      Objetos de valor imutáveis
    ├── service/                          Serviços de domínio
    ├── repository/                       Interfaces de repositório
    ├── errors/                           Erros sentinela
    └── testutil/                         Mocks in-memory para testes
```

---

## Tecnologias

- **Linguagem:** Go 1.22
- **Arquitetura:** Domain-Driven Design (DDD) + Arquitetura Limpa
- **HTTP:** Chi router + go-chi/cors
- **Documentação:** swaggo/swag (OpenAPI 2.0 / Swagger UI)
- **Padrões utilizados:** Repository Pattern, Service Layer, Value Objects, Entity, Dependency Injection
