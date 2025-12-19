package event

import (
	"errors"
	"net/http"

	"github.com/doffy007/tlab-test/internal/api"
	"github.com/gin-gonic/gin"
)

func CreateEvent(ctx *gin.Context) {
	req := &Event{}

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		api.Abort(ctx, http.StatusBadRequest, err, api.ErrPayload)
		return
	}

	event, err := Service.CreateEvent(req)
	if err != nil {
		api.Abort(ctx, http.StatusInternalServerError, err, errors.New("failed to create event"))
		return
	}

	ctx.JSON(http.StatusOK, event)
}

func GetEvent(ctx *gin.Context) {
	var err error

	eventId, err := api.GetUint64Param(ctx, "eventId", true)
	if err != nil {
		api.Abort(ctx, http.StatusBadRequest, err, errors.New("invalid event ID"))
		return
	}

	resp, err := Service.GetEvent(eventId)
	if err != nil {
		api.Abort(ctx, http.StatusInternalServerError, err, errors.New("failed to get evnet"))
		return
	}

	if resp == nil {
		api.Abort(ctx, http.StatusNotFound, api.ErrNotFound, nil)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
