package config

import (
	"covid19-service/src/main/dtos"
)

func GetCovid19CaseCacheConfig() *dtos.CacheConfig {
	return &dtos.CacheConfig{
		Size:1000,
		Policy:"lru",
		Expiry:1800,
	}
}
