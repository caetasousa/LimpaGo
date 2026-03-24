// @title           LimpaGo API
// @version         1.0
// @description     API REST para plataforma de serviços de limpeza.
// @host            localhost:8080
// @BasePath        /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Use o formato "Bearer {token}"
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "limpaGo/docs"

	"limpaGo/api/auth"
	"limpaGo/api/handler"
	"limpaGo/api/router"
	"limpaGo/api/server"
	"limpaGo/domain/service"
	"limpaGo/domain/testutil"
)

func main() {
	// Repositórios em memória (substituir por implementações de banco de dados)
	repoUsuario := testutil.NovoRepositorioUsuarioMock()
	repoPerfil := testutil.NovoRepositorioPerfilMock()
	repoLimpeza := testutil.NovoRepositorioLimpezaMock()
	repoAgenda := testutil.NovoRepositorioAgendaMock()
	repoSolicitacao := testutil.NovoRepositorioSolicitacaoMock()
	repoAvaliacao := testutil.NovoRepositorioAvaliacaoMock()
	repoFeed := testutil.NovoRepositorioFeedMock()
	repoCredencial := auth.NovoRepositorioCredencialMock()

	// Serviços de domínio
	svcAgenda := service.NovoServicoAgenda(repoAgenda)
	svcUsuario := service.NovoServicoUsuario(repoUsuario, repoPerfil)
	svcLimpeza := service.NovoServicoLimpeza(repoLimpeza)
	svcSolicitacao := service.NovoServicoSolicitacao(repoSolicitacao, repoLimpeza, svcAgenda)
	svcAvaliacao := service.NovoServicoAvaliacao(repoAvaliacao, repoSolicitacao, repoLimpeza)
	svcFeed := service.NovoServicoFeed(repoFeed)

	// Serviços de autenticação
	configJWT := auth.ConfiguracaoPadrao()
	svcToken := auth.NovoServicoToken(configJWT)
	svcAuth := auth.NovoServicoAutenticacao(repoUsuario, repoCredencial, svcUsuario, svcToken)

	// Handlers HTTP
	deps := router.Dependencias{
		Autenticacao: handler.NovoHandlerAutenticacao(svcAuth),
		ServicoToken: svcToken,
		Usuario:      handler.NovoHandlerUsuario(svcUsuario),
		Limpeza:      handler.NovoHandlerLimpeza(svcLimpeza),
		Solicitacao:  handler.NovoHandlerSolicitacao(svcSolicitacao),
		Agenda:       handler.NovoHandlerAgenda(svcAgenda),
		Avaliacao:    handler.NovoHandlerAvaliacao(svcAvaliacao),
		Feed:         handler.NovoHandlerFeed(svcFeed),
	}

	addr := ":8080"
	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}

	srv := server.Novo(addr, router.Novo(deps))

	// Iniciar servidor em goroutine
	go func() {
		log.Printf("servidor iniciado em http://localhost%s", addr)
		log.Printf("swagger ui disponível em http://localhost%s/swagger/index.html", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("erro ao iniciar servidor: %v", err)
		}
	}()

	// Aguardar sinal de encerramento
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("encerrando servidor...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("erro ao encerrar servidor: %v", err)
	}
	log.Println("servidor encerrado")
}
