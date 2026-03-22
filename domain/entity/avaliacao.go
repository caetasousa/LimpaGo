package entity

import (
	"time"

	"phresh-go/domain/valueobject"
)

type Avaliacao struct {
	ID          int
	LimpezaID   int
	FaxineiroID int
	ClienteID   int
	Nota        valueobject.Nota
	Comentario  string
	CriadoEm    time.Time
}

// AgregadoAvaliacao contém estatísticas de reputação de um faxineiro.
type AgregadoAvaliacao struct {
	FaxineiroID     int
	MediaNota       float64
	TotalAvaliacoes int
}

func NovaAvaliacao(limpezaID, faxineiroID, clienteID int, nota valueobject.Nota, comentario string) *Avaliacao {
	return &Avaliacao{
		LimpezaID:   limpezaID,
		FaxineiroID: faxineiroID,
		ClienteID:   clienteID,
		Nota:        nota,
		Comentario:  comentario,
	}
}
