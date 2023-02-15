package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	verrors "sheazuzu/common/src/errors"
	"sheazuzu/common/src/logging"
	"sheazuzu/common/src/tracing"
	"sheazuzu/common/src/utils"
	"sheazuzu/sheazuzu/src/generated/sheazuzu"
)

type sheazuzuService interface {
	FindMatchDataById(int) (sheazuzu.MatchData, error)
	UpdateMatchData(data sheazuzu.MatchData) (string, int, error)
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

func (controller *Controller) UploadMatchDataUsingPOST(w http.ResponseWriter, r *http.Request) {
	op := verrors.Op("controller: GetFindMachine")

	ctx := r.Context()

	var requestBody sheazuzu.MatchData
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		msg := "Invalid request body"
		handleError(ctx, w, verrors.E(op, verrors.HttpBadRequest, err, msg), msg)
		return
	}

	msg, id, err := controller.service.UpdateMatchData(requestBody)
	if err != nil {
		writeErrorResponse(w, op, err, "error updating MatchData", controller.logger)
		return
	}

	_ = json.NewEncoder(w).Encode(sheazuzu.UpdateResponse{
		MatchID: utils.ToIntPtr(id),
		Message: utils.ToStringPtr(msg),
	})
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

func handleError(ctx context.Context, writer http.ResponseWriter, err error, message string) {

	statusCode := verrors.HttpErrorCodeFromError(err)
	writer.WriteHeader(statusCode)
	errorCode := verrors.GetErrorCode(err)

	if statusCode == http.StatusNoContent {
		return // NoContent 204 doesn't have a response body
	}

	logging.ContextLogger(ctx).Errorw(err.Error(), "statusCode", statusCode, "errorCode", errorCode)

	_ = json.NewEncoder(writer).Encode(
		sheazuzu.ErrorResponse{
			Code:    utils.ToInt32Ptr(errorCode),
			Name:    utils.ToStringPtr(http.StatusText(statusCode)),
			Message: utils.ToStringPtrOrNil(message),
			Details: &[]string{
				fmt.Sprintf("TraceID: %s", tracing.TraceId(ctx)),
			},
		},
	)
}
