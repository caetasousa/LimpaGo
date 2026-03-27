package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"

	"limpaGo/api/auth"
	"limpaGo/api/dto"
	"limpaGo/api/handler"
	"limpaGo/api/middleware"
	"limpaGo/domain/entity"
	"limpaGo/domain/service"
	"limpaGo/domain/testutil"
	"limpaGo/domain/valueobject"
)

// --- helpers de teste ---

type ambienteTeste struct {
	repoUsuarios     *testutil.RepositorioUsuarioMock
	repoPerfis       *testutil.RepositorioPerfilMock
	repoLimpezas     *testutil.RepositorioLimpezaMock
	repoSolicitacoes *testutil.RepositorioSolicitacaoMock
	repoAgenda       *testutil.RepositorioAgendaMock
	repoAvaliacoes   *testutil.RepositorioAvaliacaoMock
	repoFeed         *testutil.RepositorioFeedMock

	svcUsuario     *service.ServicoUsuario
	svcLimpeza     *service.ServicoLimpeza
	svcAgenda      *service.ServicoAgenda
	svcSolicitacao *service.ServicoSolicitacao
	svcAvaliacao   *service.ServicoAvaliacao
	svcFeed        *service.ServicoFeed
	svcAuth  *auth.ServicoAutenticacao
	svcToken *auth.ServicoToken

	handlerAuth        *handler.HandlerAutenticacao
	handlerUsuario     *handler.HandlerUsuario
	handlerLimpeza     *handler.HandlerLimpeza
	handlerSolicitacao *handler.HandlerSolicitacao
	handlerAgenda      *handler.HandlerAgenda
	handlerAvaliacao   *handler.HandlerAvaliacao
	handlerFeed        *handler.HandlerFeed
}

func novoAmbienteTeste(t *testing.T) *ambienteTeste {
	t.Helper()

	repoUsuarios := testutil.NovoRepositorioUsuarioMock()
	repoPerfis := testutil.NovoRepositorioPerfilMock()
	repoLimpezas := testutil.NovoRepositorioLimpezaMock()
	repoSolicitacoes := testutil.NovoRepositorioSolicitacaoMock()
	repoAgenda := testutil.NovoRepositorioAgendaMock()
	repoAvaliacoes := testutil.NovoRepositorioAvaliacaoMock()
	repoFeed := testutil.NovoRepositorioFeedMock()

	svcUsuario := service.NovoServicoUsuario(repoUsuarios, repoPerfis)
	svcLimpeza := service.NovoServicoLimpeza(repoLimpezas)
	svcAgenda := service.NovoServicoAgenda(repoAgenda)
	svcSolicitacao := service.NovoServicoSolicitacao(repoSolicitacoes, repoLimpezas, svcAgenda)
	svcAvaliacao := service.NovoServicoAvaliacao(repoAvaliacoes, repoSolicitacoes, repoLimpezas)
	svcFeed := service.NovoServicoFeed(repoFeed)

	repoCredenciais := auth.NovoRepositorioCredencialMock()
	cfgJWT := auth.ConfiguracaoPadrao()
	svcToken := auth.NovoServicoToken(cfgJWT)
	svcAuth := auth.NovoServicoAutenticacao(repoUsuarios, repoCredenciais, svcUsuario, svcToken)

	return &ambienteTeste{
		repoUsuarios:     repoUsuarios,
		repoPerfis:       repoPerfis,
		repoLimpezas:     repoLimpezas,
		repoSolicitacoes: repoSolicitacoes,
		repoAgenda:       repoAgenda,
		repoAvaliacoes:   repoAvaliacoes,
		repoFeed:         repoFeed,

		svcUsuario:     svcUsuario,
		svcLimpeza:     svcLimpeza,
		svcAgenda:      svcAgenda,
		svcSolicitacao: svcSolicitacao,
		svcAvaliacao:   svcAvaliacao,
		svcFeed:        svcFeed,
		svcAuth:  svcAuth,
		svcToken: svcToken,

		handlerAuth:        handler.NovoHandlerAutenticacao(svcAuth),
		handlerUsuario:     handler.NovoHandlerUsuario(svcUsuario),
		handlerLimpeza:     handler.NovoHandlerLimpeza(svcLimpeza),
		handlerSolicitacao: handler.NovoHandlerSolicitacao(svcSolicitacao),
		handlerAgenda:      handler.NovoHandlerAgenda(svcAgenda),
		handlerAvaliacao:   handler.NovoHandlerAvaliacao(svcAvaliacao),
		handlerFeed:        handler.NovoHandlerFeed(svcFeed),
	}
}

