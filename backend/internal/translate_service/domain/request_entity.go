package domain

// OpenAIRequest represents a request payload sent to the OpenAI API containing model and messages data.
type OpenAIRequest struct {
	Model    string          `json:"model"`
	Messages []OpenAIMessage `json:"messages"`
}

// OpenAIMessage represents a message exchanged with the OpenAI API, including the role and content of the message.
type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenAIResponse represents a structured response returned from the OpenAI API containing choices and messages.
type OpenAIResponse struct {
	Choices []struct {
		Message OpenAIMessage `json:"message"`
	} `json:"choices"`
}
