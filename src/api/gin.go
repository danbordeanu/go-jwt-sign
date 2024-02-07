package api

import (
	"context"
	sharedMiddleware "dev.azure.com/coderollers/almeria/go-shared-noversion/http/middleware"
	"fmt"
	"github.com/danbordeanu/go-logger"
	"github.com/danbordeanu/go-stats/concurrency"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"jwt-sign/api/handlers"
	"jwt-sign/configuration"
	"net/http"
	"time"
)

const httpServerShutdownGracePeriodSeconds = 20

func StartGin(ctx context.Context) {
	defer concurrency.GlobalWaitGroup.Done()

	conf := configuration.AppConfig()
	log := logger.SugaredLogger()

	// Set up gin
	log.Debugf("Setting up Gin")
	if !conf.GinLogger {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()

	// Set up the middleware
	if conf.GinLogger {
		log.Warnf("Gin's logger is active! Logs will be unstructured!")
		router.Use(gin.Logger())
	}

	router.Use(gin.Recovery())
	router.Use(sharedMiddleware.CorrelationId())

	// TODO: We can move CORS to Ingress
	if conf.CorsAllowOrigins != "Disabled" {
		router.Use(cors.New(cors.Config{
			AllowOrigins: []string{conf.CorsAllowOrigins},
			AllowMethods: []string{"POST", "HEAD", "PATCH", "OPTIONS", "GET", "PUT"},
			AllowHeaders: []string{"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token",
				"Authorization", "accept", "origin", "Cache-Control", "X-Requested-With"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}))
	}
	router.Use(otelgin.Middleware("jwt-sign"))

	// let's load the html crap
	router.Static("/assets", "./assets") 
	router.LoadHTMLGlob("templates/**")

	// Set up the groups
	userAPI := router.Group("/v1")
	{

		// validate
		userAPI.POST("/validate-jwt", handlers.ValidateJwt)

		// signature validate
		userAPI.POST("/verify-signature", handlers.VerifySignature)

	}

	// Activate swagger if configured
	if conf.UseSwagger {
		log.Infof("Swagger is active, enabling endpoints")
		url := ginSwagger.URL("/swagger/doc.json") // The url pointing to API definition
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	}

	// Set up the listener
	httpSrv := &http.Server{
		Addr:    fmt.Sprintf(":%d", conf.HttpPort),
		Handler: router,
	}

	// Start the HTTP Server
	go func() {
		log.Infof("Listening on port %d", conf.HttpPort)
		if err := httpSrv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatalf("Unrecoverable HTTP Server failure: %s", err.Error())
			}
		}
	}()

	// Block until SIGTERM/SIGINT
	<-ctx.Done()

	// Clean up and shutdown the HTTP server
	cleanCtx, cancel := context.WithTimeout(context.Background(), httpServerShutdownGracePeriodSeconds*time.Second)
	defer cancel()
	log.Infof("Attempting to shutdown the HTTP server with a timeout of %d seconds", httpServerShutdownGracePeriodSeconds)
	if err := httpSrv.Shutdown(cleanCtx); err != nil {
		log.Errorf("HTTP server failed to shutdown gracefully: %s", err.Error())
	} else {
		log.Infof("HTTP Server was shutdown successfully")
	}
}
