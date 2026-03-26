# Plano: Frontend LimpaGo — Etapas por Funcionalidade

## Contexto

O backend do LimpaGo está completo: API REST com Chi, autenticação JWT, 9 repositórios PostgreSQL testados, Swagger documentado, Docker Compose funcional. O próximo passo é construir o frontend, organizado por etapas de funcionalidade que seguem o fluxo natural de uso da plataforma — do cadastro até a avaliação.

O usuário enviará o padrão de design (UI/UX) após aprovar este plano de etapas.

---

## Estratégia de testes do frontend

Três camadas de teste, aplicadas em cada etapa:

| Camada | O que testa | Ferramentas | Quando usar |
|--------|------------|-------------|-------------|
| **Unitário** | Funções puras: validações, formatações, cálculos, utils | Vitest/Jest | Sempre — em toda etapa |
| **Componente** | Componente isolado: renderiza corretamente, exibe estados, responde a eventos | Testing Library | Componentes com lógica condicional |
| **Integração** | Fluxo completo: preencher form → submit → verificar resultado (com API mockada) | Testing Library + MSW | Fluxos críticos de negócio |

**Regra**: nenhuma etapa é considerada pronta sem os testes passando.

---

## Visão geral das etapas

| Etapa | Funcionalidade | Endpoints consumidos | Depende de |
|-------|---------------|---------------------|------------|
| 1 | Cadastro (Registro) | `POST /auth/registrar` | — |
| 2 | Login + Gerenciamento de tokens | `POST /auth/login`, `POST /auth/renovar` | Etapa 1 |
| 3 | Perfil do usuário (base + cliente + faxineiro) | `GET/PUT /usuarios/eu/perfil`, perfil-cliente, perfil-faxineiro | Etapa 2 |
| 4 | Catálogo de serviços (feed público) | `GET /limpezas`, `GET /limpezas/{id}`, `GET /feed` | Etapa 2 |
| 5 | Publicar serviço (faxineiro) | `POST /limpezas`, `PUT/DELETE /limpezas/{id}`, `GET /usuarios/eu/limpezas` | Etapa 3 |
| 6 | Agenda do faxineiro | `GET/POST/DELETE /agenda/disponibilidades`, `GET/POST/DELETE /agenda/bloqueios` | Etapa 5 |
| 7 | Solicitar serviço (cliente) | `POST /solicitacoes`, `GET /usuarios/eu/solicitacoes` | Etapa 4, 6 |
| 8 | Gerenciar solicitações (faxineiro) | `GET /limpezas/{id}/solicitacoes`, aceitar, rejeitar | Etapa 7 |
| 9 | Cancelar solicitação (cliente) | `POST /solicitacoes/{id}/cancelar` | Etapa 7 |
| 10 | Avaliar serviço + Reputação | `POST /avaliacoes`, `GET /faxineiros/{id}/avaliacoes`, `GET /faxineiros/{id}/estatisticas` | Etapa 8 |

---

## Etapa 1 — Cadastro (Registro de usuário)

### O que faz
Permite que um novo usuário crie conta na plataforma com email, nome de usuário e senha.

### Regras de negócio
- **Senha**: mínimo 8 caracteres, pelo menos 1 letra maiúscula e 1 dígito
- **Email**: único no sistema (backend retorna erro `409` se duplicado)
- **Nome de usuário**: único no sistema (backend retorna erro `409` se duplicado)
- Após registro bem-sucedido, o backend retorna `ParTokens` (token de acesso + token de renovação) — o usuário já fica logado

### Endpoint
```
POST /api/v1/auth/registrar
Body: { "email": "...", "nome_usuario": "...", "senha": "..." }
Response 201: { "usuario": {...}, "tokens": { "token_acesso", "token_renovacao", "tipo_token", "expira_em" } }
Response 409: email ou nome_usuario já existe
Response 422: senha fraca
```

### Tela
- Formulário com 3 campos: email, nome de usuário, senha
- Validação client-side da senha (8+ chars, 1 maiúscula, 1 dígito) antes de enviar
- Exibir erros do backend (email duplicado, nome duplicado, senha fraca)
- Após sucesso: salvar tokens, redirecionar para completar perfil (Etapa 3)
- Link "Já tem conta? Faça login" → tela de login (Etapa 2)

### Testes

**Unitários:**
- `validarSenha("Abc12345")` → válida
- `validarSenha("abc")` → inválida (sem maiúscula, sem dígito, curta)
- `validarSenha("abcdefgh")` → inválida (sem maiúscula, sem dígito)
- `validarSenha("ABCDEFGH")` → inválida (sem dígito)
- `validarSenha("Abcdefg1")` → válida (exatamente 8 chars)
- `validarEmail("teste@email.com")` → válido
- `validarEmail("invalido")` → inválido

