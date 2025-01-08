package domain

type UserRequest struct {
	Username       string `json:"username" binding:"required"`
	Password       string `json:"password" binding:"required"`
	TurnstileToken string `json:"cf-turnstile-response" binding:"required"`
}
