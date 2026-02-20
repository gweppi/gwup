package shared

type Config struct {
	ServerUrl string `json:"server_url"`
	AuthCode string `json:"auth_code"`
}

func (c *Config) IsUndefined() bool {
	return c.ServerUrl == ""
}

const FileName = "the-file"
