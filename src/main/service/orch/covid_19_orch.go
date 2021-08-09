package orch

import (
	"covid-19/src/main/client"
	"covid-19/src/main/dtos"
	"covid-19/src/main/service"
	"strings"
)

type Covid19Orch struct {
	GoogleClient   client.GoogleClient
	Covid19Service service.Covid19Service
}

func NewCovid19Orch(googleClient client.GoogleClient, covid19Service service.Covid19Service) *Covid19Orch {
	return &Covid19Orch{
		GoogleClient:   googleClient,
		Covid19Service: covid19Service,
	}
}
func (orch *Covid19Orch) UpdateCovid19Case() (interface{}, error) {
	covid19Cases, err := orch.GoogleClient.GetIndiaCovidCases()
	if err == nil {
		err = orch.Covid19Service.UpdateCovid19Case(covid19Cases)
	}
	return covid19Cases, err
}

func (orch *Covid19Orch) GetCovid19Data(request dtos.GetCovid19CaseByPlaceRequest) (interface{}, error) {
	placeDetails, err := orch.GoogleClient.GetStateNameFromLatLng(request.Lat, request.Lng, request.ApiKey)
	if err == nil {
		state := placeDetails.Items[0].Address.State
		state = strings.Replace(state, "&", "and", -1)
		return orch.Covid19Service.GetCovidCaseByPlace(state)
	}
	return nil, err
}
