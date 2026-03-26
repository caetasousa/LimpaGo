package entity

import "fmt"

// ErroValidacao é um erro de validação de campo no nível do domínio.
type ErroValidacao struct {
	Campo    string
	Mensagem string
}

func (e *ErroValidacao) Error() string {
	return fmt.Sprintf("erro de validação no campo '%s': %s", e.Campo, e.Mensagem)
}
