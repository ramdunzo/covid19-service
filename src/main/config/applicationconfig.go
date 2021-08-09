package config

import (
	"covid-19/src/main/apis"
	"covid-19/src/main/client"
	"covid-19/src/main/dal/repos"
	"covid-19/src/main/service/impl"
	"covid-19/src/main/service/orch"
	"github.com/gorilla/mux"
)

func InitializeApplicationConfig(apiMux *mux.Router) {
	googleClient := client.NewGoogleClient()
	dbConnections := GetDatabase()
	covid19Repo := repos.NewCovid19Repo(dbConnections)
	covid19Service := impl.NewCovid19Impl(covid19Repo)
	covid19Orch := orch.NewCovid19Orch(googleClient, covid19Service)
	apis.NewCovid19Controller(apiMux, covid19Orch)
}
