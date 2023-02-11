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
	op := verrors.Op("service: Find EM by id")

	data, err := service.atbRepository.FindMatchDataByIdInDB(id)
	if err != nil {
		return sheazuzu.MatchData{}, verrors.E(op, err)
	}

	return mapper.MatchDataToBo(data), nil
}
