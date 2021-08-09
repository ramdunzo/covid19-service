package client

import (
	"covid-19/src/main/dtos"
	"encoding/json"
	"github.com/dunzoit/froyo/logger"
	"io/ioutil"
	"net/http"
	"time"
)

type GoogleClient struct {
	client *http.Client
}

func NewGoogleClient() GoogleClient {
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 20,
		},
		Timeout: time.Duration(100) * time.Second,
	}
	return GoogleClient{
		client: client,
	}
}

func (gc *GoogleClient) GetIndiaCovidCases() (*dtos.CovidCase, error) {
	url := "https://api.rootnet.in/covid19-in/stats/latest"
	req, err := http.NewRequest("GET", url, nil)
	if err == nil {
		var res *http.Response
		res, err = gc.client.Do(req)
		if err == nil {
			defer res.Body.Close()
			if res.StatusCode == 200 {
				var body []byte
				body, err = ioutil.ReadAll(res.Body)
				if err == nil {
					var data dtos.CovidCase
					err = json.Unmarshal(body, &data)
					if err == nil {
						return &data, nil
					}
				}
			}
		}
	}
	logger.Info("[GoogleClient] GetIndiaCovidCases Error : ", err)
	return nil, err
}

func (gc *GoogleClient) GetStateNameFromLatLng(lat string, lng string, apiKey string) (*dtos.ReverseGeoCoding, error) {
	url := "https://revgeocode.search.hereapi.com/v1/revgeocode?at=" + lat + "," + lng + "&apikey=" + apiKey
	req, err := http.NewRequest("GET", url, nil)
	if err == nil {
		var res *http.Response
		res, err = gc.client.Do(req)
		if err == nil {
			defer res.Body.Close()
			if res.StatusCode == 200 {
				var body []byte
				body, err = ioutil.ReadAll(res.Body)
				if err == nil {
					var data dtos.ReverseGeoCoding
					err = json.Unmarshal(body, &data)
					if err == nil {
						return &data, nil
					}
				}
			}
		}
	}
	logger.Info("[GoogleClient] GetStateNameFromLatLng  Error : ", err, "for lat : ", lat, "lng : ", lng)
	return nil, err
}
