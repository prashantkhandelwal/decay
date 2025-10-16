package server

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/contrib/static"
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

	prometheus.MustRegister(middleware.HttpRequestTotal)

	port := c.Server.PORT

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

	// Set trusted proxies
	router.SetTrustedProxies(c.Server.TrustedProxies)
	router.Use(cors.New(cors.Config{
		AllowedOrigins:   []string{"http://localhost:8989"},
		AllowedMethods:   []string{"POST", "GET"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	}))

	embedFS := EmbedFolder(Ui, "ui", true)
	router.Use(static.Serve("/", embedFS))
	log.Printf("Server started on port: %v\n", port)

	// User
	//router.POST("/login", handlers.Login())

	router.POST("/login", handlers.LoginHandler())
	router.POST("/token/refresh", handlers.RefreshHandler())
	router.POST("/logout", handlers.LogoutHandler())
	router.POST("/upload", handlers.UploadHandler(c.File))

	// TODO: Set this to 'api' later and protect with AuthMiddleware
	_ = router.Group("/api", middleware.AuthMiddleware)

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Total HTTP requests metric
	router.Use(func(c *gin.Context) {
		c.Next()
		middleware.HttpRequestTotal.WithLabelValues(c.Request.Method, c.Request.URL.Path, strconv.Itoa(c.Writer.Status())).Inc()
	})

	router.GET("/ping", handlers.Ping)

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"code": "PAGE_NOT_FOUND", "message": "Page not found",
		})
	})

	err := router.Run(":" + port)
	if err != nil {
		log.Fatalf("Error starting the server! - %v", err)
	}

	log.Println("Server running!")
}
