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

	repository.DB.Preload("AdditionalInformation").Where("id = ?", id).Find(&data)

	return data, nil
}
