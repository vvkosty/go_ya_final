CREATE TABLE IF NOT EXISTS withdraw_history
(
    id         SMALLSERIAL,
    user_id    SMALLSERIAL NOT NULL,
    order_id   varchar(20),
    withdraw   float       NOT NULL,
    created_at date        NOT NULL
)