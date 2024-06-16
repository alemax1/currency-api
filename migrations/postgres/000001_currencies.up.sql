CREATE TYPE currency_types AS ENUM ('fiat', 'crypto');

CREATE TABLE IF NOT EXISTS currencies(
    id SERIAL PRIMARY KEY,
    type currency_types NOT NULL,
    name VARCHAR UNIQUE NOT NULL,
    value_usd DECIMAL NOT NULL,
    is_available BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- INIT DATA
INSERT INTO currencies(type, name, value_usd, is_available) VALUES('fiat', 'USD', 0, false);
INSERT INTO currencies(type, name, value_usd, is_available) VALUES('fiat', 'EUR', 0, false);
INSERT INTO currencies(type, name, value_usd, is_available) VALUES('fiat', 'CNY', 0, false);
INSERT INTO currencies(type, name, value_usd, is_available) VALUES('crypto', 'USDT', 0, false);
INSERT INTO currencies(type, name, value_usd, is_available) VALUES('crypto', 'USDC', 0, false);
INSERT INTO currencies(type, name, value_usd, is_available) VALUES('crypto', 'ETH', 0, false);