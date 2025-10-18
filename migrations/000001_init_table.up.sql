CREATE TABLE IF NOT EXISTS users (
    telegram_id BIGINT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS vacancies (
    id INT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS user_vacancies (
    id SERIAL PRIMARY KEY,
    telegram_id BIGINT REFERENCES users(telegram_id) ON DELETE CASCADE,
    vacancy_id INT REFERENCES vacancies(id) ON DELETE CASCADE,
    UNIQUE(telegram_id, vacancy_id)
);

CREATE TABLE IF NOT EXISTS subscriptions (
    id SERIAL PRIMARY KEY,
    search_text TEXT NOT NULL,
    telegram_id BIGINT UNIQUE REFERENCES users(telegram_id) ON DELETE CASCADE,
    is_active BOOLEAN DEFAULT TRUE
);
