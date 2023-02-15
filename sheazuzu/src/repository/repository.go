package repository

import (
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
	"sheazuzu/sheazuzu/src/database"
	"sheazuzu/sheazuzu/src/entity"
)

type SheazuzuRepository struct {
	DB     *gorm.DB
	Mongo  *database.MongoDatabase
	logger *zap.SugaredLogger
}

func ProvideSheazuzuRepository(DB *gorm.DB, Mongo *database.MongoDatabase, logger *zap.SugaredLogger) *SheazuzuRepository {
	return &SheazuzuRepository{
		DB:     DB,
		Mongo:  Mongo,
		logger: logger,
	}
}

func (repository *SheazuzuRepository) FindMatchDataByIdInDB(id int) (entity.MatchData, error) {

	var data entity.MatchData

	db := repository.DB.Preload("AdditionalInformation").Where("id = ?", id).Find(&data)
	if db.Error != nil {
		return entity.MatchData{}, db.Error
	}

	/*
		mongoEntity, err := repository.Mongo.FindByID(context.Background(), id)
		if err != nil {
			return entity.MatchData{}, db.Error
		}

		fmt.Println("mongoEntity", mongoEntity)


	*/
	return data, nil
}

func (repository *SheazuzuRepository) UpdateMatchDataInDB(data entity.MatchData) (string, int, error) {

	db := repository.DB.Create(&data)
	if db.Error != nil {
		return "failed - mySQL", 0, db.Error
	}

	/*
		err := repository.Mongo.Save(context.Background(), data)
		if err != nil {
			return "failed - Mongo", 0, db.Error
		}

	*/

	return "successful!", data.Id, nil

}
