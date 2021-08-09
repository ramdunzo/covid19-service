package models

type Covid19 struct {
	Id                 string `bson:"_id"`
	PlaceName          string `bson:"place_name"`
	Deaths             int64  `bson:"deaths"`
	Recovered          int64  `bson:"discharged"`
	ConfirmedCovidCase int64  `bson:"covid_case"`
	UpdatedAt          string `bson:"updated_at"`
}

func NewCovid19(id, placeName, updatedAt string, deaths, recovered, confirmedCovidCase int64) *Covid19 {
	return &Covid19{
		Id:                 id,
		PlaceName:          placeName,
		Deaths:             deaths,
		Recovered:          recovered,
		ConfirmedCovidCase: confirmedCovidCase,
		UpdatedAt:          updatedAt,
	}
}
