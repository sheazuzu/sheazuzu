package controller

import (
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
	verrors "sheazuzu/common/src/errors"
	"sheazuzu/sheazuzu/src/generated/sheazuzu"
)

type sheazuzuService interface {
	FindMatchDataById(int) (sheazuzu.MatchData, error)
}

type Controller struct {
	service sheazuzuService
	logger  *zap.SugaredLogger
}

func ProvideSheazuzuAPI(service sheazuzuService, logger *zap.SugaredLogger) *Controller {
	return &Controller{
		service: service,
		logger:  logger,
	}
}

func (controller *Controller) GetMatchDataByIdUsingGET(w http.ResponseWriter, r *http.Request, params sheazuzu.GetMatchDataByIdUsingGETParams) {
	op := verrors.Op("controller: GetFindMachine")

	resultList, err := controller.service.FindMatchDataById(params.Id)
	if err != nil {
		writeErrorResponse(w, op, err, "error while getting machine by id", controller.logger)
		return
	}

	_ = json.NewEncoder(w).Encode(sheazuzu.MatchDataResponse{
		MatchData: &resultList,
	})
}

func (controller *Controller) AllMatchDataUsingGET(w http.ResponseWriter, r *http.Request) {
	panic("hello world")
}

func (controller *Controller) UploadUsingPOST(w http.ResponseWriter, r *http.Request) {
	panic("hello world")
}

func writeErrorResponse(writer http.ResponseWriter, op verrors.Op, err error, details string, logger *zap.SugaredLogger) {
	err = verrors.E(op, err)

	statusCode := verrors.HttpErrorCodeFromError(err)
	statusText := http.StatusText(statusCode)
	writer.WriteHeader(statusCode)

	errorCode := verrors.GetErrorCode(err)

	_ = json.NewEncoder(writer).Encode(sheazuzu.ErrorResponse{
		Code:        &errorCode,
		Description: &statusText,
		Details:     &[]string{details},
	})

	logger.Errorw(err.Error(),
		"details", details,
		"statusCode", statusCode,
		"errorCode", errorCode)
}
