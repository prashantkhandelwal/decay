package server

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prashantkhandelwal/decay/config"
	"github.com/prashantkhandelwal/decay/server/handlers"
	"github.com/prashantkhandelwal/decay/server/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Run(c *config.Config) {

	if _, err := os.Stat("data"); err != nil {
		log.Println("\"data\" directory not found!....Creating")
		if err := os.MkdirAll("data", os.ModePerm); err != nil {
			log.Fatalf("ERROR: Cannot create \"data\" directory - %v", err.Error())
			panic(err)
		}
	}

	// Register Prometheus metrics
	prometheus.MustRegister(middleware.HttpRequestTotal)
	prometheus.MustRegister(middleware.TotalFileUploadRequests)
	prometheus.MustRegister(middleware.SuccessfulFileUploadRequests)
	prometheus.MustRegister(middleware.FailedFileUploadRequests)

	port := c.Server.PORT
	var debugPort uint16
	if c.Debugging.EnablePprof {
		debugPort = c.Debugging.Port
		log.Printf("Debugging enabled on port: %v\n", debugPort)
	} else {
		log.Println("Debugging is disabled.")
	}

	if c.Environment != "" {
		if strings.ToLower(c.Environment) == "release" {
			log.Printf("Using environment: %v\n", c.Environment)
			gin.SetMode(gin.ReleaseMode)
		} else {
			gin.SetMode(gin.DebugMode)
		}
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.Default()

	// Router for pprof
	_debugRouter := gin.Default()

	pprof.Register(_debugRouter)

	// Set trusted proxies
	router.SetTrustedProxies(c.Server.TrustedProxies)
	router.Use(cors.New(cors.Config{
		AllowedOrigins:   []string{"http://localhost:8989"},
		AllowedMethods:   []string{"POST", "GET"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	}))

	// embedFS := EmbedFolder(Ui, "ui", true)
	// router.Use(static.Serve("/", embedFS))

	// User
	//router.POST("/login", handlers.Login())

	router.POST("/login", handlers.LoginHandler())
	router.POST("/token/refresh", handlers.RefreshHandler())
	router.POST("/logout", handlers.LogoutHandler())

	// This route is used for health checks
	router.GET("/ping", handlers.Ping)

	api := router.Group("/api", middleware.AuthMiddleware)
	api.POST("/upload", handlers.UploadHandler(c.File))

	// Prometheus metrics endpoint
	_debugRouter.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Total HTTP requests metric
	router.Use(func(c *gin.Context) {
		c.Next()
		middleware.HttpRequestTotal.WithLabelValues(c.Request.Method, c.Request.URL.Path, strconv.Itoa(c.Writer.Status())).Inc()
	})

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"code": "PAGE_NOT_FOUND", "message": "Page not found",
		})
	})

	// Running API server
	go func() {
		err := router.Run(":" + strconv.Itoa(int(port)))
		if err != nil {
			log.Fatalf("Error starting the server! - %v", err)
		}

		log.Printf("API server started on port: %v\n", port)

	}()

	// Running pprof server
	go func() {
		err := _debugRouter.Run(":" + strconv.Itoa(int(debugPort)))
		if err != nil {
			log.Fatalf("Error starting the pprof server! - %v", err)
		}

		log.Printf("Pprof server (endpoint /metrics) started on port: %v\n", debugPort)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
