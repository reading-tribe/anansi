package nettypes

type RefreshRequest struct {
	EmailAddress string `json:"emailAddress"`
	Key          string `json:"key"`
}
