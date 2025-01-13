package domain

// UserRequest represents the structure for user-related requests such as login and registration.
// Username and Password fields are required for authentication processes.
// TurnstileToken is used for handling CAPTCHA verification to enhance security.
type UserRequest struct {
	Username       string `json:"username" binding:"required"`
	Password       string `json:"password" binding:"required"`
	TurnstileToken string `json:"cf-turnstile-response"`
}
