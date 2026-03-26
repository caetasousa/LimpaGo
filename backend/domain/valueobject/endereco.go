package valueobject

// Endereco representa um endereço completo.
type Endereco struct {
	Rua         string
	Complemento string
	Bairro      string
	Cidade      string
	Estado      string
	CEP         string
}

// EstaPreenchido retorna true se pelo menos rua e cidade estão preenchidos.
func (e Endereco) EstaPreenchido() bool {
	return e.Rua != "" && e.Cidade != ""
}
