package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/doffy007/tlab-test/internal/booking"
	"github.com/doffy007/tlab-test/internal/event"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func main() {
	router := gin.New()

	router.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/status"},
	}))

	router.GET("/status", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"success": true})
	})

	v1 := router.Group("/v1")
	{
		v1.POST("event", event.CreateEvent)
		v1.GET("event/:eventId", event.GetEvent)
		v1.POST("/booking", booking.CreateBooking)
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: router,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below.
	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Info().Msgf("listen: %s", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("shutting down server")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("server forced to shutdown")
	}

	log.Info().Msg("server exiting")
}
