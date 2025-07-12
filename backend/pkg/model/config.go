package model

type Config struct {
	Host             string `env:"HOST"`
	Port             string `env:"PORT" envDefault:"8080"`
	APIKeyParam      string `env:"API_KEY_PARAM"`
	BaseEndpoint     string `env:"BASE_ENDPOINT"`
	TaskTableName    string `env:"TASK_TABLE_NAME"`
	TaskTableVersion string `env:"TASK_TABLE_VERSION"`
	Env              string `env:"ENV" envDefault:"development"`
}
