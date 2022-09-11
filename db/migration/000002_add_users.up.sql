CREATE TABLE "users" (
             "username" varchar PRIMARY KEY,
             "hashed_password" varchar NOT NULL,
             "full_name" varchar NOT NULL,
             "email" varchar UNIQUE NOT NULL,
             "created_at" timestamptz NOT NULL DEFAULT (now()),
             "password_changed_at" timestamptz NOT NULL DEFAULT ('0001-01-01T00:00:00Z')
);

ALTER TABLE "account" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");

-- CREATE UNIQUE INDEX ON "account" ("owner", "currency");

ALTER TABLE account ADD CONSTRAINT owner_currency_key UNIQUE(owner, currency);