package entity

import (
	"time"

	errosdominio "limpaGo/domain/errors"
	"limpaGo/domain/valueobject"
)

type Limpeza struct {
	ID              int
	Nome            string
	Descricao       string
	ValorHora       float64 // valor cobrado por hora para este serviço
	DuracaoEstimada float64 // duração estimada em horas
	TipoLimpeza     valueobject.TipoLimpeza
	ProfissionalID  int
	Profissional    *Usuario
	CriadoEm        time.Time
	AtualizadoEm    time.Time
}

func NovaLimpeza(profissionalID int, nome string, valorHora, duracaoEstimada float64, tipoLimpeza valueobject.TipoLimpeza) (*Limpeza, error) {
	if nome == "" {
		return nil, &ErroValidacao{Campo: "nome", Mensagem: "nome é obrigatório"}
	}
	if valorHora <= 0 {
		return nil, &ErroValidacao{Campo: "valor_hora", Mensagem: "valor por hora deve ser maior que zero"}
	}
	if duracaoEstimada <= 0 {
		return nil, &ErroValidacao{Campo: "duracao_estimada", Mensagem: "duração estimada deve ser maior que zero"}
	}
	if err := tipoLimpeza.Validar(); err != nil {
		return nil, &ErroValidacao{Campo: "tipo_limpeza", Mensagem: err.Error()}
	}

	return &Limpeza{
		ProfissionalID:  profissionalID,
		Nome:            nome,
		ValorHora:       valorHora,
		DuracaoEstimada: duracaoEstimada,
		TipoLimpeza:     tipoLimpeza,
	}, nil
}

// PrecoTotal retorna o preço total do serviço (valor/hora × duração estimada).
func (l *Limpeza) PrecoTotal() float64 {
	return l.ValorHora * l.DuracaoEstimada
}

// EPublicadoPor verifica se o profissional fornecido é quem publicou este serviço.
func (l *Limpeza) EPublicadoPor(profissionalID int) bool {
	return l.ProfissionalID == profissionalID
}

// VerificarPropriedade retorna um erro se o profissionalID não for o publicador do serviço.
func (l *Limpeza) VerificarPropriedade(profissionalID int) error {
	if !l.EPublicadoPor(profissionalID) {
		return errosdominio.ErrNaoEProfissionalDaLimpeza
	}
	return nil
}
