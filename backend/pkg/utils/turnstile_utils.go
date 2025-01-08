package utils

import (
	"encoding/json"
	"errors"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
)

type TurnstileResponse struct {
	Success     bool     `json:"success"`
	ChallengeTS string   `json:"challenge_ts"`
	Hostname    string   `json:"hostname"`
	ErrorCodes  []string `json:"error-codes"`
}

// 验证 Turnstile token
func VerifyTurnstileToken(token string) error {
	secretKey := viper.GetString("cloudflare.turnstile-key")
	if secretKey == "" {
		return errors.New("missing Turnstile secret key")
	}

	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(
			map[string]string{
				"secret":   secretKey,
				"response": token,
			},
		).
		Post("https://challenges.cloudflare.com/turnstile/v0/siteverify")

	if err != nil {
		return err
	}

	var result TurnstileResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return err
	}

	if !result.Success {
		return errors.New("Turnstile verification failed")
	}

	return nil
}
