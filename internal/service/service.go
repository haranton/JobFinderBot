package service

import (
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

func (s *Service) SearchVacancies(query string, chatId int) ([]dto.Vacancy, error) {

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
	vakanciesUser, err := s.repo.GetUserVacancies(chatId)
	if err != nil {
		return nil, err
	}

	hhVacancies = sortVacancies(hhVacancies, vakanciesUser)

	if err := s.repo.SaveVacancies(chatId, vacancies); err != nil {
		return nil, err
	}

	return hhVacancies, nil
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
