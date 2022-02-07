CREATE TABLE IF NOT EXISTS exchange_rate
(
    id               SERIAL PRIMARY KEY       NOT NULL,
    first_currency   VARCHAR(5)               NOT NULL,
    second_currency  VARCHAR(5)               NOT NULL,
    rate_value       INTEGER                  NOT NULL,
    last_update_time TIMESTAMP WITH TIME ZONE NOT NULL
);