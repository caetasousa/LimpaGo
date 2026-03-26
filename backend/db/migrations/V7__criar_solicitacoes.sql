CREATE TABLE solicitacoes (
    id                   SERIAL PRIMARY KEY,
    cliente_id           INT          NOT NULL REFERENCES usuarios(id) ON DELETE CASCADE,
    limpeza_id           INT          NOT NULL REFERENCES limpezas(id) ON DELETE CASCADE,
    status               VARCHAR(20)  NOT NULL DEFAULT 'pendente',
    data_agendada        TIMESTAMPTZ  NOT NULL,
    preco_total          NUMERIC(10,2) NOT NULL,
    multa_cancelamento   NUMERIC(10,2) NOT NULL DEFAULT 0,
    endereco_rua         VARCHAR(255) NOT NULL DEFAULT '',
    endereco_complemento VARCHAR(255) NOT NULL DEFAULT '',
    endereco_bairro      VARCHAR(255) NOT NULL DEFAULT '',
    endereco_cidade      VARCHAR(255) NOT NULL DEFAULT '',
    endereco_estado      VARCHAR(2)   NOT NULL DEFAULT '',
    endereco_cep         VARCHAR(9)   NOT NULL DEFAULT '',
    criado_em            TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    atualizado_em        TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_solicitacoes_cliente_id ON solicitacoes(cliente_id);
CREATE INDEX idx_solicitacoes_limpeza_id ON solicitacoes(limpeza_id);
CREATE INDEX idx_solicitacoes_cliente_limpeza ON solicitacoes(cliente_id, limpeza_id);
CREATE INDEX idx_solicitacoes_ativa ON solicitacoes(cliente_id, limpeza_id)
    WHERE status IN ('pendente', 'aceita');
