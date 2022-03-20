package nettypes

type RegisterRequest struct {
	EmailAddress string `json:"emailAddress"`
	Password     string `json:"password"`
}

type RegisterResponse struct {
	Message string `json:"message"`
}
