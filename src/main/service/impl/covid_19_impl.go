package impl

import (
	"covid-19/src/main/dal/models"
	"covid-19/src/main/dal/repos"
	"covid-19/src/main/dtos"
	"covid-19/src/main/service"
	"github.com/dunzoit/froyo/logger"
	"github.com/google/uuid"
	"strings"
)

type Covid19Impl struct {
	Covid19Repo *repos.Covid19Repo
}

func NewCovid19Impl(Covid19Repo *repos.Covid19Repo) service.Covid19Service {
	return &Covid19Impl{Covid19Repo: Covid19Repo}
}

func (service *Covid19Impl) UpdateCovid19Case(covid19Case *dtos.CovidCase) error {
	covid19DataState := covid19Case.Data.Regional
	//var data []*models.Covid19
	for _, state := range covid19DataState {
		id := strings.Replace(uuid.New().String(), "-", "", -1)
		covid19CaseModel := models.NewCovid19(id, state.StateName, covid19Case.LastOriginUpdate, state.Deaths, state.Recovered, state.TotalConfirmed)
		//data = append(data, covid19CaseModel)
		err := service.Covid19Repo.UpsertByPlace(covid19CaseModel)
		logger.Error("[Covid19Impl] UpdateCovid19Case error : ", err)
		return err
	}
	covid19CaseIndia := covid19Case.Data.Summary
	id := strings.Replace(uuid.New().String(), "-", "", -1)
	covid19CaseModel := models.NewCovid19(id, "India", covid19Case.LastOriginUpdate, covid19CaseIndia.Deaths, covid19CaseIndia.Recovered, covid19CaseIndia.Total)
	err := service.Covid19Repo.UpsertByPlace(covid19CaseModel)
	//data = append(data, covid19CaseModel)
	//_, err := service.Covid19Repo.BulkInsert(data)
	return err
}

func (service *Covid19Impl) GetCovidCaseByPlace(state string) (interface{}, error) {
	var response []dtos.GetCovid19CaseByPlaceResponse
	states := []string{"India"}
	states = append(states, state)
	covid19Cases, err := service.Covid19Repo.GetByStates(states)
	if err == nil {
		for _, covid19Case := range covid19Cases {
			response = append(response, dtos.GetCovid19CaseByPlaceResponse{
				Place:              covid19Case.PlaceName,
				Deaths:             covid19Case.Deaths,
				Recovered:          covid19Case.Recovered,
				ConfirmedCovidCase: covid19Case.ConfirmedCovidCase,
				UpdatedAt:          covid19Case.UpdatedAt,
			})
		}
		return response, nil
	}
	return nil, err
}
