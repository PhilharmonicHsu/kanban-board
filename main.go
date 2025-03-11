package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"time"

	"kanban-board/db"
	"kanban-board/routes"
)

func main() {
	router := gin.Default()

	db.InitMongo("mongodb://localhost:27017")
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := db.Client.Disconnect(ctx); err != nil {
			log.Printf("關閉連線時發生錯誤: %v", err)
		}
	}()

	database := db.Client.Database("kanban-board")

	loc, _ := time.LoadLocation("America/Denver")

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	routes.RegisterUserRoutes(router, database.Collection("user"), loc)
	routes.RegisterColumnRoutes(router, database)
	routes.RegisterCardRoutes(router, database.Collection("card"), loc)

	func() {
		if err := router.Run(); err != nil {
			log.Printf("執行run時發生錯誤: %v", err)
		}
	}()
}
