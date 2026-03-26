package valueobject

import "errors"

type TipoLimpeza string

const (
	TipoLimpezaPadrao      TipoLimpeza = "limpeza_padrao"
	TipoLimpezaPesada      TipoLimpeza = "limpeza_pesada"
	TipoLimpezaExpress     TipoLimpeza = "limpeza_express"
	TipoLimpezaPreMudanca  TipoLimpeza = "limpeza_pre_mudanca"
	TipoLimpezaPosObra     TipoLimpeza = "limpeza_pos_obra"
	TipoLimpezaComercial   TipoLimpeza = "limpeza_comercial"
	TipoLimpezaPassadoria  TipoLimpeza = "passadoria"
)

func (tl TipoLimpeza) Validar() error {
	switch tl {
	case TipoLimpezaPadrao, TipoLimpezaPesada, TipoLimpezaExpress,
		TipoLimpezaPreMudanca, TipoLimpezaPosObra, TipoLimpezaComercial,
		TipoLimpezaPassadoria:
		return nil
	default:
		return errors.New("tipo de limpeza inválido: deve ser limpeza_padrao, limpeza_pesada, limpeza_express, limpeza_pre_mudanca, limpeza_pos_obra, limpeza_comercial ou passadoria")
	}
}

// EResidencial retorna true se o tipo de limpeza é para residências.
func (tl TipoLimpeza) EResidencial() bool {
	switch tl {
	case TipoLimpezaPadrao, TipoLimpezaPesada, TipoLimpezaExpress,
		TipoLimpezaPreMudanca, TipoLimpezaPosObra, TipoLimpezaPassadoria:
		return true
	default:
		return false
	}
}
