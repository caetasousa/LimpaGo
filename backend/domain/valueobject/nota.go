package valueobject

import "errors"

// Nota representa uma pontuação inteira de 0 a 5 usada em avaliações.
type Nota int

func NovaNota(v int) (Nota, error) {
	if v < 0 || v > 5 {
		return 0, errors.New("a nota deve estar entre 0 e 5")
	}
	return Nota(v), nil
}
