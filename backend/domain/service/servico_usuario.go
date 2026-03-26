package service

import (
	"context"

	"limpaGo/domain/entity"
	errosdominio "limpaGo/domain/errors"
	"limpaGo/domain/repository"
	"limpaGo/domain/valueobject"
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

// --- Perfil Profissional ---

// CriarPerfilProfissional cria o perfil profissional para o usuário atuar como profissional.
func (s *ServicoUsuario) CriarPerfilProfissional(ctx context.Context, usuarioID int) (*entity.PerfilProfissional, error) {
	existenteFax, err := s.perfis.BuscarPerfilProfissional(ctx, usuarioID)
	if err != nil {
		return nil, err
	}
	if existenteFax != nil {
		return nil, errosdominio.ErrPerfilProfissionalJaExiste
	}

	perfil := entity.NovoPerfilProfissional(usuarioID)
	if err := s.perfis.SalvarPerfilProfissional(ctx, perfil); err != nil {
		return nil, err
	}
	return perfil, nil
}

// BuscarPerfilProfissional retorna o perfil profissional do profissional.
func (s *ServicoUsuario) BuscarPerfilProfissional(ctx context.Context, usuarioID int) (*entity.PerfilProfissional, error) {
	return s.perfis.BuscarPerfilProfissional(ctx, usuarioID)
}

// AtualizarPerfilProfissional atualiza os dados profissionais do profissional.
func (s *ServicoUsuario) AtualizarPerfilProfissional(ctx context.Context, usuarioID int, descricao string, anosExperiencia int, especialidades, cidadesAtendidas []string) (*entity.PerfilProfissional, error) {
	perfil, err := s.perfis.BuscarPerfilProfissional(ctx, usuarioID)
	if err != nil {
		return nil, err
	}
	if perfil == nil {
		return nil, errosdominio.ErrPerfilProfissionalNaoEncontrado
	}

	perfil.Descricao = descricao
	perfil.AnosExperiencia = anosExperiencia
	perfil.Especialidades = especialidades
	perfil.CidadesAtendidas = cidadesAtendidas

	if err := s.perfis.AtualizarPerfilProfissional(ctx, perfil); err != nil {
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
