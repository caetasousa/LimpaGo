// Package router configura todas as rotas da API com Chi.
//
// @title           LimpaGo API
// @version         1.0
// @description     API REST para plataforma de serviços de limpeza.
// @host            localhost:8080
// @BasePath        /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Use o formato "Bearer {token}"
package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	"limpaGo/api/auth"
	"limpaGo/api/handler"
	"limpaGo/api/middleware"
)

// Dependencias agrupa todos os handlers e serviços necessários para construir o router.
type Dependencias struct {
	Autenticacao *handler.HandlerAutenticacao
	ServicoToken *auth.ServicoToken
	Usuario      *handler.HandlerUsuario
	Limpeza      *handler.HandlerLimpeza
	Solicitacao  *handler.HandlerSolicitacao
	Agenda       *handler.HandlerAgenda
	Avaliacao    *handler.HandlerAvaliacao
	Feed         *handler.HandlerFeed
}

// Novo constrói e retorna o router Chi com todas as rotas registradas.
func Novo(d Dependencias) http.Handler {
	r := chi.NewRouter()

	// Middlewares globais
	r.Use(chimiddleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(cors.Handler(middleware.OpcoesCORS()))

	// Swagger UI
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	// API v1
	r.Route("/api/v1", func(r chi.Router) {
		// Autenticação (público)
		r.Post("/auth/registrar", d.Autenticacao.Registrar)
		r.Post("/auth/login", d.Autenticacao.Login)
		r.Post("/auth/renovar", d.Autenticacao.RenovarToken)

		// Usuários (autenticado)
		r.Group(func(r chi.Router) {
			r.Use(middleware.AutenticacaoJWT(d.ServicoToken))

			r.Get("/usuarios/eu/perfil", d.Usuario.BuscarMeuPerfil)
			r.Put("/usuarios/eu/perfil", d.Usuario.AtualizarMeuPerfil)
			r.Post("/usuarios/eu/perfil-faxineiro", d.Usuario.CriarPerfilFaxineiro)
			r.Get("/usuarios/eu/perfil-faxineiro", d.Usuario.BuscarPerfilFaxineiro)
			r.Put("/usuarios/eu/perfil-faxineiro", d.Usuario.AtualizarPerfilFaxineiro)
			r.Post("/usuarios/eu/perfil-cliente", d.Usuario.CriarPerfilCliente)
			r.Get("/usuarios/eu/perfil-cliente", d.Usuario.BuscarPerfilCliente)
			r.Put("/usuarios/eu/perfil-cliente", d.Usuario.AtualizarPerfilCliente)
			r.Get("/usuarios/eu/limpezas", d.Limpeza.ListarMinhasLimpezas)
			r.Get("/usuarios/eu/solicitacoes", d.Solicitacao.ListarMinhasSolicitacoes)
		})

		// Limpezas (público)
		r.Get("/limpezas", d.Limpeza.ListarCatalogo)
		r.Get("/limpezas/{id}", d.Limpeza.BuscarLimpeza)

		// Limpezas (autenticado)
		r.Group(func(r chi.Router) {
			r.Use(middleware.AutenticacaoJWT(d.ServicoToken))

			r.Post("/limpezas", d.Limpeza.CriarLimpeza)
			r.Put("/limpezas/{id}", d.Limpeza.AtualizarLimpeza)
			r.Delete("/limpezas/{id}", d.Limpeza.DeletarLimpeza)
			r.Get("/limpezas/{limpeza_id}/solicitacoes", d.Solicitacao.ListarPorLimpeza)
		})

		// Solicitações (autenticado)
		r.Group(func(r chi.Router) {
			r.Use(middleware.AutenticacaoJWT(d.ServicoToken))

			r.Post("/solicitacoes", d.Solicitacao.CriarSolicitacao)
			r.Post("/solicitacoes/{cliente_id}/{limpeza_id}/aceitar", d.Solicitacao.AceitarSolicitacao)
			r.Post("/solicitacoes/{cliente_id}/{limpeza_id}/rejeitar", d.Solicitacao.RejeitarSolicitacao)
			r.Post("/solicitacoes/{limpeza_id}/cancelar", d.Solicitacao.CancelarSolicitacao)
		})

		// Agenda (autenticado)
		r.Group(func(r chi.Router) {
			r.Use(middleware.AutenticacaoJWT(d.ServicoToken))

			r.Get("/agenda/disponibilidades", d.Agenda.ListarDisponibilidade)
			r.Post("/agenda/disponibilidades", d.Agenda.AdicionarDisponibilidade)
			r.Delete("/agenda/disponibilidades/{id}", d.Agenda.RemoverDisponibilidade)
			r.Get("/agenda/bloqueios", d.Agenda.ListarBloqueios)
			r.Post("/agenda/bloqueios", d.Agenda.CriarBloqueioPessoal)
			r.Delete("/agenda/bloqueios/{id}", d.Agenda.RemoverBloqueioPessoal)
		})

		// Avaliações (público)
		r.Get("/faxineiros/{faxineiro_id}/avaliacoes", d.Avaliacao.ListarAvaliacoes)
		r.Get("/faxineiros/{faxineiro_id}/estatisticas", d.Avaliacao.BuscarEstatisticas)

		// Avaliações (autenticado)
		r.Group(func(r chi.Router) {
			r.Use(middleware.AutenticacaoJWT(d.ServicoToken))

			r.Post("/avaliacoes", d.Avaliacao.CriarAvaliacao)
		})

		// Feed (público)
		r.Get("/feed", d.Feed.BuscarFeed)
	})

	return r
}
