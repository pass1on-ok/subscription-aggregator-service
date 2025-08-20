CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS subscriptions (
  id           uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
  service_name text        NOT NULL,
  price        integer     NOT NULL CHECK (price >= 0),
  user_id      uuid        NOT NULL,
  start_date   date        NOT NULL,
  end_date     date        NULL,
  created_at   timestamptz NOT NULL DEFAULT now(),
  updated_at   timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_subscriptions_user ON subscriptions(user_id);
CREATE INDEX IF NOT EXISTS idx_subscriptions_service ON subscriptions(service_name);
CREATE INDEX IF NOT EXISTS idx_subscriptions_start_end ON subscriptions(start_date, end_date);
