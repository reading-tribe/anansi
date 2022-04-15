package nettypes

type LogoutRequest struct {
	EmailAddress string `json:"emailAddress"`
	Key          string `json:"key"`
}
