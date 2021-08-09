package main

import (
	"covid19-service/src/main/config"
	"covid19-service/src/main/server"
)

func main(){
	config.SetUpDatabase()
	srv := server.New()
	srv.ServeHTTP()
}
