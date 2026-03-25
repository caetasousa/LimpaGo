package router_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"limpaGo/api/auth"
	"limpaGo/api/handler"
	"limpaGo/api/router"
	"limpaGo/domain/service"
	"limpaGo/domain/testutil"
)

func criarDependencias(t *testing.T) router.Dependencias {
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

	svcTokenOIDC := auth.NovoServicoTokenOIDCMock()
	sincronizacao := auth.NovoServicoSincronizacao(repoUsuarios, svcUsuario)
	cfgZitadel := auth.CarregarConfiguracaoZitadel()
	clienteZitadel := auth.NovoClienteZitadel(cfgZitadel)
	svcAuth := auth.NovoServicoAutenticacao(clienteZitadel, sincronizacao, svcTokenOIDC)

	return router.Dependencias{
		Autenticacao:     handler.NovoHandlerAutenticacao(svcAuth, cfgZitadel),
		ServicoTokenOIDC: svcTokenOIDC,
		Sincronizacao:    sincronizacao,
		Usuario:          handler.NovoHandlerUsuario(svcUsuario),
		Limpeza:          handler.NovoHandlerLimpeza(svcLimpeza),
		Solicitacao:      handler.NovoHandlerSolicitacao(svcSolicitacao),
		Agenda:           handler.NovoHandlerAgenda(svcAgenda),
		Avaliacao:        handler.NovoHandlerAvaliacao(svcAvaliacao),
		Feed:             handler.NovoHandlerFeed(svcFeed),
	}
}

func TestRouter_rotas_publicas_respondem(t *testing.T) {
	t.Parallel()
	deps := criarDependencias(t)
	r := router.Novo(deps)

	tests := []struct {
		name   string
		method string
		path   string
		want   int
	}{
		{name: "catalogo de limpezas é acessível sem autenticação", method: http.MethodGet, path: "/api/v1/limpezas", want: http.StatusOK},
		{name: "feed é acessível sem autenticação", method: http.MethodGet, path: "/api/v1/feed", want: http.StatusOK},
		{name: "config OIDC é acessível sem autenticação", method: http.MethodGet, path: "/api/v1/auth/config", want: http.StatusOK},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)
			if rec.Code != tt.want {
				t.Errorf("got %d; want %d", rec.Code, tt.want)
			}
		})
	}
}

func TestRouter_rotas_protegidas_exigem_autenticacao(t *testing.T) {
	t.Parallel()
	deps := criarDependencias(t)
	r := router.Novo(deps)

	tests := []struct {
		name   string
		method string
		path   string
	}{
		{name: "perfil do usuario exige autenticação", method: http.MethodGet, path: "/api/v1/usuarios/eu/perfil"},
		{name: "criar limpeza exige autenticação", method: http.MethodPost, path: "/api/v1/limpezas"},
		{name: "criar solicitação exige autenticação", method: http.MethodPost, path: "/api/v1/solicitacoes"},
		{name: "agenda de disponibilidades exige autenticação", method: http.MethodGet, path: "/api/v1/agenda/disponibilidades"},
		{name: "criar avaliação exige autenticação", method: http.MethodPost, path: "/api/v1/avaliacoes"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)
			if rec.Code != http.StatusUnauthorized {
				t.Errorf("got %d; want %d", rec.Code, http.StatusUnauthorized)
			}
		})
	}
}

func TestRouter_rota_inexistente_retorna_405_ou_404(t *testing.T) {
	t.Parallel()
	deps := criarDependencias(t)
	r := router.Novo(deps)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/rota-que-nao-existe", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound && rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("got %d; want 404 or 405", rec.Code)
	}
}