func (a *ambienteTeste) registrarUsuario(t *testing.T, email, nome string) (*entity.Usuario, string) {
	t.Helper()
	// No ambiente de teste, registra diretamente via ServicoUsuario (sem Zitadel)
	usuario, err := a.svcUsuario.Registrar(context.Background(), email, nome)
	if err != nil {
		t.Fatalf("erro ao registrar usuario: %v", err)
	}
	// Injeta um token fictício (os testes de handler usam contexto direto)
	return usuario, "token-teste-mock"
}

func reqComToken(method, url string, body interface{}, token string) *http.Request {
	var buf bytes.Buffer
	if body != nil {
		_ = json.NewEncoder(&buf).Encode(body)
	}
	req := httptest.NewRequest(method, url, &buf)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return req
}

func reqComContextoUsuario(req *http.Request, usuarioID int) *http.Request {
	ctx := context.WithValue(req.Context(), middleware.ChaveUsuarioID, usuarioID)
	return req.WithContext(ctx)
}

func reqComChiParams(req *http.Request, params map[string]string) *http.Request {
	rctx := chi.NewRouteContext()
	for k, v := range params {
		rctx.URLParams.Add(k, v)
	}
	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
}

// --- Testes de Autenticação ---

func TestRegistro_email_valido_retorna_201(t *testing.T) {
	t.Parallel()
	amb := novoAmbienteTeste(t)
	body := dto.RequisicaoRegistroComSenha{
		Email: "novo@email.com", NomeUsuario: "novousuario", Senha: "Senha123forte",
	}
	req := reqComToken(http.MethodPost, "/auth/registrar", body, "")
	rec := httptest.NewRecorder()

	amb.handlerAuth.Registrar(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("got %d; want %d", rec.Code, http.StatusCreated)
	}
}