**Componente:**
- Formulário renderiza 3 campos (email, nome de usuário, senha)
- Botão de submit desabilitado enquanto campos estão vazios
- Exibe mensagem de erro de senha fraca antes de enviar (validação client-side)
- Exibe mensagem de erro do backend (409: "Email já cadastrado")
- Exibe loading spinner durante requisição

**Integração:**
- Preencher formulário com dados válidos → submit → API retorna 201 → tokens salvos → redireciona para perfil
- Preencher com email duplicado → submit → API retorna 409 → exibe erro, não redireciona
- Preencher com senha fraca → validação client-side bloqueia submit

---

## Etapa 2 — Login + Gerenciamento de Tokens

### O que faz
Permite que usuário existente entre na plataforma e mantém a sessão ativa com refresh automático de tokens.

### Regras de negócio
- Login por **email + senha** (não por nome de usuário)
- Token de acesso expira em **15 minutos**
- Token de renovação expira em **7 dias**
- Se o token de acesso expirar, usar o de renovação para obter novo par
- Mensagem genérica "credenciais inválidas" (não revelar se é email ou senha errada — segurança contra enumeração)
- Usuário inativo recebe erro específico

### Endpoints
```
POST /api/v1/auth/login
Body: { "email": "...", "senha": "..." }
Response 200: { "usuario": {...}, "tokens": {...} }
Response 401: credenciais inválidas ou usuário inativo

POST /api/v1/auth/renovar
Body: { "token_renovacao": "..." }
Response 200: { "tokens": {...} }
Response 401: token de renovação inválido/expirado
```

### Tela e infraestrutura
- Formulário com 2 campos: email, senha
- Armazenar tokens de forma segura (httpOnly cookies ou localStorage com cuidado)
- Interceptor HTTP global:
  - Adiciona `Authorization: Bearer <token_acesso>` em toda requisição autenticada
  - Se receber `401`, tenta renovar com `POST /auth/renovar`
  - Se renovação falhar → redirecionar para login
- Link "Não tem conta? Cadastre-se" → tela de registro (Etapa 1)
- Botão de logout: limpar tokens e redirecionar para login

### Testes

**Unitários:**
- `salvarTokens(par)` → armazena corretamente no storage
- `obterTokenAcesso()` → retorna token salvo
- `limparTokens()` → remove ambos os tokens
- `tokenExpirado(token)` → true se expirou, false se válido
- `estaAutenticado()` → true se tem token válido

**Componente:**
- Formulário renderiza 2 campos (email, senha)
- Exibe "Credenciais inválidas" quando API retorna 401
- Exibe "Conta desativada" quando API retorna erro de usuário inativo
- Botão de logout limpa tokens e redireciona

**Integração:**
- Preencher email/senha válidos → submit → API retorna 200 → tokens salvos → redireciona para home
- Preencher credenciais erradas → submit → API retorna 401 → exibe erro genérico
- Token de acesso expira → interceptor chama /auth/renovar → novo token salvo → requisição original reenviada
- Token de renovação expirado → interceptor tenta renovar → falha → redireciona para login

---

## Etapa 3 — Perfil do Usuário

### O que faz
Permite que o usuário complete e edite seu perfil. O mesmo usuário pode ter **perfil de cliente E faxineiro** simultaneamente.

### Regras de negócio
- **Perfil base** (todos): nome completo, telefone, imagem (foto)
- **Perfil faxineiro** (opcional): descrição, anos de experiência, especialidades (lista), cidades atendidas (lista), documentos (RG, CPF, foto documento), verificação pela plataforma
- **Perfil cliente** (opcional): endereço (rua, bairro, cidade, estado, CEP), tipo de imóvel (apartamento/casa/comercial), quartos, banheiros, tamanho em m², observações, faxineiro preferido
- Especialidades aceitas: `limpeza_padrao`, `limpeza_pesada`, `limpeza_express`, `limpeza_pre_mudanca`, `limpeza_pos_obra`, `limpeza_comercial`, `passadoria`

