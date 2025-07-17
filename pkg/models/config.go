package models

type DatabaseConfig struct {
	Host     string `env:"DB_HOST"`
	Port     string `env:"DB_PORT"`
	User     string `env:"DB_USERNAME"`
	Password string `env:"DB_PASSWORD"`
	DB       string `env:"DB_NAME"`
}

type JobServiceConfig struct {
	Host string `env:"JOB_SERVICE_HOST"`
	Port string `env:"JOB_SERVICE_PORT"`
}

type CassandraConfig struct {
	Host     string `env:"CASS_HOST"`
	Port     string `env:"CASS_PORT"`
	Keyspace string `env:"CASS_KEYSPACE"`
}

type Config struct {
	Host             string `env:"HOST" envDefault:"0.0.0.0"`
	Port             string `env:"PORT" envDefault:"8080"`
	APIKeyParam      string `env:"API_KEY_PARAM"`
	BaseEndpoint     string `env:"BASE_ENDPOINT"`
	JWTSecretKey     string `env:"JWT_SECRET_KEY"`
	TaskTableName    string `env:"TASK_TABLE_NAME"`
	TaskTableVersion string `env:"TASK_TABLE_VERSION"`
	Env              string `env:"ENV" envDefault:"development"`
	Database         DatabaseConfig
	JobService       JobServiceConfig
	Cassandra        CassandraConfig
}