func TestLogin_email_inexistente_retorna_401(t *testing.T) {
	t.Parallel()
	amb := novoAmbienteTeste(t)

	body := dto.RequisicaoLogin{Email: "naoexiste@email.com", Senha: "Senha123forte"}
	req := reqComToken(http.MethodPost, "/auth/login", body, "")
	rec := httptest.NewRecorder()

	amb.handlerAuth.Login(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("got %d; want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestRenovarToken_token_invalido_retorna_401(t *testing.T) {
	t.Parallel()
	amb := novoAmbienteTeste(t)

	body := dto.RequisicaoRenovarToken{TokenRenovacao: "token-invalido"}
	req := reqComToken(http.MethodPost, "/auth/renovar", body, "")
	rec := httptest.NewRecorder()

	amb.handlerAuth.RenovarToken(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("got %d; want %d", rec.Code, http.StatusUnauthorized)
	}
}

// --- Testes de Limpeza (Serviço de Limpeza) ---

func TestCriarLimpeza_profissional_publica_servico_retorna_201(t *testing.T) {
	t.Parallel()
	amb := novoAmbienteTeste(t)
	usuario, _ := amb.registrarUsuario(t, "fax@email.com", "profissional1")

	body := dto.RequisicaoCriarLimpeza{
		Nome: "Limpeza Residencial", Descricao: "Limpeza completa",
		ValorHora: 50.0, DuracaoEstimada: 3.0, TipoLimpeza: "limpeza_padrao",
	}
	req := reqComToken(http.MethodPost, "/limpezas", body, "")
	req = reqComContextoUsuario(req, usuario.ID)
	rec := httptest.NewRecorder()

	amb.handlerLimpeza.CriarLimpeza(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("got %d; want %d", rec.Code, http.StatusCreated)
	}
	var resp dto.RespostaLimpeza
	_ = json.NewDecoder(rec.Body).Decode(&resp)
	if resp.PrecoTotal != 150.0 {
		t.Errorf("got preco_total %.2f; want 150.00", resp.PrecoTotal)
	}
}

func TestCriarLimpeza_sem_autenticacao_retorna_401(t *testing.T) {
	t.Parallel()
	amb := novoAmbienteTeste(t)
	body := dto.RequisicaoCriarLimpeza{
		Nome: "Limpeza", ValorHora: 50.0, DuracaoEstimada: 2.0, TipoLimpeza: "limpeza_padrao",
	}
	req := reqComToken(http.MethodPost, "/limpezas", body, "")
	rec := httptest.NewRecorder()

	amb.handlerLimpeza.CriarLimpeza(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("got %d; want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestListarCatalogo_retorna_lista_de_servicos(t *testing.T) {
	t.Parallel()
	amb := novoAmbienteTeste(t)
	_, _ = amb.svcLimpeza.Criar(context.Background(), 1, "Limpeza 1", "desc", 30.0, 2.0, valueobject.TipoLimpezaPadrao)
	_, _ = amb.svcLimpeza.Criar(context.Background(), 1, "Limpeza 2", "desc", 40.0, 1.5, valueobject.TipoLimpezaComercial)

	req := httptest.NewRequest(http.MethodGet, "/limpezas?pagina=1&tamanho=10", nil)
	rec := httptest.NewRecorder()

	amb.handlerLimpeza.ListarCatalogo(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("got %d; want %d", rec.Code, http.StatusOK)
	}
	var lista []dto.RespostaLimpeza
	_ = json.NewDecoder(rec.Body).Decode(&lista)
	if len(lista) != 2 {
		t.Errorf("got %d items; want 2", len(lista))
	}
}

func TestBuscarLimpeza_por_id_retorna_servico(t *testing.T) {
	t.Parallel()
	amb := novoAmbienteTeste(t)
	limpeza, _ := amb.svcLimpeza.Criar(context.Background(), 1, "Limpeza", "desc", 30.0, 2.0, valueobject.TipoLimpezaPadrao)

	req := httptest.NewRequest(http.MethodGet, "/limpezas/1", nil)
	req = reqComChiParams(req, map[string]string{"id": fmt.Sprintf("%d", limpeza.ID)})
	rec := httptest.NewRecorder()

	amb.handlerLimpeza.BuscarLimpeza(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("got %d; want %d", rec.Code, http.StatusOK)
	}
}

func TestBuscarLimpeza_inexistente_retorna_404(t *testing.T) {
	t.Parallel()
	amb := novoAmbienteTeste(t)

	req := httptest.NewRequest(http.MethodGet, "/limpezas/999", nil)
	req = reqComChiParams(req, map[string]string{"id": "999"})
	rec := httptest.NewRecorder()

	amb.handlerLimpeza.BuscarLimpeza(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("got %d; want %d", rec.Code, http.StatusNotFound)
	}
}

func TestDeletarLimpeza_profissional_remove_seu_servico(t *testing.T) {
	t.Parallel()
	amb := novoAmbienteTeste(t)
	usuario, _ := amb.registrarUsuario(t, "del@email.com", "deletar")
	limpeza, _ := amb.svcLimpeza.Criar(context.Background(), usuario.ID, "Temp", "desc", 25.0, 1.0, valueobject.TipoLimpezaPadrao)

	req := httptest.NewRequest(http.MethodDelete, "/limpezas/1", nil)
	req = reqComContextoUsuario(req, usuario.ID)
	req = reqComChiParams(req, map[string]string{"id": fmt.Sprintf("%d", limpeza.ID)})
	rec := httptest.NewRecorder()

	amb.handlerLimpeza.DeletarLimpeza(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("got %d; want %d", rec.Code, http.StatusNoContent)
	}
}

func TestListarMinhasLimpezas_profissional_ve_seus_servicos(t *testing.T) {
	t.Parallel()
	amb := novoAmbienteTeste(t)
	usuario, _ := amb.registrarUsuario(t, "fax2@email.com", "profissional2")
	_, _ = amb.svcLimpeza.Criar(context.Background(), usuario.ID, "Minha Limpeza", "desc", 30.0, 2.0, valueobject.TipoLimpezaPadrao)

	req := httptest.NewRequest(http.MethodGet, "/usuarios/eu/limpezas", nil)
	req = reqComContextoUsuario(req, usuario.ID)
	rec := httptest.NewRecorder()

	amb.handlerLimpeza.ListarMinhasLimpezas(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("got %d; want %d", rec.Code, http.StatusOK)
	}
}

// --- Testes de Perfil ---

func TestBuscarPerfil_usuario_autenticado_ve_seu_perfil(t *testing.T) {
	t.Parallel()
	amb := novoAmbienteTeste(t)
	usuario, _ := amb.registrarUsuario(t, "perfil@email.com", "perfiluser")

	req := httptest.NewRequest(http.MethodGet, "/usuarios/eu/perfil", nil)
	req = reqComContextoUsuario(req, usuario.ID)
	rec := httptest.NewRecorder()

	amb.handlerUsuario.BuscarMeuPerfil(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("got %d; want %d", rec.Code, http.StatusOK)
	}
}

func TestBuscarPerfil_sem_autenticacao_retorna_401(t *testing.T) {
	t.Parallel()
	amb := novoAmbienteTeste(t)

	req := httptest.NewRequest(http.MethodGet, "/usuarios/eu/perfil", nil)
	rec := httptest.NewRecorder()

	amb.handlerUsuario.BuscarMeuPerfil(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("got %d; want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestAtualizarPerfil_usuario_atualiza_nome_e_telefone(t *testing.T) {
	t.Parallel()
	amb := novoAmbienteTeste(t)
	usuario, _ := amb.registrarUsuario(t, "att@email.com", "attuser")

	body := dto.RequisicaoAtualizarPerfil{
		NomeCompleto: "Nome Completo", Telefone: "11999999999", Imagem: "foto.jpg",
	}
	req := reqComToken(http.MethodPut, "/usuarios/eu/perfil", body, "")
	req = reqComContextoUsuario(req, usuario.ID)
	rec := httptest.NewRecorder()

	amb.handlerUsuario.AtualizarMeuPerfil(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("got %d; want %d", rec.Code, http.StatusOK)
	}
	var resp dto.RespostaPerfil
	_ = json.NewDecoder(rec.Body).Decode(&resp)
	if resp.NomeCompleto != "Nome Completo" {
		t.Errorf("got nome %q; want %q", resp.NomeCompleto, "Nome Completo")
	}
}

func TestCriarPerfilProfissional_usuario_se_torna_profissional(t *testing.T) {
	t.Parallel()
	amb := novoAmbienteTeste(t)
	usuario, _ := amb.registrarUsuario(t, "faxp@email.com", "faxprof")

	req := httptest.NewRequest(http.MethodPost, "/usuarios/eu/perfil-profissional", nil)
	req = reqComContextoUsuario(req, usuario.ID)
	rec := httptest.NewRecorder()

	amb.handlerUsuario.CriarPerfilProfissional(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("got %d; want %d", rec.Code, http.StatusCreated)
	}
}

func TestCriarPerfilCliente_usuario_se_torna_cliente(t *testing.T) {
	t.Parallel()
	amb := novoAmbienteTeste(t)
	usuario, _ := amb.registrarUsuario(t, "clip@email.com", "cliprof")

	req := httptest.NewRequest(http.MethodPost, "/usuarios/eu/perfil-cliente", nil)
	req = reqComContextoUsuario(req, usuario.ID)
	rec := httptest.NewRecorder()

	amb.handlerUsuario.CriarPerfilCliente(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("got %d; want %d", rec.Code, http.StatusCreated)
	}
}

// --- Testes de Agenda ---

func TestAdicionarDisponibilidade_profissional_define_horario(t *testing.T) {
	t.Parallel()
	amb := novoAmbienteTeste(t)
	usuario, _ := amb.registrarUsuario(t, "ag@email.com", "aguser")

	body := dto.RequisicaoDisponibilidade{DiaSemana: 1, HoraInicio: 8, HoraFim: 17}
	req := reqComToken(http.MethodPost, "/agenda/disponibilidades", body, "")
	req = reqComContextoUsuario(req, usuario.ID)
	rec := httptest.NewRecorder()

	amb.handlerAgenda.AdicionarDisponibilidade(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("got %d; want %d", rec.Code, http.StatusCreated)
	}
}

func TestListarDisponibilidade_profissional_ve_seus_horarios(t *testing.T) {
	t.Parallel()
	amb := novoAmbienteTeste(t)
	usuario, _ := amb.registrarUsuario(t, "ld@email.com", "lduser")
	_, _ = amb.svcAgenda.AdicionarDisponibilidade(context.Background(), usuario.ID, time.Monday, 8, 12)

	req := httptest.NewRequest(http.MethodGet, "/agenda/disponibilidades", nil)
	req = reqComContextoUsuario(req, usuario.ID)
	rec := httptest.NewRecorder()

	amb.handlerAgenda.ListarDisponibilidade(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("got %d; want %d", rec.Code, http.StatusOK)
	}
}

func TestCriarBloqueioPessoal_profissional_bloqueia_horario(t *testing.T) {
	t.Parallel()
	amb := novoAmbienteTeste(t)
	usuario, _ := amb.registrarUsuario(t, "bl@email.com", "bluser")

	body := dto.RequisicaoBloqueio{
		DataInicio: time.Now().Add(24 * time.Hour),
		DataFim:    time.Now().Add(26 * time.Hour),
	}
	req := reqComToken(http.MethodPost, "/agenda/bloqueios", body, "")
	req = reqComContextoUsuario(req, usuario.ID)
	rec := httptest.NewRecorder()

	amb.handlerAgenda.CriarBloqueioPessoal(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("got %d; want %d", rec.Code, http.StatusCreated)
	}
}

func TestListarBloqueios_profissional_ve_seus_bloqueios(t *testing.T) {
	t.Parallel()
	amb := novoAmbienteTeste(t)
	usuario, _ := amb.registrarUsuario(t, "lb@email.com", "lbuser")
	_, _ = amb.svcAgenda.CriarBloqueioPessoal(context.Background(), usuario.ID, time.Now().Add(24*time.Hour), time.Now().Add(26*time.Hour))

	req := httptest.NewRequest(http.MethodGet, "/agenda/bloqueios", nil)
	req = reqComContextoUsuario(req, usuario.ID)
	rec := httptest.NewRecorder()

	amb.handlerAgenda.ListarBloqueios(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("got %d; want %d", rec.Code, http.StatusOK)
	}
}

// --- Testes de Feed ---

func TestBuscarFeed_retorna_pagina_de_atividades(t *testing.T) {
	t.Parallel()
	amb := novoAmbienteTeste(t)

	req := httptest.NewRequest(http.MethodGet, "/feed?pagina=1&tamanho=10", nil)
	rec := httptest.NewRecorder()

	amb.handlerFeed.BuscarFeed(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("got %d; want %d", rec.Code, http.StatusOK)
	}
}

// --- Testes de Avaliação ---

func TestListarAvaliacoes_de_profissional_retorna_200(t *testing.T) {
	t.Parallel()
	amb := novoAmbienteTeste(t)

	req := httptest.NewRequest(http.MethodGet, "/profissionais/1/avaliacoes", nil)
	req = reqComChiParams(req, map[string]string{"profissional_id": "1"})
	rec := httptest.NewRecorder()

	amb.handlerAvaliacao.ListarAvaliacoes(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("got %d; want %d", rec.Code, http.StatusOK)
	}
}

func TestBuscarEstatisticas_de_profissional_retorna_200(t *testing.T) {
	t.Parallel()
	amb := novoAmbienteTeste(t)

	req := httptest.NewRequest(http.MethodGet, "/profissionais/1/estatisticas", nil)
	req = reqComChiParams(req, map[string]string{"profissional_id": "1"})
	rec := httptest.NewRecorder()

	amb.handlerAvaliacao.BuscarEstatisticas(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("got %d; want %d", rec.Code, http.StatusOK)
	}
}

func TestCriarAvaliacao_sem_autenticacao_retorna_401(t *testing.T) {
	t.Parallel()
	amb := novoAmbienteTeste(t)

	body := dto.RequisicaoCriarAvaliacao{LimpezaID: 1, Nota: 5, Comentario: "Ótimo"}
	req := reqComToken(http.MethodPost, "/avaliacoes", body, "")
	rec := httptest.NewRecorder()

	amb.handlerAvaliacao.CriarAvaliacao(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("got %d; want %d", rec.Code, http.StatusUnauthorized)
	}
}

// --- Testes de Solicitação ---

func TestCriarSolicitacao_sem_autenticacao_retorna_401(t *testing.T) {
	t.Parallel()
	amb := novoAmbienteTeste(t)

	body := dto.RequisicaoCriarSolicitacao{LimpezaID: 1, DataAgendada: time.Now().Add(48 * time.Hour)}
	req := reqComToken(http.MethodPost, "/solicitacoes", body, "")
	rec := httptest.NewRecorder()

	amb.handlerSolicitacao.CriarSolicitacao(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("got %d; want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestListarMinhasSolicitacoes_cliente_ve_suas_solicitacoes(t *testing.T) {
	t.Parallel()
	amb := novoAmbienteTeste(t)
	usuario, _ := amb.registrarUsuario(t, "solic@email.com", "solicuser")

	req := httptest.NewRequest(http.MethodGet, "/usuarios/eu/solicitacoes", nil)
	req = reqComContextoUsuario(req, usuario.ID)
	rec := httptest.NewRecorder()

	amb.handlerSolicitacao.ListarMinhasSolicitacoes(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("got %d; want %d", rec.Code, http.StatusOK)
	}
}
