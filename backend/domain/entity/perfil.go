package entity

import (
	"time"

	"limpaGo/domain/valueobject"
)

// Perfil contém os dados pessoais compartilhados de qualquer usuário.
// É criado automaticamente quando um Usuario se registra.
type Perfil struct {
	UsuarioID    int
	NomeCompleto string
	Telefone     string
	Imagem       string // URL da foto de perfil
	// Desnormalizados para conveniência
	Email        string
	NomeUsuario  string
	CriadoEm     time.Time
	AtualizadoEm time.Time
}

func NovoPerfil(usuarioID int, email, nomeUsuario string) *Perfil {
	return &Perfil{
		UsuarioID:   usuarioID,
		Email:       email,
		NomeUsuario: nomeUsuario,
	}
}

// PerfilProfissional contém dados profissionais do profissional.
// Criado quando o usuário decide oferecer serviços de limpeza.
type PerfilProfissional struct {
	UsuarioID        int
	Descricao        string   // apresentação profissional / bio de trabalho
	AnosExperiencia  int
	Especialidades   []string // ex: ["limpeza_padrao", "limpeza_pesada"]
	CidadesAtendidas []string // ex: ["São Paulo", "Guarulhos"]
	// Documentação e verificação
	DocumentoRG   string // número do RG
	DocumentoCPF  string // número do CPF
	FotoDocumento string // URL da foto do documento para verificação
	Verificado    bool   // se passou pelo processo de verificação da plataforma
	CriadoEm      time.Time
	AtualizadoEm  time.Time
}

func NovoPerfilProfissional(usuarioID int) *PerfilProfissional {
	return &PerfilProfissional{
		UsuarioID: usuarioID,
	}
}

// PerfilCliente contém dados específicos de quem contrata serviços.
// Criado quando o usuário faz sua primeira solicitação ou manualmente.
type PerfilCliente struct {
	UsuarioID       int
	Endereco        valueobject.Endereco
	TipoImovel      valueobject.TipoImovel
	Quartos         int     // número de quartos (usado para estimar duração)
	Banheiros       int     // número de banheiros (usado para estimar duração)
	TamanhoImovelM2 float64 // tamanho do imóvel em metros quadrados
	Observacoes     string  // ex: "tem animais de estimação", "portaria 24h"
	// Profissional preferido
	ProfissionalPreferidoID *int // ID do profissional preferido (opcional)
	CriadoEm             time.Time
	AtualizadoEm         time.Time
}

func NovoPerfilCliente(usuarioID int) *PerfilCliente {
	return &PerfilCliente{
		UsuarioID: usuarioID,
	}
}
