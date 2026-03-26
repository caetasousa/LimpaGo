package auth

import "unicode"

// ValidarForcaSenha verifica se a senha atende os requisitos mínimos:
// - mínimo 8 caracteres
// - pelo menos uma letra maiúscula
// - pelo menos um dígito
func ValidarForcaSenha(senha string) error {
	if len(senha) < 8 {
		return ErrSenhaFraca
	}

	temMaiuscula := false
	temDigito := false

	for _, c := range senha {
		switch {
		case unicode.IsUpper(c):
			temMaiuscula = true
		case unicode.IsDigit(c):
			temDigito = true
		}
	}

	if !temMaiuscula || !temDigito {
		return ErrSenhaFraca
	}

	return nil
}
