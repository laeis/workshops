package config

const (
	CtxAuthId = "auth_id"
	CtxToken  = "auth_token"
)

type AppConfig struct {
	*DBConfig
	*JWTConfig
}

type DBConfig struct {
	Host     string `env:"POSTGRES_HOST,default=localhost"`
	Port     string `env:"POSTGRES_PORT,default=5432"`
	User     string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
	Database string `env:"POSTGRES_DB"`
}

type JWTConfig struct {
	SecretKey       string `env:"JWT_SECRET_KEY"`
	ExpirationHours int64  `env:"JWT_EXPIRATION_HOURS,default=1"`
}
