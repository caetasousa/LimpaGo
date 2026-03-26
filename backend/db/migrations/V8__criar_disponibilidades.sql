CREATE TABLE disponibilidades (
    id            SERIAL PRIMARY KEY,
    profissional_id  INT         NOT NULL REFERENCES usuarios(id) ON DELETE CASCADE,
    dia_semana    INT         NOT NULL CHECK (dia_semana BETWEEN 0 AND 6),
    hora_inicio   INT         NOT NULL CHECK (hora_inicio BETWEEN 0 AND 23),
    hora_fim      INT         NOT NULL CHECK (hora_fim BETWEEN 1 AND 24),
    criado_em     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    atualizado_em TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CHECK (hora_fim > hora_inicio)
);

CREATE INDEX idx_disponibilidades_profissional ON disponibilidades(profissional_id);
CREATE INDEX idx_disponibilidades_profissional_dia ON disponibilidades(profissional_id, dia_semana);
