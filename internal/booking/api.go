package booking

import (
	"errors"
	"net/http"

	"github.com/doffy007/tlab-test/internal/api"
	"github.com/gin-gonic/gin"
)

func CreateBooking(ctx *gin.Context) {
	req := &Booking{}

	err := ctx.ShouldBind(&req)
	if err != nil {
		api.Abort(ctx, http.StatusBadRequest, err, api.ErrPayload)
		return
	}

	booking, err := Service.CreateBooking(req)
	if err != nil {
		api.Abort(ctx, http.StatusInternalServerError, err, errors.New("failed to create event"))
		return
	}

	ctx.JSON(http.StatusOK, booking)
}
