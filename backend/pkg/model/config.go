package model

type DatabaseConfig struct {
	Host     string `env:"DB_HOST"`
	Port     string `env:"DB_PORT"`
	User     string `env:"DB_USERNAME"`
	Password string `env:"DB_PASSWORD"`
	DB       string `env:"DB_NAME"`
}

type Config struct {
	Host             string `env:"HOST"`
	Port             string `env:"PORT" envDefault:"8080"`
	APIKeyParam      string `env:"API_KEY_PARAM"`
	BaseEndpoint     string `env:"BASE_ENDPOINT"`
	TaskTableName    string `env:"TASK_TABLE_NAME"`
	TaskTableVersion string `env:"TASK_TABLE_VERSION"`
	Env              string `env:"ENV" envDefault:"development"`
	Database         DatabaseConfig
}
