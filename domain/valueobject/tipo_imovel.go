package valueobject

import "errors"

type TipoImovel string

const (
	TipoImovelApartamento TipoImovel = "apartamento"
	TipoImovelCasa        TipoImovel = "casa"
	TipoImovelComercial   TipoImovel = "comercial"
)

func (ti TipoImovel) Validar() error {
	switch ti {
	case TipoImovelApartamento, TipoImovelCasa, TipoImovelComercial:
		return nil
	default:
		return errors.New("tipo de imóvel inválido: deve ser apartamento, casa ou comercial")
	}
}
