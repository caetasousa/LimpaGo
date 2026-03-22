// @title           Phresh-Go API
// @version         1.0
// @description     API REST para plataforma de serviços de limpeza.
// @host            localhost:8080
// @BasePath        /api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-User-ID
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "phresh-go/docs"

	"phresh-go/api/handler"
	"phresh-go/api/router"
	"phresh-go/api/server"
	"phresh-go/domain/service"
	"phresh-go/domain/testutil"
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

	// Serviços de domínio
	svcAgenda := service.NovoServicoAgenda(repoAgenda)
	svcUsuario := service.NovoServicoUsuario(repoUsuario, repoPerfil)
	svcLimpeza := service.NovoServicoLimpeza(repoLimpeza)
	svcSolicitacao := service.NovoServicoSolicitacao(repoSolicitacao, repoLimpeza, svcAgenda)
	svcAvaliacao := service.NovoServicoAvaliacao(repoAvaliacao, repoSolicitacao, repoLimpeza)
	svcFeed := service.NovoServicoFeed(repoFeed)

	// Handlers HTTP
	deps := router.Dependencias{
		Usuario:     handler.NovoHandlerUsuario(svcUsuario),
		Limpeza:     handler.NovoHandlerLimpeza(svcLimpeza),
		Solicitacao: handler.NovoHandlerSolicitacao(svcSolicitacao),
		Agenda:      handler.NovoHandlerAgenda(svcAgenda),
		Avaliacao:   handler.NovoHandlerAvaliacao(svcAvaliacao),
		Feed:        handler.NovoHandlerFeed(svcFeed),
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
