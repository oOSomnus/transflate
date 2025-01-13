package utils

import (
	"encoding/json"
	"errors"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
)

// TurnstileResponse represents the response structure from Cloudflare's Turnstile verification endpoint.
// Success indicates if the verification was successful.
// ChallengeTS is the timestamp of the challenge completion in ISO 8601 format.
// Hostname is the hostname of the website where the challenge originated.
// ErrorCodes contains a list of error codes returned if the verification fails.
type TurnstileResponse struct {
	Success     bool     `json:"success"`
	ChallengeTS string   `json:"challenge_ts"`
	Hostname    string   `json:"hostname"`
	ErrorCodes  []string `json:"error-codes"`
}

// VerifyTurnstileToken validates a Cloudflare Turnstile token by sending a verification request to the Turnstile API.
// Returns an error if the secret key is missing, the API request fails, or verification is unsuccessful.
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
