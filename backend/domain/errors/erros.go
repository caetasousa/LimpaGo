package errors

import "errors"

var (
	// Usuario
	ErrEmailJaUtilizado       = errors.New("email já está em uso")
	ErrNomeUsuarioJaUtilizado = errors.New("nome de usuário já está em uso")
	ErrUsuarioNaoEncontrado   = errors.New("usuário não encontrado")

	// Perfil
	ErrPerfilNaoEncontrado          = errors.New("perfil não encontrado")
	ErrPerfilProfissionalNaoEncontrado = errors.New("perfil de profissional não encontrado")
	ErrPerfilClienteNaoEncontrado   = errors.New("perfil de cliente não encontrado")
	ErrPerfilProfissionalJaExiste      = errors.New("usuário já possui um perfil de profissional")
	ErrPerfilClienteJaExiste        = errors.New("usuário já possui um perfil de cliente")

	// Limpeza (Serviço publicado pelo profissional)
	ErrLimpezaNaoEncontrada      = errors.New("serviço de limpeza não encontrado")
	ErrNaoEProfissionalDaLimpeza = errors.New("apenas o profissional que publicou o serviço pode realizar esta ação")

	// Solicitacao (Cliente solicita o serviço de um profissional)
	ErrSolicitacaoNaoEncontrada            = errors.New("solicitação não encontrada")
	ErrSolicitacaoDuplicada                = errors.New("cliente já possui uma solicitação para este serviço")
	ErrProfissionalNaoPodeSolicitarProprio = errors.New("o profissional não pode solicitar o próprio serviço")
	ErrSolicitacaoNaoPodeSerAceita         = errors.New("apenas solicitações pendentes podem ser aceitas")
	ErrSolicitacaoNaoPodeSerCancelada      = errors.New("esta solicitação não pode ser cancelada no estado atual")
	ErrSolicitacaoNaoPodeSerRejeitada      = errors.New("apenas solicitações pendentes podem ser rejeitadas")
	ErrNaoEClienteSolicitante              = errors.New("apenas o cliente que fez a solicitação pode realizar esta ação")
	ErrNaoEProfissionalDaSolicitacao       = errors.New("apenas o profissional do serviço pode aceitar ou rejeitar solicitações")

	// Agenda
	ErrDisponibilidadeNaoEncontrada = errors.New("disponibilidade não encontrada")
	ErrBloqueioNaoEncontrado        = errors.New("bloqueio não encontrado")
	ErrAgendamentoNoPassado         = errors.New("não é possível agendar um horário no passado")
	ErrHorarioIndisponivel          = errors.New("o profissional não está disponível neste horário")
	ErrConflitoAgenda               = errors.New("já existe um serviço agendado neste horário")
	ErrNaoEProfissionalDoBloqueio   = errors.New("apenas o profissional dono do bloqueio pode realizar esta ação")
	ErrBloqueioPessoalApenas        = errors.New("apenas bloqueios pessoais podem ser removidos por esta ação")

	// Avaliacao
	ErrAvaliacaoNaoEncontrada = errors.New("avaliação não encontrada")
	ErrAvaliacaoDuplicada     = errors.New("já existe uma avaliação para esta solicitação")
	ErrSolicitacaoNaoAceita   = errors.New("a solicitação não está no estado aceita para ser avaliada")
)
