ALTER TABLE subscriptions
DROP CONSTRAINT IF EXISTS unique_user_search;

ALTER TABLE subscriptions
ADD CONSTRAINT subscriptions_telegram_id_key UNIQUE (telegram_id);