package fetcher

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"tgbot/internal/dto"
)

const baseUrl = "https://api.hh.ru/vacancies"

type Fetcher struct {
	client *http.Client
}

func NewFetcher(client *http.Client) *Fetcher {
	return &Fetcher{
		client: client,
	}
}

func (f *Fetcher) Vacancies(userInput string) ([]dto.Vacancy, error) {
	searchQuery := buildSearchQuery(userInput)

	params := url.Values{}
	params.Add("text", searchQuery)
	params.Add("per_page", "100") // максимум 100 вакансий за один запрос

	allVacancies := make([]dto.Vacancy, 0)
	page := 0

	for {
		params.Set("page", fmt.Sprintf("%d", page))
		fullURL := baseUrl + "?" + params.Encode()

		req, err := http.NewRequest(http.MethodGet, fullURL, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		resp, err := f.client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to get response: %w", err)
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to read response: %w", err)
		}

		var responseModel dto.HHResponse
		if err = json.Unmarshal(body, &responseModel); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}

		// Добавляем вакансии из этой страницы
		allVacancies = append(allVacancies, responseModel.Items...)

		// Проверяем, последняя ли страница
		if responseModel.Pages == 0 || page >= responseModel.Pages-1 {
			break
		}

		page++
	}

	return allVacancies, nil
}

func buildSearchQuery(userInput string) string {

	query := strings.TrimSpace(userInput)

	words := strings.Fields(query)
	if len(words) > 1 {
		query = strings.Join(words, " AND ")
	}

	return query
}
