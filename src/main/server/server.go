package server

import (
	"covid19-service/src/main/config"
	"covid19-service/src/main/server/handler"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// Server represents a server mux
type Server struct {
	*mux.Router
	Address string
}

// New setups & returns a server
func New() *Server {
	r := mux.NewRouter()
	//r.Use(mux.CORSMethodMiddleware(r))
	addr :=  "0.0.0.0:8000"
	s := Server{r, addr}
	s.SetupComponents()
	return &s
}

// SetupComponents of the server
func (s Server) SetupComponents() {

	s.HandleFunc(
		"/status", handler.StatusActive,
	).Methods(http.MethodGet)

	apiMux := s.PathPrefix("/api").Subrouter()

	config.InitializeApplicationConfig(apiMux)
}

func (s Server) ServeHTTP() {
	loggedRouter := handlers.LoggingHandler(os.Stdout, s.Router)
	srv := &http.Server{
		Handler: loggedRouter,
		Addr:    s.Address,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: time.Minute,
		ReadTimeout:  time.Minute,
	}

	logrus.Info("Server starting at addr: ", s.Address)
	logrus.Fatal(srv.ListenAndServe())
}