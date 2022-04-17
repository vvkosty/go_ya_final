CREATE TABLE IF NOT EXISTS users
(
    id       SMALLSERIAL,
    login    varchar(50)  NOT NULL UNIQUE,
    password varchar(255) NOT NULL
)