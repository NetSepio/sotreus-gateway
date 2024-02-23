package envconfig

import (
	"log"
	"time"

	"github.com/caarlos0/env/v6"
	_ "github.com/joho/godotenv/autoload"
)

type config struct {
	PASETO_PRIVATE_KEY string        `env:"PASETO_PRIVATE_KEY,required"`
	PASETO_EXPIRATION  time.Duration `env:"PASETO_EXPIRATION,required"`
	APP_PORT           int           `env:"APP_PORT,required"`
	AUTH_EULA          string        `env:"AUTH_EULA,required"`
	APP_NAME           string        `env:"APP_NAME,required"`
	GIN_MODE           string        `env:"GIN_MODE,required"`
	DB_HOST            string        `env:"DB_HOST,required"`
	DB_USERNAME        string        `env:"DB_USERNAME,required"`
	DB_PASSWORD        string        `env:"DB_PASSWORD,required"`
	DB_NAME            string        `env:"DB_NAME,required"`
	DB_PORT            int           `env:"DB_PORT,required"`
	ALLOWED_ORIGIN     []string      `env:"ALLOWED_ORIGIN,required" envSeparator:","`
	PASETO_SIGNED_BY   string        `env:"PASETO_SIGNED_BY,required"`
	VPN_DEPLOYER_US02  string        `env:"VPN_DEPLOYER_US02,required"`
}

var EnvVars config = config{}

func InitEnvVars() {

	if err := env.Parse(&EnvVars); err != nil {
		log.Fatalf("failed to parse EnvVars: %s", err)
	}
}
