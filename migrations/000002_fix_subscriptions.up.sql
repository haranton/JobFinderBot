ALTER TABLE subscriptions
DROP CONSTRAINT IF EXISTS subscriptions_telegram_id_key;

ALTER TABLE subscriptions
ADD CONSTRAINT unique_user_search UNIQUE (telegram_id, search_text);