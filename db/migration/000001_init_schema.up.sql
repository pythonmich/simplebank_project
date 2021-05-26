CREATE TABLE IF NOT EXISTS accounts(
       id BIGSERIAL PRIMARY KEY,
       owner VARCHAR(100) NOT NULL UNIQUE,
       balance REAL NOT NULL,
       currency VARCHAR(20) NOT NULL,
       created_at TIMESTAMPTZ NOT NULL DEFAULT (now())
);
-- entries is used to record all changes to the account balance
CREATE TABLE IF NOT EXISTS entries(
      id BIGSERIAL PRIMARY KEY,
      account_id BIGINT NOT NULL,
      FOREIGN KEY (account_id) REFERENCES accounts(id),
      amount REAL NOT NULL, -- it can be positive or negative
      created_at TIMESTAMPTZ NOT NULL DEFAULT (now())
);
-- transfers it records all the money transfers btwn two accounts
CREATE TABLE IF NOT EXISTS transfers (
         id BIGSERIAL PRIMARY KEY,
         from_account_id BIGINT NOT NULL,
         FOREIGN KEY (from_account_id) REFERENCES accounts(id),
         to_account_id BIGINT NOT NULL,
         FOREIGN KEY (to_account_id) REFERENCES accounts(id),
         amount REAL NOT NULL, -- it must be positive
         created_at TIMESTAMPTZ NOT NULL DEFAULT (now())
);

CREATE UNIQUE INDEX acc_owner ON accounts(owner);

CREATE INDEX acc_id ON entries(account_id);

CREATE INDEX from_acc_id ON transfers(from_account_id);

CREATE INDEX to_acc_id ON transfers(to_account_id);

CREATE  INDEX both_acc_id ON transfers(from_account_id, to_account_id);

COMMENT ON COLUMN entries.amount IS 'can be positive or negative';

COMMENT ON COLUMN transfers.amount IS 'must be positive';


