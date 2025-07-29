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

type RabbitConfig struct {
	Host     string `env:"RABBIT_HOST"`
	Port     string `env:"RABBIT_PORT"`
	Username string `env:"DB_USERNAME"`
	Password string `env:"RABBIT_PASSWORD"`
	Name     string `env:"RABBIT_QUEUE_NAME"`
}

type Config struct {
	Host             string `env:"HOST" envDefault:"0.0.0.0"`
	Port             string `env:"PORT" envDefault:"8080"`
	Auth0Domain      string `env:"VITE_AUTH0_DOMAIN"`
	Auth0Audience    string `env:"VITE_AUTH0_AUDIENCE"`
	APIKeyParam      string `env:"API_KEY_PARAM"`
	BaseEndpoint     string `env:"BASE_ENDPOINT"`
	JWTSecretKey     string `env:"JWT_SECRET_KEY"`
	TaskTableName    string `env:"TASK_TABLE_NAME"`
	TaskTableVersion string `env:"TASK_TABLE_VERSION"`
	Env              string `env:"ENV" envDefault:"development"`
	TempDir          string `env:"TEMP_DIR"`
	SandboxCount     int    `env:"SANDBOX_COUNT"`
	WorkerID         string `env:"WORKER_ID"`
	JWKSUrl          string `env:"JWKS_URL"`
	Database         DatabaseConfig
	JobService       JobServiceConfig
	Cassandra        CassandraConfig
	Rabbit           RabbitConfig
}