### Endpoints
```
GET  /api/v1/usuarios/eu/perfil              → perfil base
PUT  /api/v1/usuarios/eu/perfil              → atualizar perfil base
POST /api/v1/usuarios/eu/perfil-faxineiro    → criar perfil faxineiro
GET  /api/v1/usuarios/eu/perfil-faxineiro    → buscar perfil faxineiro
PUT  /api/v1/usuarios/eu/perfil-faxineiro    → atualizar perfil faxineiro
POST /api/v1/usuarios/eu/perfil-cliente      → criar perfil cliente
GET  /api/v1/usuarios/eu/perfil-cliente      → buscar perfil cliente
PUT  /api/v1/usuarios/eu/perfil-cliente      → atualizar perfil cliente
```

### Telas
- **Tela de perfil base**: formulário com nome completo, telefone, upload de imagem
- **Seção "Quero ser faxineiro"**: formulário com descrição, experiência, seleção múltipla de especialidades, cidades atendidas, upload de documentos
- **Seção "Quero ser cliente"**: formulário com endereço completo, tipo de imóvel, detalhes do imóvel, observações
- Exibir badges visuais de qual(is) perfil(is) o usuário ativou

### Testes

**Unitários:**
- `formatarTelefone("11999999999")` → "(11) 99999-9999"
- `validarCEP("01310100")` → válido
- `validarCEP("123")` → inválido
- `traduzirTipoImovel("apartamento")` → "Apartamento"
- `traduzirEspecialidade("limpeza_pesada")` → "Limpeza Pesada"

**Componente:**
- Perfil base renderiza campos nome, telefone, imagem
- Seção faxineiro mostra seleção múltipla de especialidades (7 opções)
- Seção cliente mostra campos de endereço completo
- Badge "Faxineiro" aparece quando perfil faxineiro está ativo
- Badge "Cliente" aparece quando perfil cliente está ativo
- Formulário exibe dados existentes ao carregar (modo edição)

**Integração:**
- Preencher perfil base → salvar → API retorna 200 → dados atualizados na tela
- Criar perfil faxineiro com especialidades → API retorna 201 → badge faxineiro aparece
- Criar perfil cliente com endereço → API retorna 201 → badge cliente aparece

---

## Etapa 4 — Catálogo de Serviços (Feed Público)

### O que faz
Página principal pública onde qualquer pessoa (logada ou não) pode navegar pelos serviços de limpeza disponíveis.

### Regras de negócio
- Listagem paginada (padrão 20 itens por página, máximo 100)
- Cada serviço exibe: nome, descrição, tipo de limpeza, valor/hora, duração estimada, **preço total** (valor/hora × duração)
- Feed mostra eventos recentes (criação e atualização de serviços)
- 7 tipos de limpeza para filtrar: padrão, pesada, express, pré-mudança, pós-obra, comercial, passadoria

### Endpoints
```
GET /api/v1/limpezas?pagina=1&tamanho_pagina=20    → catálogo paginado
GET /api/v1/limpezas/{id}                            → detalhe do serviço
GET /api/v1/feed?pagina=1&tamanho_pagina=20          → feed de atividades
GET /api/v1/faxineiros/{id}/estatisticas             → nota média + total avaliações
GET /api/v1/faxineiros/{id}/avaliacoes               → lista de avaliações
```

### Telas
- **Página de catálogo**: grid/lista de cards com serviços, paginação, filtro por tipo
- **Página de detalhe do serviço**: todas as informações + avaliações do faxineiro + botão "Solicitar" (se logado como cliente)
- **Feed lateral ou seção**: últimas atividades da plataforma

### Testes

**Unitários:**
- `formatarMoeda(150.5)` → "R$ 150,50"
- `calcularPrecoTotal(50, 3)` → 150 (valor/hora × duração)
- `traduzirTipoLimpeza("limpeza_pos_obra")` → "Pós-Obra"
- `formatarDuracao(2.5)` → "2h30min"

**Componente:**
- Card de serviço exibe nome, tipo, preço total formatado
- Lista vazia exibe mensagem "Nenhum serviço encontrado"
- Filtro por tipo de limpeza renderiza 7 opções
- Paginação renderiza botões anterior/próximo
- Botão "Solicitar" aparece apenas se usuário logado como cliente
- Botão "Solicitar" não aparece se usuário não logado ou se é o faxineiro dono
- Página de detalhe exibe nota média e total de avaliações

**Integração:**
- Página carrega → API retorna lista → cards renderizados com dados corretos
- Clicar filtro "Limpeza Pesada" → nova requisição com filtro → cards atualizados
- Clicar próxima página → requisição com pagina=2 → novos cards

---

## Etapa 5 — Publicar Serviço (Faxineiro)

### O que faz
Permite que um faxineiro crie, edite e remova seus serviços de limpeza.

