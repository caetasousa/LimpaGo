CREATE TABLE credenciais (
    usuario_id    INT PRIMARY KEY REFERENCES usuarios(id) ON DELETE CASCADE,
    senha_hash    VARCHAR(255) NOT NULL,
    criado_em     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    atualizado_em TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
