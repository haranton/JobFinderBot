package service

import (
	"fmt"
	"strconv"
	"tgbot/internal/dto"
	"tgbot/internal/fetcher"
	"tgbot/internal/models"
	"tgbot/internal/repo"
)

type Service struct {
	repo    *repo.Repository
	fetcher *fetcher.Fetcher
}

func NewService(repo *repo.Repository, fetcher *fetcher.Fetcher) *Service {
	return &Service{repo: repo, fetcher: fetcher}
}

//
// ========== SUBSCRIPTIONS ==========
//

// Subscriptions возвращает все активные подписки
func (s *Service) Subscriptions() ([]models.Subscription, error) {
	return s.repo.Subscriptions()
}

func (s *Service) SubscriptionsUser(telegramID int) ([]models.Subscription, error) {
	return s.repo.SubscriptionsUser(telegramID)
}

// RegisterSubscribe добавляет новую подписку для пользователя
func (s *Service) RegisterSubscribe(telegramID int, query string) (models.Subscription, error) {
	user, err := s.repo.GetUser(telegramID)
	if err != nil {
		return models.Subscription{}, fmt.Errorf("failed to get user: %w", err)
	}
	if user == (models.User{}) {
		return models.Subscription{}, fmt.Errorf("user not found: id=%d", telegramID)
	}

	sub, err := s.repo.CreateSubscribe(telegramID, query)
	if err != nil {
		return models.Subscription{}, fmt.Errorf("failed to create subscription: %w", err)
	}
	return sub, nil
}

//
// ========== VACANCIES ==========
//

// SearchVacancies ищет вакансии по запросу, фильтрует уже просмотренные и сохраняет новые
func (s *Service) SearchVacancies(query string, telegramID int) ([]dto.Vacancy, error) {
	// Получаем вакансии с hh.ru
	hhVacancies, err := s.fetcher.Vacancies(query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch vacancies: %w", err)
	}

	// Преобразуем ID в int для записи в базу
	vacancies := make([]models.Vacancy, 0, len(hhVacancies))
	for _, hhVac := range hhVacancies {
		idInt, convErr := strconv.Atoi(hhVac.Id)
		if convErr != nil {
			continue // пропускаем некорректный ID
		}
		vacancies = append(vacancies, models.Vacancy{ID: idInt})
	}

	// Получаем уже сохранённые вакансии пользователя
	userVacancies, err := s.repo.GetUserVacancies(telegramID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user vacancies: %w", err)
	}

	// Фильтруем новые
	newVacancies := filterNewVacancies(hhVacancies, userVacancies)

	// Сохраняем вакансии для пользователя
	if err := s.repo.SaveVacancies(telegramID, vacancies); err != nil {
		return nil, fmt.Errorf("failed to save vacancies: %w", err)
	}

	return newVacancies, nil
}

//
// ========== USERS ==========
//

// RegisterUser создаёт нового пользователя
func (s *Service) RegisterUser(telegramID int) (models.User, error) {
	user, err := s.repo.GetUser(telegramID)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	if user != (models.User{}) {
		return models.User{}, fmt.Errorf("user already exists: id=%d", telegramID)
	}

	newUser, err := s.repo.CreateUser(telegramID)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to create user: %w", err)
	}
	return newUser, nil
}

// DeleteSubscribes удаляет все подписки пользователя
func (s *Service) DeleteSubscribes(telegramID int) error {

	isExist, err := s.userIsExist(telegramID)

	if err != nil {
		return fmt.Errorf("failed check user exist %d: %w", telegramID, err)
	}

	if !isExist {
		return fmt.Errorf("failed user dont exist %d: %w", telegramID, err)
	}
	if err := s.repo.DeleteUserSubscriptions(telegramID); err != nil {
		return fmt.Errorf("failed to delete subscriptions for user %d: %w", telegramID, err)
	}
	return nil
}

// UserIsExist проверяет, существует ли пользователь
func (s *Service) userIsExist(telegramID int) (bool, error) {
	user, err := s.repo.GetUser(telegramID)
	if err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}
	return user != (models.User{}), nil
}

//
// ========== HELPERS ==========
//

// filterNewVacancies фильтрует вакансии, которые пользователь ещё не видел
func filterNewVacancies(all []dto.Vacancy, userVacancies []models.UserVacancy) []dto.Vacancy {
	seen := make(map[int]struct{}, len(userVacancies))
	for _, uv := range userVacancies {
		seen[uv.VacancyID] = struct{}{}
	}

	var filtered []dto.Vacancy
	for _, v := range all {
		idInt, err := strconv.Atoi(v.Id)
		if err != nil {
			continue
		}
		if _, exists := seen[idInt]; !exists {
			filtered = append(filtered, v)
		}
	}

	return filtered
}
