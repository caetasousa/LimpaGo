package dto_test

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"limpaGo/api/auth"
	"limpaGo/api/dto"
	"limpaGo/domain/entity"
	errosdominio "limpaGo/domain/errors"
	"limpaGo/domain/valueobject"
)

func TestMapearErroDominio_erros_de_dominio_mapeiam_para_status_correto(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		err        error
		wantStatus int
	}{
		{name: "usuario não encontrado retorna 404", err: errosdominio.ErrUsuarioNaoEncontrado, wantStatus: http.StatusNotFound},
		{name: "limpeza não encontrada retorna 404", err: errosdominio.ErrLimpezaNaoEncontrada, wantStatus: http.StatusNotFound},
		{name: "solicitação não encontrada retorna 404", err: errosdominio.ErrSolicitacaoNaoEncontrada, wantStatus: http.StatusNotFound},
		{name: "perfil não encontrado retorna 404", err: errosdominio.ErrPerfilNaoEncontrado, wantStatus: http.StatusNotFound},
		{name: "email duplicado retorna 409", err: errosdominio.ErrEmailJaUtilizado, wantStatus: http.StatusConflict},
		{name: "nome usuario duplicado retorna 409", err: errosdominio.ErrNomeUsuarioJaUtilizado, wantStatus: http.StatusConflict},
		{name: "solicitação duplicada retorna 409", err: errosdominio.ErrSolicitacaoDuplicada, wantStatus: http.StatusConflict},
		{name: "não é faxineiro da limpeza retorna 403", err: errosdominio.ErrNaoEFaxineiroDaLimpeza, wantStatus: http.StatusForbidden},
		{name: "faxineiro não pode solicitar próprio retorna 422", err: errosdominio.ErrFaxineiroNaoPodeSolicitarProprio, wantStatus: http.StatusUnprocessableEntity},
		{name: "horário indisponível retorna 422", err: errosdominio.ErrHorarioIndisponivel, wantStatus: http.StatusUnprocessableEntity},
		{name: "credenciais inválidas retorna 401", err: auth.ErrCredenciaisInvalidas, wantStatus: http.StatusUnauthorized},
		{name: "token inválido retorna 401", err: auth.ErrTokenInvalido, wantStatus: http.StatusUnauthorized},
		{name: "usuario inativo retorna 403", err: auth.ErrUsuarioInativo, wantStatus: http.StatusForbidden},
		{name: "senha fraca retorna 422", err: auth.ErrSenhaFraca, wantStatus: http.StatusUnprocessableEntity},
		{name: "erro desconhecido retorna 500", err: errors.New("erro genérico"), wantStatus: http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			status, resp := dto.MapearErroDominio(tt.err)
			if status != tt.wantStatus {
				t.Errorf("got status %d; want %d", status, tt.wantStatus)
			}
			if resp.Codigo != tt.wantStatus {
				t.Errorf("got codigo %d; want %d", resp.Codigo, tt.wantStatus)
			}
		})
	}
}

func TestMapearErroDominio_erro_validacao_retorna_422_com_detalhes(t *testing.T) {
	t.Parallel()
	err := &entity.ErroValidacao{Campo: "nome", Mensagem: "obrigatório"}
	status, resp := dto.MapearErroDominio(err)

	if status != http.StatusUnprocessableEntity {
		t.Errorf("got status %d; want %d", status, http.StatusUnprocessableEntity)
	}
	if resp.Detalhes["nome"] != "obrigatório" {
		t.Errorf("got detalhes[nome] %q; want %q", resp.Detalhes["nome"], "obrigatório")
	}
}

func TestDeUsuario_converte_entidade_para_dto(t *testing.T) {
	t.Parallel()
	u := &entity.Usuario{ID: 1, Email: "a@b.com", NomeUsuario: "ab", Ativo: true}
	resp := dto.DeUsuario(u)

	if resp.ID != 1 || resp.Email != "a@b.com" || resp.NomeUsuario != "ab" || !resp.Ativo {
		t.Errorf("got %+v; want id=1 email=a@b.com nome=ab ativo=true", resp)
	}
}

func TestDeLimpeza_calcula_preco_total(t *testing.T) {
	t.Parallel()
	l := &entity.Limpeza{
		ID: 1, FaxineiroID: 2, Nome: "Residencial",
		ValorHora: 50.0, DuracaoEstimada: 3.0,
		TipoLimpeza: valueobject.TipoLimpezaPadrao,
	}
	resp := dto.DeLimpeza(l)

	if resp.PrecoTotal != 150.0 {
		t.Errorf("got preco_total %.2f; want 150.00", resp.PrecoTotal)
	}
}

