package middleware

import (
	"github.com/go-chi/cors"
)

// ConfigurarCORS retorna o middleware de CORS configurado para desenvolvimento.
// TODO: restringir AllowedOrigins para produção.
func ConfigurarCORS() func(next interface{ ServeHTTP(interface{}, interface{}) }) interface{ ServeHTTP(interface{}, interface{}) } {
	return nil // substituído pelo cors.Handler direto no router
}

// opcoesCORS retorna as opções de CORS para uso no router.
func OpcoesCORS() cors.Options {
	return cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-User-ID"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}
}
