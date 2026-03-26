package server_test

import (
	"net/http"
	"testing"
	"time"

	"limpaGo/api/server"
)

func TestNovo_servidor_configurado_com_timeouts_corretos(t *testing.T) {
	t.Parallel()
	srv := server.Novo(":8080", http.DefaultServeMux)

	if srv.Addr != ":8080" {
		t.Errorf("got addr %q; want %q", srv.Addr, ":8080")
	}
	if srv.ReadTimeout != 15*time.Second {
		t.Errorf("got ReadTimeout %v; want %v", srv.ReadTimeout, 15*time.Second)
	}
	if srv.WriteTimeout != 15*time.Second {
		t.Errorf("got WriteTimeout %v; want %v", srv.WriteTimeout, 15*time.Second)
	}
	if srv.IdleTimeout != 60*time.Second {
		t.Errorf("got IdleTimeout %v; want %v", srv.IdleTimeout, 60*time.Second)
	}
}

func TestNovo_servidor_usa_handler_fornecido(t *testing.T) {
	t.Parallel()
	mux := http.NewServeMux()
	srv := server.Novo(":3000", mux)

	if srv.Handler != mux {
		t.Error("expected server handler to be the provided mux")
	}
}
