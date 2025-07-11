package model

type Config struct {
	Host             string `env:"HOST"`
	Port             string `env:"PORT" envDefault:"8080"`
	TaskTableName    string `env:"TASK_TABLE_NAME"`
	TaskTableVersion string `env:"TASK_TABLE_VERSION"`
	Env              string `env:"ENV" envDefault:"development"`
}
