CREATE TABLE IF NOT EXISTS rate
(
    id               SERIAL PRIMARY KEY NOT NULL,
    first_currency   VARCHAR(5)         NOT NULL,
    second_currency  VARCHAR(5)         NOT NULL,
    value            REAL               NOT NULL,
    last_update_time TIMESTAMP          NOT NULL
);