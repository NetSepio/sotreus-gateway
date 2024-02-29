package regions

import "github.com/NetSepio/sotreus-gateway/config/envconfig"

type Region struct {
	Name       string
	Code       string
	ServerHttp string
}

var Regions map[string]Region

//todo
//Store regions in persistent DB
//

func InitRegions() {
	if envconfig.EnvVars.APP_ENVIRONMENT == "dev" {
		Regions = map[string]Region{
			"us02": {
				Name:       "US 2",
				Code:       "us02",
				ServerHttp: envconfig.EnvVars.VPN_DEPLOYER_US02,
			},
		}
	} else {
		Regions = map[string]Region{
			"us01": {
				Name:       "US 01",
				Code:       "us01",
				ServerHttp: envconfig.EnvVars.VPN_DEPLOYER_US01,
			},
			"eu01": {
				Name:       "EU 01",
				Code:       "eu01",
				ServerHttp: envconfig.EnvVars.VPN_DEPLOYER_EU01,
			},
			"in01": {
				Name:       "IN 01",
				Code:       "in01",
				ServerHttp: envconfig.EnvVars.VPN_DEPLOYER_IN01,
			},
		}

	}
}
