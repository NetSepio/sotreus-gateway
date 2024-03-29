package envconfig

import (
	"log"
	"time"

	"github.com/caarlos0/env/v6"
	_ "github.com/joho/godotenv/autoload"
)

type config struct {
	PASETO_PRIVATE_KEY    string        `env:"PASETO_PRIVATE_KEY,required"`
	PASETO_EXPIRATION     time.Duration `env:"PASETO_EXPIRATION,required"`
	APP_NAME              string        `env:"APP_NAME,required"`
	APP_ENVIRONMENT       string        `env:"APP_ENVIRONMENT,required"`
	AUTH_EULA             string        `env:"AUTH_EULA,required"`
	HTTP_PORT             string        `env:"HTTP_PORT,required"`
	GIN_MODE              string        `env:"GIN_MODE,required"`
	DB_HOST               string        `env:"DB_HOST,required"`
	DB_USERNAME           string        `env:"DB_USERNAME,required"`
	DB_PASSWORD           string        `env:"DB_PASSWORD,required"`
	DB_NAME               string        `env:"DB_NAME,required"`
	DB_PORT               int           `env:"DB_PORT,required"`
	ALLOWED_ORIGIN        []string      `env:"ALLOWED_ORIGIN,required" envSeparator:","`
	PASETO_SIGNED_BY      string        `env:"PASETO_SIGNED_BY,required"`
	VPN_DEPLOYER_US02     string        `env:"VPN_DEPLOYER_US02,required"`
	VPN_DEPLOYER_US01     string        `env:"VPN_DEPLOYER_US01,required"`
	VPN_DEPLOYER_EU01     string        `env:"VPN_DEPLOYER_EU01,required"`
	VPN_DEPLOYER_IN01     string        `env:"VPN_DEPLOYER_IN01,required"`
	STRIPE_WEBHOOK_SECRET string        `env:"STRIPE_WEBHOOK_SECRET,required"`
	STRIPE_SECRET_KEY     string        `env:"STRIPE_SECRET_KEY,required"`
}

var EnvVars config = config{}

func InitEnvVars() {

	if err := env.Parse(&EnvVars); err != nil {
		log.Fatalf("failed to parse EnvVars: %s", err)
	}
}
