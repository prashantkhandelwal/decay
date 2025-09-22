package server

import (
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/prashantkhandelwal/decay/config"
	"github.com/prashantkhandelwal/decay/server/handlers"
)

func Run(c *config.Config) {

	if _, err := os.Stat("data"); err != nil {
		log.Println("\"data\" directory not found!....Creating")
		if err := os.MkdirAll("data", os.ModePerm); err != nil {
			log.Fatalf("ERROR: Cannot create \"data\" directory - %v", err.Error())
			panic(err)
		}
	}

	// err := config.InitDB()
	// if err != nil {
	// 	log.Fatalf("ERROR: Unable to configure database - %v", err.Error())
	// 	panic(err)
	// }

	port := c.Port

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

	embedFS := EmbedFolder(Ui, "ui", true)
	router.Use(static.Serve("/", embedFS))
	log.Printf("Server started on port: %v\n", port)

	// User
	//router.POST("/login", handlers.Login())

	router.GET("/ping", handlers.Ping)
	//router.POST("/save", handlers.SaveBookmark())

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
