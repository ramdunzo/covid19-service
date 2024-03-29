package apis

import (
	"covid19-service/src/main/dtos"
	"covid19-service/src/main/service/orch"
	"encoding/json"
	"github.com/dunzoit/froyo/logger"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

type Covid19Controller struct {
	covid19Orch *orch.Covid19Orch
}

func NewCovid19Controller(apiMux *mux.Router, covid19Orch *orch.Covid19Orch) *Covid19Controller {
	controller := Covid19Controller{covid19Orch: covid19Orch}
	v0mux := apiMux.PathPrefix("/v1/covid").Subrouter()
	v0mux.HandleFunc("/update", controller.UpdateCovidCase).Methods("POST")
	v0mux.HandleFunc("/cases", controller.GetCovidCase).Methods("GET")
	return &controller
}

func (controller *Covid19Controller) UpdateCovidCase(w http.ResponseWriter, r *http.Request) {
	logger.Info("[Covid19Controller] UpdateCovidCase request received *****")
	response, err := controller.covid19Orch.UpdateCovid19Case()
	if err != nil {
		_ = json.NewEncoder(w).Encode(err)
	}
	_ = json.NewEncoder(w).Encode(response)

}

func (controller *Covid19Controller) GetCovidCase(w http.ResponseWriter, r *http.Request) {
	logger.Info("[Covid19Controller] GetCovidCase request received *****")
	reqBody, _ := ioutil.ReadAll(r.Body)
	var request dtos.GetCovid19CaseByPlaceRequest
	_ = json.Unmarshal(reqBody, &request)
	response, err := controller.covid19Orch.GetCovid19DataCache(request.Lat, request.Lng)
	if err != nil {
		_ = json.NewEncoder(w).Encode(err)
	}
	_ = json.NewEncoder(w).Encode(response)
}
