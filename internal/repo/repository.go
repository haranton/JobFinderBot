package repo

import (
	"context"
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
	return &Repository{db: db}
}

// CreateSubscribe добавляет новую подписку для пользователя.
func (repo *Repository) CreateSubscribe(ctx context.Context, userId int, subscribeQuery string) (models.Subscription, error) {
	var sub models.Subscription

	query := `
		INSERT INTO subscriptions (telegram_id, search_text)
		VALUES ($1, $2)
		RETURNING id, telegram_id, search_text;
	`

	if err := repo.db.GetContext(ctx, &sub, query, userId, subscribeQuery); err != nil {
		return models.Subscription{}, fmt.Errorf("failed to create subscription: %w", err)
	}

	return sub, nil
}

// Subscriptions возвращает все подписки из базы.
func (repo *Repository) Subscriptions(ctx context.Context) ([]models.Subscription, error) {
	var subscriptions []models.Subscription

	query := `
		SELECT search_text, telegram_id
		FROM subscriptions
	`

	if err := repo.db.SelectContext(ctx, &subscriptions, query); err != nil {
		return nil, fmt.Errorf("failed to get subscriptions: %w", err)
	}

	return subscriptions, nil
}

// SubscriptionsUser возвращает подписки конкретного пользователя.
func (repo *Repository) SubscriptionsUser(ctx context.Context, userId int) ([]models.Subscription, error) {
	var subscriptions []models.Subscription

	query := `
		SELECT search_text, telegram_id
		FROM subscriptions
		WHERE telegram_id = $1
	`

	if err := repo.db.SelectContext(ctx, &subscriptions, query, userId); err != nil {
		return nil, fmt.Errorf("failed to get subscriptions for user: %w", err)
	}

	return subscriptions, nil
}

// GetUserVacancies возвращает вакансии, связанные с пользователем.
func (repo *Repository) GetUserVacancies(ctx context.Context, userId int) ([]models.UserVacancy, error) {
	var userVacancies []models.UserVacancy

	query := `
		SELECT vacancy_id 
		FROM user_vacancies 
		WHERE telegram_id = $1
	`

	if err := repo.db.SelectContext(ctx, &userVacancies, query, userId); err != nil {
		return nil, fmt.Errorf("failed to get user vacancies: %w", err)
	}

	return userVacancies, nil
}

// SaveVacancies сохраняет вакансии и их связи с пользователем в транзакции.
func (repo *Repository) SaveVacancies(ctx context.Context, userId int, vacancies []models.Vacancy) error {
	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	// Если возникнет ошибка — откатываем транзакцию
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	for _, v := range vacancies {
		// 1. Сохраняем вакансию (если нет — добавляем)
		if _, err = tx.ExecContext(ctx, `
			INSERT INTO vacancies (id)
			VALUES ($1)
			ON CONFLICT (id) DO NOTHING
		`, v.ID); err != nil {
			return fmt.Errorf("failed to insert vacancy: %w", err)
		}

		// 2. Сохраняем связь между пользователем и вакансией
		if _, err = tx.ExecContext(ctx, `
			INSERT INTO user_vacancies (telegram_id, vacancy_id)
			VALUES ($1, $2)
			ON CONFLICT (telegram_id, vacancy_id) DO NOTHING
		`, userId, v.ID); err != nil {
			return fmt.Errorf("failed to link user and vacancy: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetUser получает пользователя по telegram_id.
func (repo *Repository) GetUser(ctx context.Context, userId int) (models.User, error) {
	var user models.User

	query := `
		SELECT telegram_id, created_at
		FROM users
		WHERE telegram_id = $1
	`

	err := repo.db.GetContext(ctx, &user, query, userId)
	if errors.Is(err, sql.ErrNoRows) {
		return models.User{}, nil
	}
	if err != nil {
		return models.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// CreateUser создаёт нового пользователя.
func (repo *Repository) CreateUser(ctx context.Context, userId int) (models.User, error) {
	var user models.User

	query := `
		INSERT INTO users (telegram_id)
		VALUES ($1)
		ON CONFLICT (telegram_id) DO NOTHING
		RETURNING telegram_id, created_at
	`

	if err := repo.db.GetContext(ctx, &user, query, userId); err != nil {
		return models.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// DeleteUserSubscriptions удаляет все подписки пользователя.
func (repo *Repository) DeleteUserSubscriptions(ctx context.Context, userId int) error {
	query := `
		DELETE FROM subscriptions 
		WHERE telegram_id = $1
	`

	if _, err := repo.db.ExecContext(ctx, query, userId); err != nil {
		return fmt.Errorf("failed to delete subscriptions for user %d: %w", userId, err)
	}

	return nil
}

func (r *Repository) Close() error {
	if r.db != nil {
		return r.db.Close()
	}
	return nil
}
