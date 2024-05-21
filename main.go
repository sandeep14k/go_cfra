package main

import (
	"log"
	"os"

	routes "gin/routes"
	"gin/worker"

	"gin/store"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	mongoStore := new(store.MongoStore)
	mongoStore.OpenConnectionWithMongoDB()
	var wg sync.WaitGroup

	wg.Add(2)

	go worker.PerformWork(mongoStore, &wg)
	go server1(&wg)
	wg.Wait()
}
func server1(wg *sync.WaitGroup) {
	defer wg.Done()
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	router := gin.New()
	router.Use(gin.Logger())

	routes.AuthRoutes(router)
	routes.UserRoutes(router)

	router.GET("/api-1", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "Access granted for api-1"})
	})

	router.GET("/api-2", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "Access granted for api-2"})
	})

	router.Run(":" + port)
}
