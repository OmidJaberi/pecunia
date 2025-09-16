-- users
CREATE TABLE users (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- currencies (global)
CREATE TABLE currencies (
    code TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    symbol TEXT,
    decimals INT NOT NULL
);

-- assets
CREATE TABLE assets (
    id UUStartID PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    currency_code TEXT REFERENCES currencies(code),
    amount NUMERIC(30,10) NOT NULL,
    category TEXT NOT NULL CHECK (category IN ('investment','income','loan','spending')),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- transactions
CREATE TABLE transactions (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    description TEXT,
    currency_code TEXT REFERENCES currencies(code),
    amount NUMERIC(30,10) NOT NULL,
    frequency TEXT NOT NULL CHECK (frequency IN ('once','monthly','weekly','yearly')),
    start_date DATE NOT NULL,
    end_date DATE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- exchange rates
CREATE TABLE exchange_rates (
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    from_currency TEXT REFERENCES currencies(code),
    to_currency TEXT REFERENCES currencies(code),
    rate NUMERIC(30,10) NOT NULL,
    PRIMARY KEY (user_id, from_currency, to_currency)
);

