CREATE TABLE perfis_profissional (
    usuario_id        INT PRIMARY KEY REFERENCES usuarios(id) ON DELETE CASCADE,
    descricao         TEXT         NOT NULL DEFAULT '',
    anos_experiencia  INT          NOT NULL DEFAULT 0,
    especialidades    TEXT[]       NOT NULL DEFAULT '{}',
    cidades_atendidas TEXT[]       NOT NULL DEFAULT '{}',
    documento_rg      VARCHAR(20)  NOT NULL DEFAULT '',
    documento_cpf     VARCHAR(14)  NOT NULL DEFAULT '',
    foto_documento    TEXT         NOT NULL DEFAULT '',
    verificado        BOOLEAN      NOT NULL DEFAULT FALSE,
    criado_em         TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    atualizado_em     TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
