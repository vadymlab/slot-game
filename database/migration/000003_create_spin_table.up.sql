CREATE TABLE spins
(
    id         SERIAL PRIMARY KEY,
    user_id    INTEGER        NOT NULL,
    bet_amount NUMERIC(10, 2) NOT NULL,
    win_amount NUMERIC(10, 2) NOT NULL,
    created_at TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,

    -- Foreign key constraint to users table
    CONSTRAINT fk_user
        FOREIGN KEY (user_id)
            REFERENCES users (id)
            ON DELETE SET NULL
            ON UPDATE CASCADE
);
