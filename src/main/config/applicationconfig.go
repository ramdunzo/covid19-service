package config

import (
	"covid19-service/src/main/apis"
	"covid19-service/src/main/client"
	"covid19-service/src/main/dal/repos"
	"covid19-service/src/main/service/impl"
	"covid19-service/src/main/service/orch"
	"github.com/gorilla/mux"
)

func InitializeApplicationConfig(apiMux *mux.Router) {
	googleClient := client.NewGoogleClient()
	dbConnections := GetDatabase()
	covid19Repo := repos.NewCovid19Repo(dbConnections)
	getCovid19CaseCacheConfig := GetCovid19CaseCacheConfig()
	covid19Service := impl.NewCovid19Impl(covid19Repo)
	covid19Orch := orch.NewCovid19Orch(googleClient, covid19Service, getCovid19CaseCacheConfig)
	apis.NewCovid19Controller(apiMux, covid19Orch)
}