### Regras de negócio
- Apenas usuários com **perfil de faxineiro** podem publicar
- Campos obrigatórios: nome, descrição, valor por hora (> 0), duração estimada em horas (> 0), tipo de limpeza
- **Preço total = valor/hora × duração estimada** (calculado automaticamente, não editável)
- Faxineiro só pode editar/deletar seus **próprios** serviços (backend retorna `403` se tentar alterar serviço de outro)

### Endpoints
```
POST   /api/v1/limpezas                      → criar serviço
PUT    /api/v1/limpezas/{id}                  → atualizar serviço
DELETE /api/v1/limpezas/{id}                  → deletar serviço
GET    /api/v1/usuarios/eu/limpezas           → meus serviços publicados
```

### Telas
- **Dashboard do faxineiro**: lista dos serviços publicados com opções editar/excluir
- **Formulário criar/editar serviço**: nome, descrição, valor/hora, duração estimada, select de tipo de limpeza, preview do preço total calculado
- Confirmação antes de deletar

### Testes

**Unitários:**
- `calcularPrecoTotal(0, 3)` → erro (valor deve ser > 0)
- `calcularPrecoTotal(50, 0)` → erro (duração deve ser > 0)
- `calcularPrecoTotal(75, 2)` → 150
- `validarFormularioLimpeza({nome: "", ...})` → erro campo obrigatório

**Componente:**
- Formulário renderiza campos: nome, descrição, valor/hora, duração, tipo
- Preview do preço total atualiza em tempo real ao digitar valor/hora e duração
- Select de tipo de limpeza lista 7 opções
- Botão editar abre formulário com dados preenchidos
- Modal de confirmação aparece ao clicar excluir
- Dashboard lista apenas serviços do faxineiro logado

**Integração:**
- Preencher formulário válido → submit → API retorna 201 → serviço aparece na lista
- Editar serviço existente → submit → API retorna 200 → dados atualizados na lista
- Confirmar exclusão → API retorna 200 → serviço removido da lista
- Usuário sem perfil faxineiro → tela exibe mensagem para criar perfil primeiro

---

## Etapa 6 — Agenda do Faxineiro

### O que faz
Permite que o faxineiro gerencie sua disponibilidade semanal e bloqueios de horário.

### Regras de negócio
- **Disponibilidade**: blocos recorrentes por dia da semana (0=Domingo a 6=Sábado)
  - Cada bloco: dia da semana, hora início (0-23), hora fim (1-24)
  - Múltiplos blocos por dia permitidos (ex: manhã 8-12 e tarde 14-18)
  - Pode adicionar e remover livremente
- **Bloqueio pessoal**: período específico (data/hora início e fim) onde o faxineiro não está disponível
  - DataFim deve ser posterior a DataInicio
  - Não pode criar bloqueios no passado
  - Pode remover livremente
- **Bloqueio de serviço**: criado automaticamente pelo sistema quando uma solicitação é aceita (não aparece para criação manual, mas aparece na listagem com flag `e_pessoal: false`)
- A disponibilidade é verificada em **dois momentos**: quando o cliente solicita E quando o faxineiro aceita

### Endpoints
```
GET    /api/v1/agenda/disponibilidades           → listar disponibilidades
POST   /api/v1/agenda/disponibilidades           → adicionar disponibilidade
DELETE /api/v1/agenda/disponibilidades/{id}       → remover disponibilidade
GET    /api/v1/agenda/bloqueios                   → listar bloqueios (pessoais + serviço)
POST   /api/v1/agenda/bloqueios                   → criar bloqueio pessoal
DELETE /api/v1/agenda/bloqueios/{id}              → remover bloqueio pessoal
```

### Telas
- **Calendário semanal de disponibilidade**: visualização em grade (dias × horas), adicionar/remover blocos
- **Lista de bloqueios**: separar visualmente bloqueios pessoais (editáveis) e bloqueios de serviço (somente leitura, com link para a solicitação)
- Formulário para criar bloqueio pessoal com date-time picker

### Testes

**Unitários:**
- `traduzirDiaSemana(0)` → "Domingo"
- `traduzirDiaSemana(6)` → "Sábado"
- `formatarHorario(8, 12)` → "08:00 - 12:00"
- `validarBloqueio(inicio, fim)` → erro se fim <= inicio
- `validarBloqueio(passado, futuro)` → erro se inicio no passado
- `ePessoal(bloqueio)` → true se solicitacao_id é null

