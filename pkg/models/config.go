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

type ScheduleServiceConfig struct {
	Auth0ClientID     string `env:"SCHEDULING_AUTH0_CLIENT_ID"`
	Auth0ClientSecret string `env:"SCHEDULING_AUTH0_CLIENT_SECRET"`
}

type CassandraConfig struct {
	Host     string `env:"CASS_HOST"`
	Port     string `env:"CASS_PORT"`
	Keyspace string `env:"CASS_KEYSPACE"`
}

type RabbitConfig struct {
	Host     string `env:"RABBIT_HOST"`
	Port     string `env:"RABBIT_PORT"`
	Username string `env:"RABBIT_USERNAME"`
	Password string `env:"RABBIT_PASSWORD"`
	Name     string `env:"RABBIT_QUEUE_NAME"`
}

type Auth0Config struct {
	Domain   string `env:"VITE_AUTH0_DOMAIN"`
	Audience string `env:"VITE_AUTH0_AUDIENCE"`
	JWKSUrl  string `env:"JWKS_URL"`
}

type Auth0WorkerConfig struct {
	ClientId     string `env:"WORKER_AUTH0_CLIENT_ID"`
	ClientSecret string `env:"WORKER_AUTH0_CLIENT_SECRET"`
}

type WorkerConfig struct {
	ID string `env:"WORKER_ID"`
}

type Config struct {
	BaseEndpoint     string `env:"BASE_ENDPOINT"`
	JWTSecretKey     string `env:"JWT_SECRET_KEY"`
	TaskTableName    string `env:"TASK_TABLE_NAME"`
	TaskTableVersion string `env:"TASK_TABLE_VERSION"`
	Env              string `env:"ENV" envDefault:"development"`
	JobAPIEndpoint   string `env:"REPORTER_URL"`
	SandboxCount     int    `env:"SANDBOX_COUNT"`
	TempDir          string `env:"TEMP_DIR"`
	WorkerID         string `env:"WORKER_ID"`
	Auth0            Auth0Config
	Auth0Worker      Auth0WorkerConfig
	Database         DatabaseConfig
	JobService       JobServiceConfig
	Cassandra        CassandraConfig
	Rabbit           RabbitConfig
	Schedule         ScheduleServiceConfig
	Worker           WorkerConfig
}
