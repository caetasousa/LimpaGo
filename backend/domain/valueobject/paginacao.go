package valueobject

// Paginacao encapsula parâmetros de paginação com validação embutida.
type Paginacao struct {
	Pagina        int
	TamanhoPagina int
}

// NovaPaginacao cria uma Paginacao com valores válidos.
// Pagina mínima: 1. TamanhoPagina: entre 1 e 100, padrão 20.
func NovaPaginacao(pagina, tamanhoPagina int) Paginacao {
	if pagina < 1 {
		pagina = 1
	}
	if tamanhoPagina < 1 || tamanhoPagina > 100 {
		tamanhoPagina = 20
	}
	return Paginacao{Pagina: pagina, TamanhoPagina: tamanhoPagina}
}