**Componente:**
- Grade semanal renderiza 7 colunas (Dom-Sáb) com slots de horário
- Blocos de disponibilidade existentes aparecem destacados na grade
- Bloqueio pessoal exibe botão remover; bloqueio de serviço exibe apenas "Serviço agendado"
- Date-time picker não permite selecionar datas no passado
- Formulário de bloqueio valida que fim > início antes de enviar

**Integração:**
- Adicionar disponibilidade (Segunda, 8-12) → API retorna 201 → bloco aparece na grade
- Remover disponibilidade → API retorna 200 → bloco some da grade
- Criar bloqueio pessoal → API retorna 201 → aparece na lista com badge "Pessoal"
- Listar bloqueios → API retorna pessoais + serviço → renderiza ambos com visual distinto

---

## Etapa 7 — Solicitar Serviço (Cliente)

### O que faz
Permite que um cliente solicite um serviço de limpeza publicado por um faxineiro.

### Regras de negócio
- Apenas usuários com **perfil de cliente** podem solicitar
- **Faxineiro NÃO pode solicitar seu próprio serviço** (backend retorna erro)
- Data agendada deve ser no **futuro**
- **Apenas 1 solicitação ativa** (pendente ou aceita) por cliente por serviço
- Preço total é **capturado no momento da criação** e não muda depois
- Verificação de disponibilidade automática:
  1. Disponibilidade semanal do faxineiro cobre a duração?
  2. Nenhum bloqueio (pessoal ou serviço) conflita com o horário?
- Status inicial: `pendente`
- Endereço pode ser informado manualmente ou copiado do perfil do cliente

### Endpoints
```
POST /api/v1/solicitacoes
Body: { "limpeza_id": 1, "data_agendada": "2026-04-01T10:00:00Z" }
Response 201: { solicitação completa com status "pendente" }
Response 409: já existe solicitação ativa para este serviço
Response 422: data no passado, sem disponibilidade, conflito de agenda

GET /api/v1/usuarios/eu/solicitacoes   → minhas solicitações como cliente
```

### Telas
- **Botão "Solicitar"** na página de detalhe do serviço (Etapa 4) → abre modal/página com:
  - Resumo do serviço (nome, preço total, duração)
  - Date-time picker para data agendada
  - Endereço (pré-preenchido do perfil ou editar)
  - Confirmação antes de enviar
- **Minhas solicitações (cliente)**: lista com status visual (pendente=amarelo, aceita=verde, rejeitada=vermelho, cancelada=cinza, concluída=azul)

### Testes

**Unitários:**
- `dataNoFuturo("2020-01-01")` → false
- `dataNoFuturo("2030-01-01")` → true
- `corDoStatus("pendente")` → "amarelo"
- `corDoStatus("aceita")` → "verde"
- `corDoStatus("rejeitada")` → "vermelho"
- `corDoStatus("cancelada")` → "cinza"
- `corDoStatus("concluida")` → "azul"
- `formatarDataAgendada("2026-04-01T10:00:00Z")` → "01/04/2026 às 10:00"

**Componente:**
- Modal de solicitação exibe resumo do serviço (nome, preço, duração)
- Date-time picker não permite selecionar data no passado
- Endereço pré-preenchido do perfil do cliente (se existir)
- Lista de solicitações renderiza badges coloridos por status
- Botão "Solicitar" desabilitado durante loading

**Integração:**
- Preencher data futura → submit → API retorna 201 → solicitação aparece como "pendente" na lista
- Tentar solicitar novamente o mesmo serviço → API retorna 409 → exibe "Já existe solicitação ativa"
- Selecionar data sem disponibilidade → API retorna 422 → exibe "Faxineiro indisponível neste horário"

---

## Etapa 8 — Gerenciar Solicitações (Faxineiro)

### O que faz
Permite que o faxineiro veja, aceite ou rejeite solicitações recebidas para seus serviços.

### Regras de negócio
- Apenas o faxineiro **dono do serviço** pode aceitar/rejeitar
- Só pode aceitar/rejeitar solicitações com status `pendente`
- Ao **aceitar**:
  - Sistema verifica disponibilidade **novamente** (pode ter mudado desde a criação)
  - Se disponível → status muda para `aceita` + cria bloqueio de serviço na agenda
  - Se indisponível → erro, solicitação permanece `pendente`
- Ao **rejeitar**: status muda para `rejeitada` (terminal, sem reversão)

### Endpoints
```
GET  /api/v1/limpezas/{limpeza_id}/solicitacoes                      → solicitações do meu serviço
POST /api/v1/solicitacoes/{cliente_id}/{limpeza_id}/aceitar           → aceitar
POST /api/v1/solicitacoes/{cliente_id}/{limpeza_id}/rejeitar          → rejeitar
```

