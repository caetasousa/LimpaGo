package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"limpaGo/api/auth"
	"limpaGo/api/middleware"
)

func novoSvcToken(t *testing.T) *auth.ServicoToken {
	t.Helper()
	return auth.NovoServicoToken(auth.ConfiguracaoPadrao())
}

func TestAutenticacaoJWT_requisicao_sem_token_retorna_401(t *testing.T) {
	t.Parallel()
	mw := middleware.AutenticacaoJWT(novoSvcToken(t))
	h := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("got %d; want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestAutenticacaoJWT_token_invalido_retorna_401(t *testing.T) {
	t.Parallel()
	mw := middleware.AutenticacaoJWT(novoSvcToken(t))
	h := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer token-invalido-qualquer")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("got %d; want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestAutenticacaoJWT_header_sem_bearer_retorna_401(t *testing.T) {
	t.Parallel()
	mw := middleware.AutenticacaoJWT(novoSvcToken(t))
	h := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Basic abc123")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("got %d; want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestObterUsuarioID_contexto_sem_id_retorna_false(t *testing.T) {
	t.Parallel()
	_, ok := middleware.ObterUsuarioID(context.Background())
	if ok {
		t.Error("expected ok=false for empty context; got true")
	}
}

func TestRecuperacao_handler_com_panic_retorna_500(t *testing.T) {
	t.Parallel()
	h := middleware.Recuperacao(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("erro inesperado no teste")
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("got %d; want %d", rec.Code, http.StatusInternalServerError)
	}
}

func TestRecuperacao_handler_normal_passa_sem_erro(t *testing.T) {
	t.Parallel()
	h := middleware.Recuperacao(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("got %d; want %d", rec.Code, http.StatusOK)
	}
}

func TestLogger_registra_status_da_resposta(t *testing.T) {
	t.Parallel()
	h := middleware.Logger(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}))

	req := httptest.NewRequest(http.MethodPost, "/teste", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("got %d; want %d", rec.Code, http.StatusCreated)
	}
}

func TestOpcoesCORS_permite_metodos_esperados(t *testing.T) {
	t.Parallel()
	opts := middleware.OpcoesCORS()

	esperados := map[string]bool{
		"GET":     false,
		"POST":    false,
		"PUT":     false,
		"DELETE":  false,
		"OPTIONS": false,
	}
	for _, m := range opts.AllowedMethods {
		esperados[m] = true
	}
	for metodo, encontrado := range esperados {
		if !encontrado {
			t.Errorf("método %s não está nos AllowedMethods", metodo)
		}
	}
}
