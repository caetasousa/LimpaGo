package entity

import (
	"testing"

	"limpaGo/domain/valueobject"
)

func TestNovaAvaliacao(t *testing.T) {
	t.Parallel()

	nota, _ := valueobject.NovaNota(5)
	a := NovaAvaliacao(10, 1, 2, nota, "Excelente serviço")

	if a.LimpezaID != 10 {
		t.Errorf("LimpezaID = %d; want 10", a.LimpezaID)
	}
	if a.ProfissionalID != 1 {
		t.Errorf("ProfissionalID = %d; want 1", a.ProfissionalID)
	}
	if a.ClienteID != 2 {
		t.Errorf("ClienteID = %d; want 2", a.ClienteID)
	}
	if int(a.Nota) != 5 {
		t.Errorf("Nota = %d; want 5", a.Nota)
	}
	if a.Comentario != "Excelente serviço" {
		t.Errorf("Comentario = %q; want %q", a.Comentario, "Excelente serviço")
	}
}