### Telas
- **Painel de solicitações do faxineiro**: agrupado por serviço, com badges de status
- Cada solicitação pendente mostra: nome do cliente, data agendada, endereço, preço total + botões Aceitar/Rejeitar
- Confirmação antes de aceitar/rejeitar
- Feedback visual claro quando aceite falha por indisponibilidade de agenda

### Testes

**Unitários:**
- `podeTerAcao("pendente")` → true (aceitar/rejeitar)
- `podeTerAcao("aceita")` → false
- `podeTerAcao("rejeitada")` → false
- `podeTerAcao("cancelada")` → false

**Componente:**
- Solicitação pendente exibe botões Aceitar e Rejeitar
- Solicitação aceita/rejeitada/cancelada NÃO exibe botões de ação
- Modal de confirmação aparece ao clicar Aceitar ("Confirma aceitar?")
- Modal de confirmação aparece ao clicar Rejeitar ("Confirma rejeitar? Essa ação não pode ser desfeita")
- Exibe mensagem de erro quando aceite falha por conflito de agenda

**Integração:**
- Clicar Aceitar → confirmar → API retorna 200 → status muda para "aceita" na lista
- Clicar Rejeitar → confirmar → API retorna 200 → status muda para "rejeitada" na lista
- Clicar Aceitar → API retorna 422 (agenda indisponível) → exibe erro, status permanece "pendente"

---

## Etapa 9 — Cancelar Solicitação (Cliente)

### O que faz
Permite que o cliente cancele uma solicitação pendente ou aceita, com regras de multa.

### Regras de negócio
- **Solicitação pendente**: cancelamento sem custo
- **Solicitação aceita**:
  - Se faltam **≥ 24 horas** para o serviço → cancelamento sem custo
  - Se faltam **< 24 horas** → multa de **20% do preço total**
- Cancelar solicitação aceita → remove automaticamente o bloqueio de serviço da agenda do faxineiro
- Apenas o **cliente que criou** pode cancelar
- Status `rejeitada`, `cancelada` e `concluída` não podem ser cancelados

### Endpoint
```
POST /api/v1/solicitacoes/{limpeza_id}/cancelar
Response 200: { solicitação com status "cancelada", multa_cancelamento: 0 ou valor }
Response 403: não é o cliente da solicitação
Response 422: status não permite cancelamento
```

### Tela
- Botão "Cancelar" na lista de solicitações do cliente (Etapa 7)
- **Modal de confirmação** que mostra:
  - Se pendente: "Cancelamento gratuito"
  - Se aceita com ≥ 24h: "Cancelamento gratuito"
  - Se aceita com < 24h: "Atenção: multa de 20% (R$ X,XX) será aplicada"
- Após cancelamento: atualizar status visual na lista

### Testes

**Unitários:**
- `calcularMulta(200, "pendente", qualquerData)` → 0
- `calcularMulta(200, "aceita", dataEm25horas)` → 0 (>= 24h)
- `calcularMulta(200, "aceita", dataEm23horas)` → 40 (20% de 200)
- `calcularMulta(200, "aceita", dataEm1hora)` → 40
- `podeCancelar("pendente")` → true
- `podeCancelar("aceita")` → true
- `podeCancelar("rejeitada")` → false
- `podeCancelar("concluida")` → false
- `podeCancelar("cancelada")` → false
- `formatarMulta(40)` → "R$ 40,00"
- `mensagemCancelamento("pendente", 0)` → "Cancelamento gratuito"
- `mensagemCancelamento("aceita", 40)` → "Atenção: multa de 20% (R$ 40,00) será aplicada"

**Componente:**
- Botão "Cancelar" aparece apenas em solicitações pendentes ou aceitas
- Botão "Cancelar" NÃO aparece em rejeitada, cancelada, concluída
- Modal exibe "Cancelamento gratuito" para solicitação pendente
- Modal exibe aviso de multa com valor calculado para aceita < 24h
- Modal exibe "Cancelamento gratuito" para aceita ≥ 24h

**Integração:**
- Cancelar solicitação pendente → API retorna 200 com multa=0 → status muda para "cancelada"
- Cancelar solicitação aceita < 24h → API retorna 200 com multa=40 → exibe valor da multa → status "cancelada"

---

## Etapa 10 — Avaliar Serviço + Reputação

### O que faz
Permite que o cliente avalie um serviço aceito, marcando-o como concluído. Exibe reputação pública do faxineiro.

