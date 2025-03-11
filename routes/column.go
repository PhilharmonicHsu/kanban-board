package routes

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"kanban-board/models"
	"log"
)

type Column struct {
	ID          int32  `bson:"_id" json:"id"`
	Name        string `bson:"name" json:"name"`
	Description string `bson:"description" json:"description"`
}

func RegisterColumnRoutes(router *gin.Engine, kanbanBoardDB *mongo.Database) {
	columnsGroup := router.Group("/columns")
	{
		columnsGroup.GET("", func(c *gin.Context) {
			columnsCollection := kanbanBoardDB.Collection("column")
			reqCtx := c.Request.Context()

			pipeline := mongo.Pipeline{
				{{"$sort", bson.D{
					{"sort", 1},
				}}},
			}
			cursor, err := columnsCollection.Aggregate(reqCtx, pipeline)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
			}
			defer func() {
				if err := cursor.Close(reqCtx); err != nil {
					log.Printf("關閉 cursor 時發生錯誤: %v", err)
				}
			}()

			var results []Column
			if err := cursor.All(reqCtx, &results); err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
			}

			c.JSON(200, results)
		})

		columnsGroup.GET("/:columnId/cards", func(context *gin.Context) {
			cardsCollection := kanbanBoardDB.Collection("card")

			var uri struct {
				ColumnId int64 `uri:"columnId" binding:"required"`
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
						{"column_id", uri.ColumnId},
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

			context.JSON(200, results)
		})
	}
}
