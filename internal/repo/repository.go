package repo

import (
	"database/sql"
	"errors"
	"fmt"
	"tgbot/internal/models"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (repo *Repository) CreateOrUpdateSubscribe(userId int, subscribeQuery string) (models.Subscription, error) {
	var sub models.Subscription

	query := `
		INSERT INTO subscriptions (telegram_id, search_text)
		VALUES ($1, $2)
		ON CONFLICT (telegram_id)
		DO UPDATE SET 
			search_text = EXCLUDED.search_text
		RETURNING id, telegram_id, search_text;
	`

	err := repo.db.Get(&sub, query, userId, subscribeQuery)
	if err != nil {
		return models.Subscription{}, fmt.Errorf("failed to create or update subscription: %w", err)
	}

	return sub, nil
}

func (repo *Repository) Subscription(userId int) (models.Subscription, error) {
	var subscription models.Subscription

	query := `
		SELECT search_text
		FROM subscriptions
		WHERE telegram_id = $1
	`

	err := repo.db.Get(&subscription, query, userId)

	if errors.Is(err, sql.ErrNoRows) {
		return models.Subscription{}, nil
	}

	if err != nil {
		return models.Subscription{}, fmt.Errorf("failed to get subscription from db: %w", err)
	}

	return subscription, nil
}

func (repo *Repository) GetUserVacancies(userId int) ([]models.UserVacancy, error) {
	var userVacancy []models.UserVacancy
	query := `
		SELECT vacancy_id from user_vacancies where telegram_id = $1
	`
	err := repo.db.Select(&userVacancy, query, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get vacansies user from db, err: %w", err)
	}

	return userVacancy, nil
}

func (repo *Repository) SaveVacancies(userId int, vacancies []models.Vacancy) error {
	tx, err := repo.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	for _, v := range vacancies {
		// 1. Сохраняем вакансию (если нет — добавляем)
		_, err = tx.Exec(`
			INSERT INTO vacancies (id)
			VALUES ($1)
			ON CONFLICT (id) DO NOTHING
		`, v.ID)
		if err != nil {
			return fmt.Errorf("failed to insert vacancy: %w", err)
		}

		// 2. Сохраняем связь между пользователем и вакансией
		_, err = tx.Exec(`
			INSERT INTO user_vacancies (telegram_id, vacancy_id)
			VALUES ($1, $2)
			ON CONFLICT (telegram_id, vacancy_id) DO NOTHING
		`, userId, v.ID)
		if err != nil {
			return fmt.Errorf("failed to link user and vacancy: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (repo *Repository) GetUser(userId int) (models.User, error) {
	var user models.User

	query := `
		SELECT telegram_id, created_at
		FROM users
		WHERE telegram_id = $1
	`

	err := repo.db.Get(&user, query, userId)

	if errors.Is(err, sql.ErrNoRows) {
		return models.User{}, nil
	}

	if err != nil {
		return models.User{}, fmt.Errorf("failed to get user from db: %w", err)
	}

	return user, nil
}

func (repo *Repository) CreateUser(userId int) (models.User, error) {
	var user models.User

	query := `
		INSERT INTO users (telegram_id)
		VALUES ($1)
		ON CONFLICT (telegram_id) DO NOTHING
		RETURNING telegram_id, created_at
	`

	err := repo.db.Get(&user, query, userId)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}
