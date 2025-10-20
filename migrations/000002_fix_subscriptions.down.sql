ALTER TABLE subscriptions
DROP CONSTRAINT IF EXISTS unique_user_search;

ALTER TABLE subscriptions
ADD CONSTRAINT unique_subscription_per_user UNIQUE (telegram_id);