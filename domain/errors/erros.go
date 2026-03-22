package errors

import "errors"

var (
	// Usuario
	ErrEmailJaUtilizado       = errors.New("email já está em uso")
	ErrNomeUsuarioJaUtilizado = errors.New("nome de usuário já está em uso")
	ErrUsuarioNaoEncontrado   = errors.New("usuário não encontrado")

	// Perfil
	ErrPerfilNaoEncontrado          = errors.New("perfil não encontrado")
	ErrPerfilFaxineiroNaoEncontrado = errors.New("perfil de faxineiro não encontrado")
	ErrPerfilClienteNaoEncontrado   = errors.New("perfil de cliente não encontrado")
	ErrPerfilFaxineiroJaExiste      = errors.New("usuário já possui um perfil de faxineiro")
	ErrPerfilClienteJaExiste        = errors.New("usuário já possui um perfil de cliente")

	// Limpeza (Serviço publicado pelo faxineiro)
	ErrLimpezaNaoEncontrada   = errors.New("serviço de limpeza não encontrado")
	ErrNaoEFaxineiroDaLimpeza = errors.New("apenas o faxineiro que publicou o serviço pode realizar esta ação")

	// Solicitacao (Cliente solicita o serviço de um faxineiro)
	ErrSolicitacaoNaoEncontrada         = errors.New("solicitação não encontrada")
	ErrSolicitacaoDuplicada             = errors.New("cliente já possui uma solicitação para este serviço")
	ErrFaxineiroNaoPodeSolicitarProprio = errors.New("o faxineiro não pode solicitar o próprio serviço")
	ErrSolicitacaoNaoPodeSerAceita      = errors.New("apenas solicitações pendentes podem ser aceitas")
	ErrSolicitacaoNaoPodeSerCancelada   = errors.New("esta solicitação não pode ser cancelada no estado atual")
	ErrSolicitacaoNaoPodeSerRejeitada   = errors.New("apenas solicitações pendentes podem ser rejeitadas")
	ErrNaoEClienteSolicitante           = errors.New("apenas o cliente que fez a solicitação pode realizar esta ação")
	ErrNaoEFaxineiroDaSolicitacao       = errors.New("apenas o faxineiro do serviço pode aceitar ou rejeitar solicitações")

	// Agenda
	ErrDisponibilidadeNaoEncontrada = errors.New("disponibilidade não encontrada")
	ErrBloqueioNaoEncontrado        = errors.New("bloqueio não encontrado")
	ErrAgendamentoNoPassado         = errors.New("não é possível agendar um horário no passado")
	ErrHorarioIndisponivel          = errors.New("o faxineiro não está disponível neste horário")
	ErrConflitoAgenda               = errors.New("já existe um serviço agendado neste horário")
	ErrNaoEFaxineiroDoBloqueio      = errors.New("apenas o faxineiro dono do bloqueio pode realizar esta ação")
	ErrBloqueioPessoalApenas        = errors.New("apenas bloqueios pessoais podem ser removidos por esta ação")

	// Avaliacao
	ErrAvaliacaoNaoEncontrada = errors.New("avaliação não encontrada")
	ErrAvaliacaoDuplicada     = errors.New("já existe uma avaliação para esta solicitação")
	ErrSolicitacaoNaoAceita   = errors.New("a solicitação não está no estado aceita para ser avaliada")
)
