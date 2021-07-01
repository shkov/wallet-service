CREATE TABLE IF NOT EXISTS accounts (
  id BIGINT PRIMARY KEY,
  balance VARCHAR(32),
  created_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS payments (
  id BIGSERIAL PRIMARY KEY,
  from_account_id BIGINT REFERENCES accounts (id),
  to_account_id BIGINT REFERENCES accounts (id),
  amount VARCHAR(32),
  created_at TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS payments_from_account_id_idx on payments (from_account_id);

CREATE INDEX IF NOT EXISTS payments_to_account_id_idx on payments (to_account_id);
