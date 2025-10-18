package models

import "time"

// Таблица users
type User struct {
	TelegramID int64     `db:"telegram_id" json:"telegram_id"`
	ChatID     int64     `db:"chat_id" json:"chat_id"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

// Таблица vacancies
type Vacancy struct {
	ID        int       `db:"id" json:"id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// Таблица user_vacancies (связь M:N между users и vacancies)
type UserVacancy struct {
	ID         int   `db:"id" json:"id"`
	TelegramID int64 `db:"telegram_id" json:"telegram_id"`
	VacancyID  int   `db:"vacancy_id" json:"vacancy_id"`
}

// Таблица subscriptions
type Subscription struct {
	ID         int       `db:"id" json:"id"`
	SearchText string    `db:"search_text" json:"search_text"`
	TelegramID int64     `db:"telegram_id" json:"telegram_id"`
	IsActive   bool      `db:"is_active" json:"is_active"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}
