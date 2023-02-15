package repository

import (
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
	"sheazuzu/sheazuzu/src/entity"
)

type SheazuzuRepository struct {
	DB     *gorm.DB
	logger *zap.SugaredLogger
}

func ProvideSheazuzuRepository(DB *gorm.DB, logger *zap.SugaredLogger) *SheazuzuRepository {
	return &SheazuzuRepository{
		DB:     DB,
		logger: logger,
	}
}

func (repository *SheazuzuRepository) FindMatchDataByIdInDB(id int) (entity.MatchData, error) {

	var data entity.MatchData

	db := repository.DB.Preload("AdditionalInformation").Where("id = ?", id).Find(&data)
	if db.Error != nil {
		return entity.MatchData{}, db.Error
	}

	return data, nil
}

func (repository *SheazuzuRepository) UpdateMatchDataInDB(data entity.MatchData) (string, int, error) {

	db := repository.DB.Create(&data)
	if db.Error != nil {
		return "failed", 0, db.Error
	}
	return "successful!", data.Id, nil
}
