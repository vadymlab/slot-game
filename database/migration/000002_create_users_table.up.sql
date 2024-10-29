CREATE TABLE "users"
(
    id          SERIAL PRIMARY KEY,
    external_id UUID                NOT NULL DEFAULT uuid_generate_v4(),
    login       VARCHAR(255) UNIQUE NOT NULL,
    password    VARCHAR(255)        NOT NULL,
    balance     NUMERIC                      DEFAULT NULL,
    created_at  TIMESTAMP           NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP           NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMP

);

CREATE INDEX idx_user_deletedat ON "users" (deleted_at);