### Regras de negócio
- Apenas solicitações com status `aceita` podem ser avaliadas
- **Uma avaliação por solicitação** (backend retorna `409` se duplicada)
- Nota: inteiro de **0 a 5** (obrigatória)
- Comentário: texto livre (opcional)
- Criar avaliação → solicitação muda automaticamente para `concluída`
- Estatísticas do faxineiro são públicas: média das notas + total de avaliações

### Endpoints
```
POST /api/v1/avaliacoes
Body: { "limpeza_id": 1, "nota": 5, "comentario": "Excelente!" }
Response 201: { avaliação criada }
Response 409: já avaliou esta solicitação

GET /api/v1/faxineiros/{faxineiro_id}/avaliacoes      → lista de avaliações
GET /api/v1/faxineiros/{faxineiro_id}/estatisticas     → média + total
```

### Telas
- **Botão "Avaliar"** aparece apenas em solicitações com status `aceita`
- **Modal/página de avaliação**: estrelas (0-5) + campo de comentário + botão enviar
- **Seção de reputação** no perfil público do faxineiro: nota média, total de avaliações, lista de comentários
- Exibir reputação também na página de detalhe do serviço (Etapa 4)

### Testes

**Unitários:**
- `validarNota(0)` → válida
- `validarNota(5)` → válida
- `validarNota(6)` → inválida
- `validarNota(-1)` → inválida
- `formatarMediaNota(4.333)` → "4.3"
- `renderizarEstrelas(4)` → "★★★★☆"

**Componente:**
- Botão "Avaliar" aparece apenas em solicitações com status `aceita`
- Botão "Avaliar" NÃO aparece em pendente, rejeitada, cancelada, concluída
- Componente de estrelas renderiza 5 estrelas clicáveis
- Clicar na estrela 3 → seleciona notas 1, 2, 3 (preenchidas)
- Campo de comentário é opcional (formulário válido sem ele)
- Seção de estatísticas exibe média formatada e total
- Lista de avaliações exibe nota em estrelas + comentário + nome do cliente

**Integração:**
- Selecionar nota 5 + comentário → submit → API retorna 201 → solicitação muda para "concluída"
- Tentar avaliar novamente → API retorna 409 → exibe "Você já avaliou este serviço"
- Página do faxineiro carrega → API retorna estatísticas → exibe média e total corretos

---

## Melhorias transversais (aplicadas ao longo das etapas)

### UX — Onboarding guiado (Etapa 1-3)
Após cadastro, wizard de 2-3 passos: "Você quer contratar ou prestar serviços?" → direciona para criar perfil cliente ou faxineiro. Evita que o usuário fique perdido após o registro.

**Testes:**
- **Componente:** wizard renderiza passo 1 ("Contratar" / "Prestar serviços"), clicar "Prestar serviços" avança para formulário de faxineiro
- **Integração:** completar wizard → perfil criado via API → redireciona para dashboard correto

### UX — Notificações em tempo real (Etapa 7-8)
Quando o faxineiro aceita/rejeita uma solicitação, o cliente vê a atualização sem recarregar. WebSocket ou polling curto para atualizar status.

**Testes:**
- **Unitário:** `processarEvento({ tipo: "solicitacao_aceita", id: 1 })` → atualiza status local
- **Componente:** badge de status muda de "pendente" para "aceita" quando evento chega
- **Integração:** simular evento WebSocket → lista de solicitações atualiza automaticamente

### UX — Busca com filtros combinados (Etapa 4)
No catálogo, filtrar por: cidade, faixa de preço (mín-máx), nota mínima do faxineiro, tipo de limpeza. Filtros combinam entre si.

**Testes:**
- **Unitário:** `construirQueryFiltros({ cidade: "SP", precoMin: 50, tipoLimpeza: "pesada" })` → query string correta
- **Componente:** painel de filtros renderiza campos de cidade, faixa de preço, nota mínima, tipo
- **Integração:** aplicar filtro cidade + tipo → API chamada com parâmetros combinados → cards filtrados

### Segurança — Proteção de rotas (Etapa 2)
Componente de guarda que redireciona:
- Para login se não autenticado
- Para "criar perfil faxineiro" se tentar acessar área de faxineiro sem perfil
- Para "criar perfil cliente" se tentar acessar área de cliente sem perfil

**Testes:**
- **Componente:** usuário não autenticado tenta acessar `/agenda` → renderiza redirect para `/login`
- **Componente:** usuário sem perfil faxineiro acessa `/meus-servicos` → renderiza redirect para criar perfil
- **Integração:** navegar para rota protegida sem token → redireciona para login → fazer login → redireciona de volta para rota original

