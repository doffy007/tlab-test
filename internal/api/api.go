package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	db "github.com/doffy007/tlab-test/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype/zeronull"
	"github.com/rs/zerolog/log"
)

const (
	UserID = "userId"
	CredID = "credId"
	Role   = "role"
	Token  = "token"
)

var (
	ErrPayload           = errors.New("something wrong with payload format")
	ErrQuery             = errors.New("something wrong with query format")
	ErrUnauthorized      = errors.New("please signin")
	ErrUnverified        = errors.New("please verify your account")
	ErrInsufficientRoles = errors.New("insufficient roles")
	ErrNotFound          = errors.New("not found")
	ErrNoAccess          = errors.New("has no access")
	ErrUpstream          = errors.New("upstream error")
	ErrNotImplemented    = errors.New("not implemented")
)

type Response struct {
	Results any   `json:"results"`
	Total   int64 `json:"total"`
}

func (resp *Response) MarshalJSON() ([]byte, error) {
	type Alias Response
	if resp.Total == 0 {
		resp.Results = []string{}
	}

	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(resp),
	})
}

func Abort(ctx *gin.Context, code int, err, respErr error) {
	log.Error().
		Err(err).
		Str("method", ctx.Request.Method).
		Str("path", ctx.FullPath()).
		Msg("")

	if respErr == nil {
		respErr = err
	}

	ctx.AbortWithStatusJSON(code, gin.H{"error": respErr.Error()})
}

func HasActiveSession(ctx *gin.Context) bool {
	return ctx.GetString(UserID) != "0"
}

func GetUserID(ctx *gin.Context, abort bool) (userID uint64, err error) {
	userID, err = strconv.ParseUint(ctx.GetString(UserID), 10, 64)
	if err != nil {
		log.Debug().Err(err).Msg("invalid user id")
		if abort {
			Abort(ctx, http.StatusBadRequest, fmt.Errorf("invalid user id: %s", ctx.GetString(UserID)), nil)
		}
		return
	}
	return
}

func GetStringUserID(ctx *gin.Context, abort bool) (userID string, err error) {
	userID = ctx.GetString(UserID)
	if userID == "" {
		err = fmt.Errorf("invalid user id: %s", userID)
		log.Debug().Err(err).Msg("invalid user id")
		if abort {
			Abort(ctx, http.StatusBadRequest, err, nil)
		}
		return
	}
	return
}

func GetInt64Param(ctx *gin.Context, k string, abort bool) (v int64, err error) {
	v, err = strconv.ParseInt(ctx.Param(k), 10, 64)
	if err != nil {
		if abort {
			Abort(ctx, http.StatusBadRequest, fmt.Errorf("invalid numeric param %s: %s", k, ctx.Param(k)), nil)
		}
		return
	}
	return
}

func GetIntParam(ctx *gin.Context, k string, abort bool) (v int, err error) {
	v, err = strconv.Atoi(ctx.Param(k))
	if err != nil {
		if abort {
			Abort(ctx, http.StatusBadRequest, fmt.Errorf("invalid numeric param %s: %s", k, ctx.Param(k)), nil)
		}
		return
	}
	return
}

func GetUint64(ctx *gin.Context, k string) (v uint64, err error) {
	s := ctx.GetString(k)
	if s != "" {
		v, err = strconv.ParseUint(ctx.GetString(k), 10, 64)
		if err != nil {
			return
		}
		return
	}

	return 0, nil
}

func GetUint64Param(ctx *gin.Context, k string, abort bool) (v uint64, err error) {
	v, err = strconv.ParseUint(ctx.Param(k), 10, 64)
	if err != nil {
		if abort {
			Abort(ctx, http.StatusBadRequest, fmt.Errorf("invalid numeric param %s: %s", k, ctx.Param(k)), nil)
		}
		return
	}
	return
}

func GetUint64Query(ctx *gin.Context, k string, abort bool) (v uint64, err error) {
	if ctx.Query(k) != "" {
		v, err = strconv.ParseUint(ctx.Query(k), 10, 64)
		if err != nil {
			if abort {
				Abort(ctx, http.StatusBadRequest, fmt.Errorf("invalid numeric query %s: %s", k, ctx.Query(k)), nil)
			}
			return
		}
	}

	return
}

func GetBoolQuery(ctx *gin.Context, k string, abort bool) (v bool, err error) {
	if ctx.Query(k) != "" {
		v, err = strconv.ParseBool(ctx.Query(k))
		if err != nil {
			if abort {
				Abort(ctx, http.StatusBadRequest, fmt.Errorf("invalid bool query %s: %s", k, ctx.Query(k)), nil)
			}
			return
		}
	}

	return
}

func GetStringParam(ctx *gin.Context, k string, abort bool) (v string, err error) {
	v = ctx.Param(k)
	if v == "" {
		if abort {
			Abort(ctx, http.StatusBadRequest, fmt.Errorf("invalid string param %s: %s", k, ctx.Param(k)), nil)
		}
		return
	}
	return
}

func trackUser(userID, ip, ua string) error {
	return db.Service.Commit(nil, func(tx pgx.Tx) error {
		_, err := tx.Exec(
			context.Background(),
			`UPDATE "user"
			SET 
				last_seen = $2,
				ip_address = $3,
				user_agent = $4
			WHERE id = $1`,
			userID,
			time.Now(),
			zeronull.Text(ip),
			zeronull.Text(ua),
		)
		return err
	})
}

func SanitizeXFFHeader(ctx *gin.Context) {
	ips := strings.Split(ctx.Request.Header.Get("X-Forwarded-For"), ",")
	if len(ips) > 3 {
		ips = ips[len(ips)-3:]
		ctx.Request.Header.Set("X-Forwarded-For", strings.Join(ips, ","))
	}
}
