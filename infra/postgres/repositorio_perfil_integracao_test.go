//go:build integration

package postgres_test

import (
	"context"
	"errors"
	"testing"

	"limpaGo/domain/entity"
	errosdominio "limpaGo/domain/errors"
	"limpaGo/domain/valueobject"
	"limpaGo/infra/postgres"
)

func TestPerfil_GerenciarPerfilBaseDoUsuario(t *testing.T) {
	db := criarBancoTeste(t)
	t.Cleanup(func() { limparTabelas(t, db) })
	repo := postgres.NovoRepositorioPerfilPG(db)
	ctx := context.Background()

	usuarioID := inserirUsuario(t, db, "perfil@teste.com", "perfilteste")

	t.Run("usuario cria perfil com nome e telefone", func(t *testing.T) {
		p := entity.NovoPerfil(usuarioID, "perfil@teste.com", "perfilteste")
		p.NomeCompleto = "João da Silva"
		p.Telefone = "11999999999"

		if err := repo.Salvar(ctx, p); err != nil {
			t.Fatalf("Salvar() error: %v", err)
		}
		if p.CriadoEm.IsZero() {
			t.Error("CriadoEm zerado; want preenchido")
		}
	})

	t.Run("perfil criado e recuperado com dados corretos", func(t *testing.T) {
		got, err := repo.BuscarPorUsuarioID(ctx, usuarioID)
		if err != nil {
			t.Fatalf("BuscarPorUsuarioID() error: %v", err)
		}
		if got == nil {
			t.Fatal("got nil; want perfil")
		}
		if got.NomeCompleto != "João da Silva" {
			t.Errorf("NomeCompleto = %q; want %q", got.NomeCompleto, "João da Silva")
		}
	})

	t.Run("usuario atualiza nome e telefone do perfil", func(t *testing.T) {
		p := &entity.Perfil{UsuarioID: usuarioID, NomeCompleto: "João Silva Atualizado", Telefone: "11888888888"}
		if err := repo.Atualizar(ctx, p); err != nil {
			t.Fatalf("Atualizar() error: %v", err)
		}
		if p.AtualizadoEm.IsZero() {
			t.Error("AtualizadoEm zerado; want preenchido")
		}
	})

	t.Run("perfil inexistente retorna erro de nao encontrado", func(t *testing.T) {
		_, err := repo.BuscarPorUsuarioID(ctx, 999999)
		if !errors.Is(err, errosdominio.ErrPerfilNaoEncontrado) {
			t.Errorf("got %v; want %v", err, errosdominio.ErrPerfilNaoEncontrado)
		}
	})
}

