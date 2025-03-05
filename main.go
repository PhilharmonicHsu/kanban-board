package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type User struct {
	ID        int32     `bson:"_id" json:"id"`
	Account   string    `bson:"account" json:"account"`
	Password  string    `bson:"password" json:"password"`
	Name      string    `bson:"name" json:"name"`
	Avatar    *string   `bson:"avatar" json:"avatar"` // 有可能為null時，就使用指標
	CreateAt  time.Time `bson:"created_at" json:"create_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

type UserResponse struct {
	ID        int32   `bson:"_id" json:"id"`
	Account   string  `bson:"account" json:"account"`
	Password  string  `bson:"password" json:"password"`
	Name      string  `bson:"name" json:"name"`
	Avatar    *string `bson:"avatar" json:"avatar"`
	CreateAt  string  `bson:"created_at" json:"create_at"`
	UpdatedAt string  `bson:"updated_at" json:"updated_at"`
}

func main() {
	r := gin.Default()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatal("連線失敗:", err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		log.Fatal("無法連線到 MongoDB:", err)
	}

	database := client.Database("kanban-board")
	usersCollection := database.Collection("user")

	loc, _ := time.LoadLocation("America/Denver") // UTC-7 的時區

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/user", func(c *gin.Context) {
		var results []UserResponse

		cursor, err := usersCollection.Find(ctx, bson.D{})

		if err != nil {
			c.JSON(500, gin.H{
				"message": err,
			})

			return
		}

		defer func() {
			if err := cursor.Close(ctx); err != nil {
				log.Printf("關閉 cursor 時發生錯誤: %v", err)
			}
		}()

		for cursor.Next(ctx) {
			var user User
			if err := cursor.Decode(&user); err != nil {
				c.JSON(500, gin.H{
					"message": "解碼文件失敗",
				})

				return
			}

			userResponse := UserResponse{
				ID:        user.ID,
				Account:   user.Account,
				Password:  user.Password,
				Name:      user.Name,
				Avatar:    user.Avatar,
				CreateAt:  user.CreateAt.In(loc).Format("2006-01-02 15:04:05"),
				UpdatedAt: user.UpdatedAt.In(loc).Format("2006-01-02 15:04:05"),
			}

			results = append(results, userResponse)
		}

		if err := cursor.Err(); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, results)
	})

	func() {
		if err := r.Run(); err != nil {
			log.Printf("執行run時發生錯誤: %v", err)
		}
	}()
}
