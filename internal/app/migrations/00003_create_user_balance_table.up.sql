CREATE TABLE IF NOT EXISTS user_balance
(
    id         SMALLSERIAL,
    user_id    SMALLSERIAL NOT NULL,
    balance    float         NOT NULL,
    withdraw   float         NOT NULL,
    updated_at date        NOT NULL
)