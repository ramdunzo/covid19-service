package repos

import (
	"context"
	"covid19-service/src/main/dal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Covid19Repo struct {
	collection *mongo.Collection
}

func NewCovid19Repo(dbConnections *mongo.Database) *Covid19Repo {
	return &Covid19Repo{
		collection: dbConnections.Collection("covidcases"),
	}
}

func (repo *Covid19Repo) BulkInsert(covid19 []*models.Covid19) ([]*models.Covid19, error) {
	docs := make([]interface{}, len(covid19))

	for i := 0; i < len(covid19); i++ {
		docs[i] = covid19[i]
	}
	_, err := repo.collection.InsertMany(context.Background(), docs)
	return covid19, err

}

func (repo *Covid19Repo) UpsertByPlace(covid19 *models.Covid19) error {
	opts := options.Update().SetUpsert(true)
	filter := bson.D{{"place_name", covid19.PlaceName}}
	update := bson.M{
		"$set": covid19,
	}
	_, err := repo.collection.UpdateOne(nil, filter, update, opts)

	return err
}

func (repo *Covid19Repo) GetByStates(places []string) ([]models.Covid19, error) {

	var covidData []models.Covid19
	query := bson.M{"place_name": bson.M{"$in": places}}
	res, err := repo.collection.Find(nil, query)

	if err == nil {
		if err = res.All(nil, &covidData); err != nil {
			return covidData, err
		}
	}
	return covidData, nil
}
