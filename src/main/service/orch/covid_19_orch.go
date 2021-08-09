package orch

import (
	"covid19-service/src/main/client"
	"covid19-service/src/main/dtos"
	"covid19-service/src/main/service"
	"covid19-service/src/main/utils"
	"github.com/dunzoit/froyo/logger"
	"github.com/goburrow/cache"
	"strings"
)

type Covid19Orch struct {
	GoogleClient             client.GoogleClient
	Covid19Service           service.Covid19Service
	CacheConfig              *dtos.CacheConfig
	GetCovid19CasesDataCache cache.LoadingCache
}

func NewCovid19Orch(googleClient client.GoogleClient, covid19Service service.Covid19Service, cacheConfig *dtos.CacheConfig) *Covid19Orch {
	orchObj :=  &Covid19Orch{
		GoogleClient:   googleClient,
		Covid19Service: covid19Service,
		CacheConfig:    cacheConfig,
	}
	orchObj.GetCovid19CasesDataCache = utils.GetLoadingCache(cacheConfig, func(key cache.Key) (cache.Value, error) {
		return getCovid19CasesDataCache(orchObj, key)
	})
	return orchObj
}
func (orch *Covid19Orch) UpdateCovid19Case() (interface{}, error) {
	covid19Cases, err := orch.GoogleClient.GetIndiaCovidCases()
	if err == nil {
		err = orch.Covid19Service.UpdateCovid19Case(covid19Cases)
	}
	return covid19Cases, err
}
func getCovid19CasesDataCache(orch *Covid19Orch, key cache.Key) (interface{}, error){
	data := strings.Split(key.(string), "/")
	return orch.GetCovid19Data(data[0], data[1])
}

func (orch *Covid19Orch) GetCovid19DataCache(lat string, lng string) (interface{}, error) {
	cacheKey := lat + "/" + lng
	value, err := orch.GetCovid19CasesDataCache.Get(cacheKey)
	if err != nil {
		logger.Error("Error : ", err)
		return nil, err
	}
	return value, nil
}

func (orch *Covid19Orch) GetCovid19Data(lat string, lng string) (interface{}, error) {
	placeDetails, err := orch.GoogleClient.GetStateNameFromLatLng(lat, lng)
	if err == nil {
		state := placeDetails.Items[0].Address.State
		state = strings.Replace(state, "&", "and", -1)
		return orch.Covid19Service.GetCovidCaseByPlace(state)
	}
	return nil, err
}
