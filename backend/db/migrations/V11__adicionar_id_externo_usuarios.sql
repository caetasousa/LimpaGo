-- Adiciona coluna id_externo para mapear o sub do Zitadel ao usuário interno.
-- Permite sincronização just-in-time na primeira autenticação via Zitadel.
ALTER TABLE usuarios ADD COLUMN IF NOT EXISTS id_externo VARCHAR(255) UNIQUE;

CREATE INDEX IF NOT EXISTS idx_usuarios_id_externo ON usuarios(id_externo);
