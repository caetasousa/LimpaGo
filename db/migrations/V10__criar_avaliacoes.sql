CREATE TABLE avaliacoes (
    id           SERIAL PRIMARY KEY,
    limpeza_id   INT         NOT NULL REFERENCES limpezas(id) ON DELETE CASCADE,
    faxineiro_id INT         NOT NULL REFERENCES usuarios(id) ON DELETE CASCADE,
    cliente_id   INT         NOT NULL REFERENCES usuarios(id) ON DELETE CASCADE,
    nota         INT         NOT NULL CHECK (nota BETWEEN 0 AND 5),
    comentario   TEXT        NOT NULL DEFAULT '',
    criado_em    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (cliente_id, limpeza_id)
);

CREATE INDEX idx_avaliacoes_faxineiro ON avaliacoes(faxineiro_id);
