package valueobject

import "errors"

type TipoEventoFeed string

const (
	TipoEventoFeedCriacao    TipoEventoFeed = "criacao"
	TipoEventoFeedAtualizacao TipoEventoFeed = "atualizacao"
)

func (e TipoEventoFeed) Validar() error {
	switch e {
	case TipoEventoFeedCriacao, TipoEventoFeedAtualizacao:
		return nil
	default:
		return errors.New("tipo de evento do feed inválido")
	}
}
