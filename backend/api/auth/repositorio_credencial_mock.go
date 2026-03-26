package auth

import (
	"context"
	"errors"
	"sync"
	"time"
)

// RepositorioCredencialMock é uma implementação in-memory de RepositorioCredencial para desenvolvimento e testes.
type RepositorioCredencialMock struct {
	mu   sync.RWMutex
	dados map[int]*Credencial
}

// NovoRepositorioCredencialMock cria um novo mock in-memory.
func NovoRepositorioCredencialMock() *RepositorioCredencialMock {
	return &RepositorioCredencialMock{
		dados: make(map[int]*Credencial),
	}
}

func (r *RepositorioCredencialMock) BuscarPorUsuarioID(ctx context.Context, usuarioID int) (*Credencial, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	c, ok := r.dados[usuarioID]
	if !ok {
		return nil, errors.New("credencial não encontrada")
	}
	copia := *c
	return &copia, nil
}

func (r *RepositorioCredencialMock) Salvar(ctx context.Context, credencial *Credencial) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	agora := time.Now()
	if _, existe := r.dados[credencial.UsuarioID]; !existe {
		credencial.CriadoEm = agora
	}
	credencial.AtualizadoEm = agora

	copia := *credencial
	r.dados[credencial.UsuarioID] = &copia
	return nil
}
