package routes

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"kanban-board/models"
	"log"
	"time"
)

func RegisterUserRoutes(router *gin.Engine, usersCollection *mongo.Collection, loc *time.Location) {
	userGroup := router.Group("/user")
	{
		userGroup.GET("", func(c *gin.Context) {
			reqCtx := c.Request.Context()
			var results []models.UserResponse

			cursor, err := usersCollection.Find(reqCtx, bson.D{})
			if err != nil {
				c.JSON(500, gin.H{"message": err.Error()})

				return
			}

			defer func() {
				if err := cursor.Close(reqCtx); err != nil {
					log.Printf("關閉 cursor 時發生錯誤: %v", err)
				}
			}()

			for cursor.Next(reqCtx) {
				var user models.User
				if err := cursor.Decode(&user); err != nil {
					c.JSON(500, gin.H{
						"message": "解碼文件失敗",
					})

					return
				}

				userResponse := models.UserResponse{
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
	}
}
