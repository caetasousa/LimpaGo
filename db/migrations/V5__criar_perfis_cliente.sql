CREATE TABLE perfis_cliente (
    usuario_id             INT PRIMARY KEY REFERENCES usuarios(id) ON DELETE CASCADE,
    endereco_rua           VARCHAR(255)       NOT NULL DEFAULT '',
    endereco_complemento   VARCHAR(255)       NOT NULL DEFAULT '',
    endereco_bairro        VARCHAR(255)       NOT NULL DEFAULT '',
    endereco_cidade        VARCHAR(255)       NOT NULL DEFAULT '',
    endereco_estado        VARCHAR(2)         NOT NULL DEFAULT '',
    endereco_cep           VARCHAR(9)         NOT NULL DEFAULT '',
    tipo_imovel            VARCHAR(20)        NOT NULL DEFAULT '',
    quartos                INT                NOT NULL DEFAULT 0,
    banheiros              INT                NOT NULL DEFAULT 0,
    tamanho_imovel_m2      DOUBLE PRECISION   NOT NULL DEFAULT 0,
    observacoes            TEXT               NOT NULL DEFAULT '',
    profissional_preferido_id INT REFERENCES usuarios(id) ON DELETE SET NULL,
    criado_em              TIMESTAMPTZ        NOT NULL DEFAULT NOW(),
    atualizado_em          TIMESTAMPTZ        NOT NULL DEFAULT NOW()
);
