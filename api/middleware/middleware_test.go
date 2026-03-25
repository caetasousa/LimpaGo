package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"limpaGo/api/auth"
	"limpaGo/api/middleware"
	"limpaGo/domain/entity"
)

func configTeste() auth.ConfiguracaoJWT {
	return auth.ConfiguracaoJWT{
		SegredoAcesso:    []byte("segredo-acesso-teste"),
		SegredoRenovacao: []byte("segredo-renovacao-teste"),
		DuracaoAcesso:    15 * 60 * 1e9, // 15 min
		DuracaoRenovacao: 24 * 60 * 60 * 1e9,
		Emissor:          "teste",
	}
}

func TestAutenticacaoJWT_requisicao_sem_token_retorna_401(t *testing.T) {
	t.Parallel()
	svcToken := auth.NovoServicoToken(configTeste())
	mw := middleware.AutenticacaoJWT(svcToken)
	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("got %d; want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestAutenticacaoJWT_token_invalido_retorna_401(t *testing.T) {
	t.Parallel()
	svcToken := auth.NovoServicoToken(configTeste())
	mw := middleware.AutenticacaoJWT(svcToken)
	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer token-invalido-qualquer")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("got %d; want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestAutenticacaoJWT_token_valido_injeta_usuario_no_contexto(t *testing.T) {
	t.Parallel()
	svcToken := auth.NovoServicoToken(configTeste())

	usuario := &entity.Usuario{ID: 42, Email: "teste@email.com", NomeUsuario: "teste", Ativo: true}
	token, err := svcToken.GerarTokenAcesso(usuario)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var capturadoID int
	var capturadoOK bool
	mw := middleware.AutenticacaoJWT(svcToken)
	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturadoID, capturadoOK = middleware.ObterUsuarioID(r.Context())
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("got status %d; want %d", rec.Code, http.StatusOK)
	}
	if !capturadoOK {
		t.Fatal("expected usuario_id in context; got nothing")
	}
	if capturadoID != 42 {
		t.Errorf("got usuario_id %d; want 42", capturadoID)
	}
}

func TestAutenticacaoJWT_header_sem_bearer_retorna_401(t *testing.T) {
	t.Parallel()
	svcToken := auth.NovoServicoToken(configTeste())
	mw := middleware.AutenticacaoJWT(svcToken)
	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Basic abc123")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

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
	handler := middleware.Recuperacao(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("erro inesperado no teste")
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("got %d; want %d", rec.Code, http.StatusInternalServerError)
	}
}

func TestRecuperacao_handler_normal_passa_sem_erro(t *testing.T) {
	t.Parallel()
	handler := middleware.Recuperacao(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("got %d; want %d", rec.Code, http.StatusOK)
	}
}

func TestLogger_registra_status_da_resposta(t *testing.T) {
	t.Parallel()
	handler := middleware.Logger(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}))

	req := httptest.NewRequest(http.MethodPost, "/teste", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

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
