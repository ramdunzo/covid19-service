package main

import (
	"covid-19/src/main/config"
	"covid-19/src/main/server"
)

func main(){
	config.SetUpDatabase()
	srv := server.New()
	srv.ServeHTTP()
}
