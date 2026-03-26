CREATE TABLE bloqueios (
    id             SERIAL PRIMARY KEY,
    profissional_id   INT         NOT NULL REFERENCES usuarios(id) ON DELETE CASCADE,
    solicitacao_id INT REFERENCES solicitacoes(id) ON DELETE SET NULL,
    data_inicio    TIMESTAMPTZ NOT NULL,
    data_fim       TIMESTAMPTZ NOT NULL,
    criado_em      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CHECK (data_fim > data_inicio)
);

CREATE INDEX idx_bloqueios_profissional ON bloqueios(profissional_id);
CREATE INDEX idx_bloqueios_profissional_periodo ON bloqueios(profissional_id, data_inicio, data_fim);
CREATE INDEX idx_bloqueios_solicitacao ON bloqueios(solicitacao_id)
    WHERE solicitacao_id IS NOT NULL;