### Segurança — Rate limiting visual (transversal)
Se o backend retornar `429 Too Many Requests`, exibir mensagem amigável: "Muitas tentativas. Aguarde um momento e tente novamente."

**Testes:**
- **Componente:** quando API retorna 429 → exibe mensagem amigável, não erro genérico
- **Componente:** mensagem desaparece após tempo configurável

### Performance — Cache de dados estáticos (Etapa 3-4)
Tipos de limpeza (7), especialidades (7), dias da semana, tipos de imóvel — cachear no client. Não recarregar a cada navegação.

**Testes:**
- **Unitário:** `obterTiposLimpeza()` → retorna do cache na segunda chamada (sem fetch)

### Performance — Skeleton loading (transversal)
Em vez de spinner genérico, mostrar esqueletos (placeholders) dos cards/formulários enquanto carrega.

**Testes:**
- **Componente:** durante loading, renderiza skeleton com mesma estrutura do card final
- **Componente:** após loading, skeleton é substituído pelo conteúdo real

### Performance — Paginação infinita no feed (Etapa 4)
Scroll infinito no feed de atividades em vez de botões de página. Carrega próxima página ao atingir o fim da lista.

**Testes:**
- **Integração:** scroll até o fim → API chamada com pagina=2 → novos itens adicionados à lista (sem substituir)
- **Componente:** exibe "Carregando mais..." durante fetch da próxima página
- **Componente:** exibe "Não há mais atividades" quando última página

### Negócio — Resumo financeiro do faxineiro (Etapa 8)
Dashboard com: total ganho (soma dos preços de solicitações concluídas), serviços realizados (count), nota média, tudo em cards resumo no topo do painel.

**Testes:**
- **Unitário:** `calcularTotalGanho([{preco: 100, status: "concluida"}, {preco: 200, status: "cancelada"}])` → 100 (só concluídas)
- **Unitário:** `contarServicosRealizados(solicitacoes)` → conta apenas status "concluida"
- **Componente:** dashboard exibe 3 cards: "Total ganho: R$ X", "Serviços realizados: N", "Nota média: X.X"

### Negócio — Histórico de solicitações com abas (Etapa 7-8)
Separar "Ativas" (pendente + aceita) e "Finalizadas" (concluída + cancelada + rejeitada) em abas. Facilita encontrar o que importa.

**Testes:**
- **Unitário:** `filtrarAtivas(solicitacoes)` → retorna apenas pendentes e aceitas
- **Unitário:** `filtrarFinalizadas(solicitacoes)` → retorna concluídas, canceladas, rejeitadas
- **Componente:** aba "Ativas" renderiza apenas solicitações pendentes/aceitas
- **Componente:** aba "Finalizadas" renderiza apenas concluídas/canceladas/rejeitadas
- **Componente:** badge na aba "Ativas" mostra contagem de pendentes

---

## Estrutura técnica transversal (definida com o design)

Itens que serão decididos quando o usuário enviar o padrão de design:

- **Framework/biblioteca** (React, Next.js, Vue, etc.)
- **Gerenciamento de estado** (Context API, Zustand, Redux, etc.)
- **Estilização** (Tailwind, CSS Modules, Styled Components, etc.)
- **Estrutura de pastas** do frontend
- **Componentes reutilizáveis** (botões, cards, modals, forms)
- **Tema visual** (cores, tipografia, espaçamento)
- **Biblioteca de testes** (Vitest + Testing Library + MSW recomendados)

---

## Verificação por etapa

Cada etapa será considerada pronta quando:

1. Tela funcional conectada ao backend real
2. Tratamento de todos os erros do backend (409, 401, 403, 422, 429)
3. Validação client-side onde aplicável
4. Feedback visual para loading (skeleton), sucesso e erro
5. Responsividade (mobile-first)
6. **Todos os testes passando** (unitários + componente + integração da etapa)
7. Documentação atualizada (DOCUMENTACAO.md + README.md)

---

## Alteração no CLAUDE.md

Adicionar seção **"Convenções de teste do frontend"** com as regras:

- Três camadas obrigatórias: unitário, componente, integração
- Testes unitários para toda função pura (validação, formatação, cálculo)
- Testes de componente para todo componente com lógica condicional
- Testes de integração para fluxos críticos (cadastro, login, solicitar, cancelar, avaliar)
- Rodar `npm test` (ou equivalente) antes de qualquer commit do frontend
- Nomes de teste descritivos alinhados com a documentação (mesmo padrão do backend)
