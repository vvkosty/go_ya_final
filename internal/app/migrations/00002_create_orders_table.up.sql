CREATE TABLE IF NOT EXISTS orders
(
    id          varchar(20) unique NOT NULL,
    user_id     SMALLSERIAL NOT NULL,
    accrual     float,
    status      varchar(50),
    uploaded_at timestamp
)