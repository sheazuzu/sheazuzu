package service

import (
	"go.uber.org/zap"
	verrors "sheazuzu/common/src/errors"
	"sheazuzu/sheazuzu/src/entity"
	"sheazuzu/sheazuzu/src/generated/sheazuzu"
	"sheazuzu/sheazuzu/src/mapper"
)

type sheazuzuRepository interface {
	FindMatchDataByIdInDB(int) (entity.MatchData, error)
	UpdateMatchDataInDB(data entity.MatchData) (string, int, error)
}

type Service struct {
	atbRepository sheazuzuRepository
	logger        *zap.SugaredLogger
}

func ProvideSheazuzuService(sheazuzuRepository sheazuzuRepository, logger *zap.SugaredLogger) *Service {
	return &Service{
		atbRepository: sheazuzuRepository,
		logger:        logger,
	}
}

func (service *Service) FindMatchDataById(id int) (sheazuzu.MatchData, error) {
	op := verrors.Op("service: Find MatchData by id")

	data, err := service.atbRepository.FindMatchDataByIdInDB(id)
	if err != nil {
		return sheazuzu.MatchData{}, verrors.E(op, err)
	}

	return mapper.MatchDataToBo(data), nil
}

func (service *Service) UpdateMatchData(data sheazuzu.MatchData) (string, int, error) {
	op := verrors.Op("service: Update MatchData")

	msg, id, err := service.atbRepository.UpdateMatchDataInDB(mapper.BoToMatchData(data))
	if err != nil {
		return "", 0, verrors.E(op, err)
	}
	return msg, id, nil
}
