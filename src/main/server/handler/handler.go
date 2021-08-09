package handler

import (
	"encoding/json"
	"net/http"

	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/sirupsen/logrus"
)

// Fields represents JSON
type Fields map[string]interface{}

// Response is written to http.ResponseWriter
type Response struct {
	Code    int
	Payload interface{}
}

var isHealthy = true

// Make creates a http handler from a request handler func
func Make(
	f func(req *http.Request) Response,
) func(w http.ResponseWriter, req *http.Request) {
	handler := func(w http.ResponseWriter, req *http.Request) {

		setupResponse(&w, req)
		if (*req).Method == "OPTIONS" {
			req.Body.Close()
			return
		}

		res := f(req)
		JSON, err := json.Marshal(res.Payload)
		if err != nil {
			logrus.WithError(err).Fatal("json marshal failed")
		}

		w.WriteHeader(res.Code)
		w.Write(JSON)
		req.Body.Close()
	}

	sentryHandler := sentryhttp.New(sentryhttp.Options{})
	handler = sentryHandler.HandleFunc(handler)
	return handler
}

func StatusActive(w http.ResponseWriter, req *http.Request) {
	Make(func(req *http.Request) Response {
		if isHealthy {
			return Response{
				Code:    http.StatusOK,
				Payload: activeHealthCheckResponse,
			}
		}

		return Response{
			Code:    http.StatusInternalServerError,
			Payload: inactiveHealthCheckResponse,
		}
	})(w, req)
}

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Content-Type", "application/json")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
	(*w).Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Origin, Accept, x-requested-with, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
