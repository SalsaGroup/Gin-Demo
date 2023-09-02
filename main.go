package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math/rand"
	"os"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	r := gin.Default()

	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGO_URI")))	if err != nil {
		fmt.Println(err)
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Disconnect(ctx)

	collection := client.Database("test").Collection("phrases")

	r.GET("/", func(c *gin.Context) {
		reqCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		cursor, err := collection.Find(reqCtx, bson.M{})
		if err != nil {
			fmt.Println(err)
			return
		}
		var phrases []bson.M
		err = cursor.All(reqCtx, &phrases)
		if err != nil {
			fmt.Println(err)
			return
		}

		if len(phrases) > 0 {
			randomIndex := rand.Intn(len(phrases))
			randomPhrase := phrases[randomIndex]

			c.JSON(200, gin.H{
				"message": randomPhrase["phrase"],
			})
		} else {
			c.JSON(200, gin.H{
				"message": "No phrases in the database",
			})
		}
	})

	r.Run()
}
