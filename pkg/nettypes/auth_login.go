package nettypes

type LoginRequest struct {
	EmailAddress string `json:"emailAddress"`
	Password     string `json:"password"`
}

type LoginResponse struct {
	Token   string `json:"token"`
	Message string `json:"message"`
}
