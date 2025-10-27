package ais

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type AIS struct {
	url string
}

type SendOrderExecutorRequest struct {
	OrderId    int `json:"order_id"`
	ExecutorId int `json:"executor_id"`
}

func (a *AIS) SendOrderExecutor(ctx context.Context, orderId, executorId int) error {
	const op = "AIS.SendOrderExecutor"

	body := bytes.NewBuffer(nil)

	if err := json.NewEncoder(body).Encode(&SendOrderExecutorRequest{
		OrderId:    orderId,
		ExecutorId: executorId,
	}); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, a.url, body)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s: %s", op, resp.Status)
	}

	return nil
}
