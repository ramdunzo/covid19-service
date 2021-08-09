package utils

import (
	"covid19-service/src/main/dtos"
	"github.com/goburrow/cache"
	"time"
)

func GetLoadingCache(cacheConfig *dtos.CacheConfig, loaderFunction cache.LoaderFunc) cache.LoadingCache {

	return cache.NewLoadingCache(loaderFunction,

		cache.WithMaximumSize(cacheConfig.Size),
		cache.WithPolicy(cacheConfig.Policy),
		cache.WithExpireAfterAccess(time.Duration(cacheConfig.Expiry) * time.Second),
		cache.WithRefreshAfterWrite(time.Duration(cacheConfig.Expiry) * time.Second),
	)
}