func TestDeSolicitacao_converte_todos_os_campos(t *testing.T) {
	t.Parallel()
	data := time.Date(2026, 4, 1, 10, 0, 0, 0, time.UTC)
	s := &entity.Solicitacao{
		ID: 1, ClienteID: 2, LimpezaID: 3,
		Status:   valueobject.StatusSolicitacaoPendente,
		DataAgendada: data,
		PrecoTotal:   100.0,
		Endereco: valueobject.Endereco{Cidade: "SP"},
	}
	resp := dto.DeSolicitacao(s)

	if resp.Status != "pendente" {
		t.Errorf("got status %q; want %q", resp.Status, "pendente")
	}
	if resp.Endereco.Cidade != "SP" {
		t.Errorf("got cidade %q; want %q", resp.Endereco.Cidade, "SP")
	}
}

func TestEnderecoDTO_converte_ida_e_volta(t *testing.T) {
	t.Parallel()
	original := valueobject.Endereco{
		Rua: "Rua A", Complemento: "Apt 1", Bairro: "Centro",
		Cidade: "SP", Estado: "SP", CEP: "01000-000",
	}
	dtoEnd := dto.DeEndereco(original)
	volta := dtoEnd.ParaEndereco()

	if volta != original {
		t.Errorf("got %+v; want %+v", volta, original)
	}
}

func TestDeParTokens_converte_tokens(t *testing.T) {
	t.Parallel()
	p := &auth.ParTokens{
		TokenAcesso:    "acesso-123",
		TokenRenovacao: "renovacao-456",
		TipoToken:      "Bearer",
		ExpiraEm:       1234567890,
	}
	resp := dto.DeParTokens(p)

	if resp.TokenAcesso != "acesso-123" || resp.TokenRenovacao != "renovacao-456" {
		t.Errorf("got %+v; want acesso-123 / renovacao-456", resp)
	}
}

func TestDeAgregadoAvaliacao_converte_estatisticas(t *testing.T) {
	t.Parallel()
	a := &entity.AgregadoAvaliacao{FaxineiroID: 1, MediaNota: 4.5, TotalAvaliacoes: 10}
	resp := dto.DeAgregadoAvaliacao(a)

	if resp.MediaNota != 4.5 || resp.TotalAvaliacoes != 10 {
		t.Errorf("got media=%.1f total=%d; want 4.5 / 10", resp.MediaNota, resp.TotalAvaliacoes)
	}
}

func TestDePaginaFeed_converte_pagina_com_itens(t *testing.T) {
	t.Parallel()
	limpeza := &entity.Limpeza{
		ID: 1, FaxineiroID: 2, Nome: "Limpeza",
		ValorHora: 30, DuracaoEstimada: 2,
		TipoLimpeza: valueobject.TipoLimpezaPadrao,
	}
	pagina := &entity.PaginaFeed{
		Itens: []*entity.ItemFeed{
			{Limpeza: limpeza, TipoEvento: valueobject.TipoEventoFeedCriacao, DataEvento: time.Now(), NumeroLinha: 1},
		},
		TotalItens:    1,
		Pagina:        1,
		TamanhoPagina: 20,
	}
	resp := dto.DePaginaFeed(pagina)

	if len(resp.Itens) != 1 {
		t.Fatalf("got %d itens; want 1", len(resp.Itens))
	}
	if resp.Itens[0].Limpeza.Nome != "Limpeza" {
		t.Errorf("got nome %q; want %q", resp.Itens[0].Limpeza.Nome, "Limpeza")
	}
}

func TestNovaRespostaErro_cria_resposta_simples(t *testing.T) {
	t.Parallel()
	resp := dto.NovaRespostaErro(400, "erro")
	if resp.Codigo != 400 || resp.Mensagem != "erro" {
		t.Errorf("got %+v; want codigo=400 mensagem=erro", resp)
	}
}

func TestDeLimpezaLista_converte_lista_vazia(t *testing.T) {
	t.Parallel()
	resp := dto.DeLimpezaLista(nil)
	if len(resp) != 0 {
		t.Errorf("got %d items; want 0", len(resp))
	}
}
