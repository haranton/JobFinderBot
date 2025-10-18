package repo

import (
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

func (repo *Repository) GetUserVacancies(chatId int) ([]models.UserVacancy, error) {
	var userVacancy []models.UserVacancy
	query := `
		SELECT vacancy_id from user_vacancies where telegram_id = $1
	`
	err := repo.db.Select(&userVacancy, query, chatId)
	if err != nil {
		return nil, fmt.Errorf("failed to get vacansies user from db, err: %w", err)
	}

	return userVacancy, nil
}

func (repo *Repository) SaveVacancies(telegramID int, vacancies []models.Vacancy) error {
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
		`, telegramID, v.ID)
		if err != nil {
			return fmt.Errorf("failed to link user and vacancy: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