func TestPerfil_FaxineiroCadastraPerfilProfissionalComEspecialidades(t *testing.T) {
	db := criarBancoTeste(t)
	t.Cleanup(func() { limparTabelas(t, db) })
	repo := postgres.NovoRepositorioPerfilPG(db)
	ctx := context.Background()

	usuarioID := inserirUsuario(t, db, "faxineiro@teste.com", "faxineiroteste")

	t.Run("faxineiro registra especialidades e cidades atendidas", func(t *testing.T) {
		p := entity.NovoPerfilFaxineiro(usuarioID)
		p.Descricao = "Faxineiro experiente"
		p.AnosExperiencia = 5
		p.Especialidades = []string{"limpeza_padrao", "limpeza_pesada"}
		p.CidadesAtendidas = []string{"São Paulo", "Guarulhos"}
		p.DocumentoCPF = "123.456.789-00"
		p.DocumentoRG = "12.345.678-9"

		if err := repo.SalvarPerfilFaxineiro(ctx, p); err != nil {
			t.Fatalf("SalvarPerfilFaxineiro() error: %v", err)
		}
		if p.CriadoEm.IsZero() {
			t.Error("CriadoEm zerado; want preenchido")
		}
	})

	t.Run("especialidades e cidades sao preservadas ao consultar perfil", func(t *testing.T) {
		got, err := repo.BuscarPerfilFaxineiro(ctx, usuarioID)
		if err != nil {
			t.Fatalf("BuscarPerfilFaxineiro() error: %v", err)
		}
		if got == nil {
			t.Fatal("got nil; want perfil faxineiro")
		}
		if len(got.Especialidades) != 2 {
			t.Errorf("Especialidades len = %d; want 2", len(got.Especialidades))
		}
		if len(got.CidadesAtendidas) != 2 {
			t.Errorf("CidadesAtendidas len = %d; want 2", len(got.CidadesAtendidas))
		}
		if got.AnosExperiencia != 5 {
			t.Errorf("AnosExperiencia = %d; want 5", got.AnosExperiencia)
		}
	})

	t.Run("faxineiro atualiza experiencia e especialidades", func(t *testing.T) {
		p := &entity.PerfilFaxineiro{
			UsuarioID:        usuarioID,
			Descricao:        "Atualizado",
			AnosExperiencia:  7,
			Especialidades:   []string{"limpeza_padrao"},
			CidadesAtendidas: []string{"Campinas"},
			Verificado:       true,
		}
		if err := repo.AtualizarPerfilFaxineiro(ctx, p); err != nil {
			t.Fatalf("AtualizarPerfilFaxineiro() error: %v", err)
		}
	})

	t.Run("perfil faxineiro inexistente retorna erro de nao encontrado", func(t *testing.T) {
		_, err := repo.BuscarPerfilFaxineiro(ctx, 999999)
		if !errors.Is(err, errosdominio.ErrPerfilFaxineiroNaoEncontrado) {
			t.Errorf("got %v; want %v", err, errosdominio.ErrPerfilFaxineiroNaoEncontrado)
		}
	})
}

func TestPerfil_ClienteCadastraEnderecoETipoDeImovel(t *testing.T) {
	db := criarBancoTeste(t)
	t.Cleanup(func() { limparTabelas(t, db) })
	repo := postgres.NovoRepositorioPerfilPG(db)
	ctx := context.Background()

	usuarioID := inserirUsuario(t, db, "cliente@teste.com", "clienteteste")

	t.Run("cliente registra endereco e detalhes do imovel", func(t *testing.T) {
		p := entity.NovoPerfilCliente(usuarioID)
		p.Endereco = valueobject.Endereco{
			Rua:    "Rua das Flores",
			Bairro: "Centro",
			Cidade: "São Paulo",
			Estado: "SP",
			CEP:    "01310-100",
		}
		p.TipoImovel = valueobject.TipoImovelApartamento
		p.Quartos = 2
		p.Banheiros = 1
		p.TamanhoImovelM2 = 60.5

		if err := repo.SalvarPerfilCliente(ctx, p); err != nil {
			t.Fatalf("SalvarPerfilCliente() error: %v", err)
		}
		if p.CriadoEm.IsZero() {
			t.Error("CriadoEm zerado; want preenchido")
		}
	})

	t.Run("endereco e tipo de imovel sao preservados ao consultar", func(t *testing.T) {
		got, err := repo.BuscarPerfilCliente(ctx, usuarioID)
		if err != nil {
			t.Fatalf("BuscarPerfilCliente() error: %v", err)
		}
		if got == nil {
			t.Fatal("got nil; want perfil cliente")
		}
		if got.Endereco.Cidade != "São Paulo" {
			t.Errorf("Cidade = %q; want %q", got.Endereco.Cidade, "São Paulo")
		}
		if got.TipoImovel != valueobject.TipoImovelApartamento {
			t.Errorf("TipoImovel = %q; want %q", got.TipoImovel, valueobject.TipoImovelApartamento)
		}
		if got.Quartos != 2 {
			t.Errorf("Quartos = %d; want 2", got.Quartos)
		}
	})

	t.Run("perfil cliente inexistente retorna erro de nao encontrado", func(t *testing.T) {
		_, err := repo.BuscarPerfilCliente(ctx, 999999)
		if !errors.Is(err, errosdominio.ErrPerfilClienteNaoEncontrado) {
			t.Errorf("got %v; want %v", err, errosdominio.ErrPerfilClienteNaoEncontrado)
		}
	})
}
