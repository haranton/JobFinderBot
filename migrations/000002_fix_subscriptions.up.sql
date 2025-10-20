ALTER TABLE subscriptions
DROP CONSTRAINT IF EXISTS unique_subscription_per_user;

ALTER TABLE subscriptions
ADD CONSTRAINT unique_user_search UNIQUE (telegram_id, search_text);