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
	return &Service{
		repo:    repo,
		fetcher: fetcher,
	}
}

func (s *Service) SearchVacancies(query string, userId int) ([]dto.Vacancy, error) {

	hhVacancies, err := s.fetcher.Vacancies(query)
	if err != nil {
		return nil, err
	}

	var vacancies []models.Vacancy
	for _, hhVac := range hhVacancies {

		IdInt, _ := strconv.Atoi(hhVac.Id)
		vacancy := models.Vacancy{
			ID: IdInt,
		}
		vacancies = append(vacancies, vacancy)
	}
	//проверяем повторящие элементы
	vakanciesUser, err := s.repo.GetUserVacancies(userId)
	if err != nil {
		return nil, err
	}

	hhVacancies = sortVacancies(hhVacancies, vakanciesUser)

	if err := s.repo.SaveVacancies(userId, vacancies); err != nil {
		return nil, err
	}

	return hhVacancies, nil
}

func (s *Service) RegisterSubscribe(userId int, subscribeQuery string) (models.Subscription, error) {

	user, err := s.repo.GetUser(userId)

	if err != nil {
		return models.Subscription{}, err
	}

	if user == (models.User{}) {
		return models.Subscription{}, fmt.Errorf("user dont find: user id: %v", userId)
	}

	subscribe, err := s.repo.CreateOrUpdateSubscribe(userId, subscribeQuery)
	if err != nil {
		return models.Subscription{}, err
	}

	return subscribe, nil

}

func (s *Service) RegisterUser(userId int) (models.User, error) {
	user, err := s.repo.GetUser(userId)

	if err != nil {
		return models.User{}, err
	}

	if user != (models.User{}) {
		return user, fmt.Errorf("user already exist")
	}

	user, err = s.repo.CreateUser(userId)
	if err != nil {
		return models.User{}, err
	}

	return user, nil

}

func sortVacancies(vacancies []dto.Vacancy, userVacancies []models.UserVacancy) []dto.Vacancy {

	seen := make(map[int]struct{}, len(userVacancies))
	for _, uv := range userVacancies {
		seen[uv.VacancyID] = struct{}{}
	}

	var filtered []dto.Vacancy
	for _, v := range vacancies {
		IdInt, _ := strconv.Atoi(v.Id)
		if _, exists := seen[IdInt]; !exists {
			filtered = append(filtered, v)
		}
	}

	return filtered
}

// func (s *Service) MarkVacancySent(userID int, vacancyID string) error {
// 	// Бизнес-логика отметки отправки
// 	return s.userVacancyRepo.MarkSent(userID, vacancyID)
// }
