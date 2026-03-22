package service

import (
	"context"

	"phresh-go/domain/entity"
	errosdominio "phresh-go/domain/errors"
	"phresh-go/domain/repository"
	"phresh-go/domain/valueobject"
)

type ServicoUsuario struct {
	usuarios repository.RepositorioUsuario
	perfis   repository.RepositorioPerfil
}

func NovoServicoUsuario(usuarios repository.RepositorioUsuario, perfis repository.RepositorioPerfil) *ServicoUsuario {
	return &ServicoUsuario{usuarios: usuarios, perfis: perfis}
}

// Registrar cria um novo usuário e cria automaticamente seu perfil base.
func (s *ServicoUsuario) Registrar(ctx context.Context, email, nomeUsuario string) (*entity.Usuario, error) {
	existente, err := s.usuarios.BuscarPorEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if existente != nil {
		return nil, errosdominio.ErrEmailJaUtilizado
	}

	existente, err = s.usuarios.BuscarPorNomeUsuario(ctx, nomeUsuario)
	if err != nil {
		return nil, err
	}
	if existente != nil {
		return nil, errosdominio.ErrNomeUsuarioJaUtilizado
	}

	usuario, err := entity.NovoUsuario(email, nomeUsuario)
	if err != nil {
		return nil, err
	}

	if err := s.usuarios.Salvar(ctx, usuario); err != nil {
		return nil, err
	}

	// Cria perfil base automaticamente no registro
	perfil := entity.NovoPerfil(usuario.ID, usuario.Email, usuario.NomeUsuario)
	if err := s.perfis.Salvar(ctx, perfil); err != nil {
		return nil, err
	}
	usuario.Perfil = perfil

	return usuario, nil
}

// --- Perfil Base ---

// BuscarPerfil retorna o perfil base do usuário.
func (s *ServicoUsuario) BuscarPerfil(ctx context.Context, usuarioID int) (*entity.Perfil, error) {
	return s.perfis.BuscarPorUsuarioID(ctx, usuarioID)
}

// AtualizarPerfil atualiza os dados pessoais do perfil base.
func (s *ServicoUsuario) AtualizarPerfil(ctx context.Context, usuarioID int, nomeCompleto, telefone, imagem string) (*entity.Perfil, error) {
	perfil, err := s.perfis.BuscarPorUsuarioID(ctx, usuarioID)
	if err != nil {
		return nil, err
	}

	perfil.NomeCompleto = nomeCompleto
	perfil.Telefone = telefone
	perfil.Imagem = imagem

	if err := s.perfis.Atualizar(ctx, perfil); err != nil {
		return nil, err
	}
	return perfil, nil
}

// --- Perfil Faxineiro ---

// CriarPerfilFaxineiro cria o perfil profissional para o usuário atuar como faxineiro.
func (s *ServicoUsuario) CriarPerfilFaxineiro(ctx context.Context, usuarioID int) (*entity.PerfilFaxineiro, error) {
	existenteFax, err := s.perfis.BuscarPerfilFaxineiro(ctx, usuarioID)
	if err != nil {
		return nil, err
	}
	if existenteFax != nil {
		return nil, errosdominio.ErrPerfilFaxineiroJaExiste
	}

	perfil := entity.NovoPerfilFaxineiro(usuarioID)
	if err := s.perfis.SalvarPerfilFaxineiro(ctx, perfil); err != nil {
		return nil, err
	}
	return perfil, nil
}

// BuscarPerfilFaxineiro retorna o perfil profissional do faxineiro.
func (s *ServicoUsuario) BuscarPerfilFaxineiro(ctx context.Context, usuarioID int) (*entity.PerfilFaxineiro, error) {
	return s.perfis.BuscarPerfilFaxineiro(ctx, usuarioID)
}

// AtualizarPerfilFaxineiro atualiza os dados profissionais do faxineiro.
func (s *ServicoUsuario) AtualizarPerfilFaxineiro(ctx context.Context, usuarioID int, descricao string, anosExperiencia int, especialidades, cidadesAtendidas []string) (*entity.PerfilFaxineiro, error) {
	perfil, err := s.perfis.BuscarPerfilFaxineiro(ctx, usuarioID)
	if err != nil {
		return nil, err
	}
	if perfil == nil {
		return nil, errosdominio.ErrPerfilFaxineiroNaoEncontrado
	}

	perfil.Descricao = descricao
	perfil.AnosExperiencia = anosExperiencia
	perfil.Especialidades = especialidades
	perfil.CidadesAtendidas = cidadesAtendidas

	if err := s.perfis.AtualizarPerfilFaxineiro(ctx, perfil); err != nil {
		return nil, err
	}
	return perfil, nil
}

// --- Perfil Cliente ---

// CriarPerfilCliente cria o perfil de cliente para o usuário contratar serviços.
func (s *ServicoUsuario) CriarPerfilCliente(ctx context.Context, usuarioID int) (*entity.PerfilCliente, error) {
	existenteCli, err := s.perfis.BuscarPerfilCliente(ctx, usuarioID)
	if err != nil {
		return nil, err
	}
	if existenteCli != nil {
		return nil, errosdominio.ErrPerfilClienteJaExiste
	}

	perfil := entity.NovoPerfilCliente(usuarioID)
	if err := s.perfis.SalvarPerfilCliente(ctx, perfil); err != nil {
		return nil, err
	}
	return perfil, nil
}

// BuscarPerfilCliente retorna o perfil de cliente do usuário.
func (s *ServicoUsuario) BuscarPerfilCliente(ctx context.Context, usuarioID int) (*entity.PerfilCliente, error) {
	return s.perfis.BuscarPerfilCliente(ctx, usuarioID)
}

// AtualizarPerfilCliente atualiza os dados do imóvel e preferências do cliente.
func (s *ServicoUsuario) AtualizarPerfilCliente(ctx context.Context, usuarioID int, endereco valueobject.Endereco, tipoImovel valueobject.TipoImovel, quartos, banheiros int, tamanhoImovelM2 float64, observacoes string) (*entity.PerfilCliente, error) {
	perfil, err := s.perfis.BuscarPerfilCliente(ctx, usuarioID)
	if err != nil {
		return nil, err
	}
	if perfil == nil {
		return nil, errosdominio.ErrPerfilClienteNaoEncontrado
	}

	perfil.Endereco = endereco
	perfil.TipoImovel = tipoImovel
	perfil.Quartos = quartos
	perfil.Banheiros = banheiros
	perfil.TamanhoImovelM2 = tamanhoImovelM2
	perfil.Observacoes = observacoes

	if err := s.perfis.AtualizarPerfilCliente(ctx, perfil); err != nil {
		return nil, err
	}
	return perfil, nil
}
