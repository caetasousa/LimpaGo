CREATE TABLE perfis (
    usuario_id    INT PRIMARY KEY REFERENCES usuarios(id) ON DELETE CASCADE,
    nome_completo VARCHAR(255) NOT NULL DEFAULT '',
    telefone      VARCHAR(20)  NOT NULL DEFAULT '',
    imagem        TEXT         NOT NULL DEFAULT '',
    email         VARCHAR(255) NOT NULL DEFAULT '',
    nome_usuario  VARCHAR(100) NOT NULL DEFAULT '',
    criado_em     TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    atualizado_em TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
