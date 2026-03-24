CREATE TABLE limpezas (
    id               SERIAL PRIMARY KEY,
    nome             VARCHAR(255) NOT NULL,
    descricao        TEXT         NOT NULL DEFAULT '',
    valor_hora       NUMERIC(10,2) NOT NULL CHECK (valor_hora > 0),
    duracao_estimada NUMERIC(10,2) NOT NULL CHECK (duracao_estimada > 0),
    tipo_limpeza     VARCHAR(30)  NOT NULL,
    faxineiro_id     INT          NOT NULL REFERENCES usuarios(id) ON DELETE CASCADE,
    criado_em        TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    atualizado_em    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_limpezas_faxineiro_id ON limpezas(faxineiro_id);
CREATE INDEX idx_limpezas_atualizado_em ON limpezas(atualizado_em DESC);
