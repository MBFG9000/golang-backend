package modules

type PostgreConfig struct {
	Host     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

type AuthMiddlewareConfig struct {
	ApiKeyHeader string
	ValidAPIKey  string
}
