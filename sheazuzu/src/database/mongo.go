package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	mongoClient "go.mongodb.org/mongo-driver/mongo"
	"go.opencensus.io/trace"
	verrors "sheazuzu/common/src/errors"
	"sheazuzu/common/src/mongo"
	"sheazuzu/common/src/tracing"
	"sheazuzu/sheazuzu/src/entity"
	"sheazuzu/sheazuzu/src/generated/sheazuzu"
)

type MongoDatabase struct {
	mongo *mongo.Database
}

func NewMongoDatabase(mongo *mongo.Database) *MongoDatabase {
	return &MongoDatabase{mongo: mongo}
}

const matchDataSet = "matchData"

func (database *MongoDatabase) Save(ctx context.Context, matchData entity.MatchData) error {
	op := verrors.Op("MongoDB: Save MatchData")

	ctx, span := tracing.StartSpan(ctx, "Save Match Data")
	span.AddAttributes(trace.Int64Attribute("Id", int64(matchData.Id)))
	defer span.End()

	_, err := database.mongo.Database.Collection(matchDataSet).InsertOne(ctx, matchData)

	if err != nil {
		return verrors.E(op, err)
	}

	return nil
}

func (database *MongoDatabase) FindByID(ctx context.Context, id int) (sheazuzu.MatchData, error) {
	op := verrors.Op("MongoDB: Fetch Storage Entry by MarketingCode")
	info := []verrors.Info{
		{
			Name: "id",
			Val:  id,
		},
	}

	ctx, span := tracing.StartSpan(ctx, "Fetch match data by id")
	span.AddAttributes(trace.StringAttribute("id", string(id)))
	defer span.End()

	// option 'i' makes search case-insensitive
	filter := bson.D{
		{
			Key:   "_id",
			Value: id,
		},
	}

	result := database.mongo.Database.Collection(matchDataSet).FindOne(ctx, filter)
	err := result.Err()
	if err != nil {
		if err == mongoClient.ErrNoDocuments {
			return sheazuzu.MatchData{}, nil
		}
		return sheazuzu.MatchData{}, verrors.E(err, op, info)
	}

	var data sheazuzu.MatchData
	err = result.Decode(&data)
	if err != nil {
		return sheazuzu.MatchData{}, verrors.E(op, info, err)
	}

	return data, nil
}
