package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"tgbot/internal/dto"
)

type Poller struct {
	client *http.Client
	token  string
}

func NewPoller(client *http.Client, token string) *Poller {
	return &Poller{client: client, token: token}
}

func (p *Poller) Start(ctx context.Context, out chan<- dto.Update) error {
	offset := 0
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		updates, err := p.getUpdates(ctx, offset)
		if err != nil {
			// простой экспоненциальный бэкофф
			select {
			case <-time.After(1 * time.Second):
			case <-ctx.Done():
				return nil
			}
			continue
		}
		for _, u := range updates {
			select {
			case out <- u:
				offset = u.UpdateID + 1
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func (p *Poller) getUpdates(ctx context.Context, offset int) ([]dto.Update, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s%s/getUpdates?offset=%d&timeout=30", telegramAPI, p.token, offset), nil)
	if err != nil {
		return nil, err
	}
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result struct {
		OK     bool         `json:"ok"`
		Result []dto.Update `json:"result"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return result.Result, nil
}
