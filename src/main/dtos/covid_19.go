package dtos

type GetCovid19CaseByPlaceRequest struct {
	Lat    string `json:"lat"`
	Lng    string `json:"lng"`
}
type GetCovid19CaseByPlaceResponse struct {
	Place              string `json:"place"`
	Deaths             int64  `json:"deaths"`
	Recovered          int64  `json:"recovered"`
	ConfirmedCovidCase int64  `json:"covid_case"`
	UpdatedAt          string `json:"last_updated_at"`
}
