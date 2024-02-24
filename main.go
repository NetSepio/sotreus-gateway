package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/NetSepio/sotreus-gateway/api"
	"github.com/NetSepio/sotreus-gateway/config/constants"
	"github.com/NetSepio/sotreus-gateway/config/dbconfig"
	"github.com/NetSepio/sotreus-gateway/config/envconfig"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	ginApp := gin.Default()
	envconfig.InitEnvVars()
	dbconfig.DbInit()
	constants.InitConstants()
	// cors middleware
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Authorization", "Content-Type"}
	ginApp.Use(cors.New(config))

	ginApp.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"status": 404, "message": "Invalid Endpoint Request"})
	})
	api.ApplyRoutes(ginApp)
	ginApp.Run(":" + os.Getenv("HTTP_PORT"))

	// wait for a SIGINT or SIGTERM signal
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	fmt.Println("Received signal, shutting down...")
}
