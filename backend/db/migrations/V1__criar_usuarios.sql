CREATE TABLE usuarios (
    id               SERIAL PRIMARY KEY,
    email            VARCHAR(255) NOT NULL UNIQUE,
    nome_usuario     VARCHAR(100) NOT NULL UNIQUE,
    email_verificado BOOLEAN NOT NULL DEFAULT FALSE,
    ativo            BOOLEAN NOT NULL DEFAULT TRUE,
    super_usuario    BOOLEAN NOT NULL DEFAULT FALSE,
    criado_em        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    atualizado_em    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_usuarios_email ON usuarios(email);
CREATE INDEX idx_usuarios_nome_usuario ON usuarios(nome_usuario);
