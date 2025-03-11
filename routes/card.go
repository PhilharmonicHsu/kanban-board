package routes

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"kanban-board/models"
	"log"
	"time"
)

func RegisterCardRoutes(router *gin.Engine, cardsCollection *mongo.Collection, loc *time.Location) {
	cardGroup := router.Group("/card")
	{
		cardGroup.GET("/:id", func(context *gin.Context) {
			var uri struct {
				ID int64 `uri:"id" binding:"required"`
			}

			if err := context.ShouldBindUri(&uri); err != nil {
				context.JSON(400, gin.H{
					"message": err.Error(),
				})

				return
			}

			requestContext := context.Request.Context()
			pipeline := mongo.Pipeline{
				{
					{"$match", bson.D{
						{"_id", uri.ID},
					}},
				},
				{
					{"$lookup", bson.D{
						{"from", "user"},
						{"localField", "member_ids"},
						{"foreignField", "_id"},
						{"as", "members"},
					}},
				},
				{
					{"$lookup", bson.D{
						{"from", "field"},
						{"localField", "field_ids"},
						{"foreignField", "_id"},
						{"as", "fields"},
					}},
				},
				{
					{"$project", bson.D{
						{"member_ids", 0},
						{"field_ids", 0},
					}},
				},
			}

			cursor, err := cardsCollection.Aggregate(requestContext, pipeline)
			if err != nil {
				log.Printf("關閉 cursor 時發生錯誤: %v", err)
			}

			var results []models.RoughCard

			if err = cursor.All(requestContext, &results); err != nil {
				context.JSON(500, gin.H{"error": err.Error()})
			}

			if len(results) == 0 {
				context.JSON(404, gin.H{"message": "找不到文件"})
			}

			context.JSON(200, results[0])
		})
	}
}
