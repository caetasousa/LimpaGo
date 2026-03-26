package entity

import (
	"time"

	"limpaGo/domain/valueobject"
)

type Avaliacao struct {
	ID          int
	LimpezaID   int
	ProfissionalID int
	ClienteID   int
	Nota        valueobject.Nota
	Comentario  string
	CriadoEm    time.Time
}

// AgregadoAvaliacao contém estatísticas de reputação de um profissional.
type AgregadoAvaliacao struct {
	ProfissionalID     int
	MediaNota       float64
	TotalAvaliacoes int
}

func NovaAvaliacao(limpezaID, profissionalID, clienteID int, nota valueobject.Nota, comentario string) *Avaliacao {
	return &Avaliacao{
		LimpezaID:   limpezaID,
		ProfissionalID: profissionalID,
		ClienteID:   clienteID,
		Nota:        nota,
		Comentario:  comentario,
	}
}
