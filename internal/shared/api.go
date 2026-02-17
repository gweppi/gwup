package shared

type ServerInfo struct {
	Status string `json:"status"`
	Version string `json:"version"`
	RequiresAuth bool `json:"requires_auth"`
}
