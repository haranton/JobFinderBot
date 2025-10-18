package models

import "time"

type User struct {
	ID         int       `db:"id" json:"id"`
	TelegramID int64     `db:"telegram_id" json:"telegram_id"`
	ChatID     int64     `db:"chat_id" json:"chat_id"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

type Vacancy struct {
	ID        int       `db:"id" json:"id"`
	VacancyID string    `db:"vacancy_id" json:"vacancy_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type UserVacancy struct {
	ID        int `db:"id" json:"id"`
	UserID    int `db:"user_id" json:"user_id"`
	VacancyID int `db:"vacancy_id" json:"vacancy_id"`
}

type Subscription struct {
	ID         int       `db:"id" json:"id"`
	SearchText string    `db:"search_text" json:"search_text"`
	UserID     int       `db:"user_id" json:"user_id"`
	IsActive   bool      `db:"is_active" json:"is_active"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}
