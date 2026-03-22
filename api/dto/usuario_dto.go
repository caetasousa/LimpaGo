package dto

import (
	"phresh-go/domain/entity"
	"phresh-go/domain/valueobject"
)

// RequisicaoRegistro representa o corpo da requisição de registro de usuário.
type RequisicaoRegistro struct {
	Email       string `json:"email"`
	NomeUsuario string `json:"nome_usuario"`
}

// RespostaUsuario representa o usuário na resposta da API.
type RespostaUsuario struct {
	ID          int    `json:"id"`
	Email       string `json:"email"`
	NomeUsuario string `json:"nome_usuario"`
	Ativo       bool   `json:"ativo"`
}

// DeUsuario converte uma entidade Usuario para RespostaUsuario.
func DeUsuario(u *entity.Usuario) RespostaUsuario {
	return RespostaUsuario{
		ID:          u.ID,
		Email:       u.Email,
		NomeUsuario: u.NomeUsuario,
		Ativo:       u.Ativo,
	}
}

// RespostaPerfil representa o perfil base na resposta da API.
type RespostaPerfil struct {
	UsuarioID    int    `json:"usuario_id"`
	NomeCompleto string `json:"nome_completo"`
	Telefone     string `json:"telefone"`
	Imagem       string `json:"imagem"`
	Email        string `json:"email"`
	NomeUsuario  string `json:"nome_usuario"`
}

// DePerfil converte uma entidade Perfil para RespostaPerfil.
func DePerfil(p *entity.Perfil) RespostaPerfil {
	return RespostaPerfil{
		UsuarioID:    p.UsuarioID,
		NomeCompleto: p.NomeCompleto,
		Telefone:     p.Telefone,
		Imagem:       p.Imagem,
		Email:        p.Email,
		NomeUsuario:  p.NomeUsuario,
	}
}

// RequisicaoAtualizarPerfil representa o corpo para atualizar o perfil base.
type RequisicaoAtualizarPerfil struct {
	NomeCompleto string `json:"nome_completo"`
	Telefone     string `json:"telefone"`
	Imagem       string `json:"imagem"`
}

// RespostaPerfilFaxineiro representa o perfil do faxineiro na resposta da API.
type RespostaPerfilFaxineiro struct {
	UsuarioID        int      `json:"usuario_id"`
	Descricao        string   `json:"descricao"`
	AnosExperiencia  int      `json:"anos_experiencia"`
	Especialidades   []string `json:"especialidades"`
	CidadesAtendidas []string `json:"cidades_atendidas"`
	Verificado       bool     `json:"verificado"`
}

// DePerfilFaxineiro converte uma entidade PerfilFaxineiro para RespostaPerfilFaxineiro.
func DePerfilFaxineiro(p *entity.PerfilFaxineiro) RespostaPerfilFaxineiro {
	return RespostaPerfilFaxineiro{
		UsuarioID:        p.UsuarioID,
		Descricao:        p.Descricao,
		AnosExperiencia:  p.AnosExperiencia,
		Especialidades:   p.Especialidades,
		CidadesAtendidas: p.CidadesAtendidas,
		Verificado:       p.Verificado,
	}
}

// RequisicaoAtualizarPerfilFaxineiro representa o corpo para atualizar o perfil do faxineiro.
type RequisicaoAtualizarPerfilFaxineiro struct {
	Descricao        string   `json:"descricao"`
	AnosExperiencia  int      `json:"anos_experiencia"`
	Especialidades   []string `json:"especialidades"`
	CidadesAtendidas []string `json:"cidades_atendidas"`
}

// EnderecoDTO representa um endereço nos DTOs.
type EnderecoDTO struct {
	Rua         string `json:"rua"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Cidade      string `json:"cidade"`
	Estado      string `json:"estado"`
	CEP         string `json:"cep"`
}

// ParaEndereco converte EnderecoDTO para valueobject.Endereco.
func (e EnderecoDTO) ParaEndereco() valueobject.Endereco {
	return valueobject.Endereco{
		Rua:         e.Rua,
		Complemento: e.Complemento,
		Bairro:      e.Bairro,
		Cidade:      e.Cidade,
		Estado:      e.Estado,
		CEP:         e.CEP,
	}
}

// DeEndereco converte valueobject.Endereco para EnderecoDTO.
func DeEndereco(e valueobject.Endereco) EnderecoDTO {
	return EnderecoDTO{
		Rua:         e.Rua,
		Complemento: e.Complemento,
		Bairro:      e.Bairro,
		Cidade:      e.Cidade,
		Estado:      e.Estado,
		CEP:         e.CEP,
	}
}

// RespostaPerfilCliente representa o perfil do cliente na resposta da API.
type RespostaPerfilCliente struct {
	UsuarioID           int         `json:"usuario_id"`
	Endereco            EnderecoDTO `json:"endereco"`
	TipoImovel          string      `json:"tipo_imovel"`
	Quartos             int         `json:"quartos"`
	Banheiros           int         `json:"banheiros"`
	TamanhoImovelM2     float64     `json:"tamanho_imovel_m2"`
	Observacoes         string      `json:"observacoes"`
	FaxineiroPreferidoID *int       `json:"faxineiro_preferido_id,omitempty"`
}

// DePerfilCliente converte uma entidade PerfilCliente para RespostaPerfilCliente.
func DePerfilCliente(p *entity.PerfilCliente) RespostaPerfilCliente {
	return RespostaPerfilCliente{
		UsuarioID:            p.UsuarioID,
		Endereco:             DeEndereco(p.Endereco),
		TipoImovel:           string(p.TipoImovel),
		Quartos:              p.Quartos,
		Banheiros:            p.Banheiros,
		TamanhoImovelM2:      p.TamanhoImovelM2,
		Observacoes:          p.Observacoes,
		FaxineiroPreferidoID: p.FaxineiroPreferidoID,
	}
}

// RequisicaoAtualizarPerfilCliente representa o corpo para atualizar o perfil do cliente.
type RequisicaoAtualizarPerfilCliente struct {
	Endereco        EnderecoDTO `json:"endereco"`
	TipoImovel      string      `json:"tipo_imovel"`
	Quartos         int         `json:"quartos"`
	Banheiros       int         `json:"banheiros"`
	TamanhoImovelM2 float64     `json:"tamanho_imovel_m2"`
	Observacoes     string      `json:"observacoes"`
}
