package auth

import "context"

// RepositorioCredencial define o contrato de persistência para credenciais.
type RepositorioCredencial interface {
	BuscarPorUsuarioID(ctx context.Context, usuarioID int) (*Credencial, error)
	Salvar(ctx context.Context, credencial *Credencial) error
}
