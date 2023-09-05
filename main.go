package main

import (
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math/rand"
	"os"
	"time"
)

func main() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	rand.Seed(time.Now().UnixNano())
	r := gin.Default()

	// Enable CORS
	r.Use(cors.Default())

	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
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

	// Change the collection to 'recipes'
	collection := client.Database("test").Collection("recipes")

	r.GET("/recipe", func(c *gin.Context) {
		reqCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		cursor, err := collection.Find(reqCtx, bson.M{})
		if err != nil {
			fmt.Println(err)
			return
		}
		var recipes []bson.M
		err = cursor.All(reqCtx, &recipes)
		if err != nil {
			fmt.Println(err)
			return
		}

		if len(recipes) > 0 {
			randomIndex := rand.Intn(len(recipes))
			randomRecipe := recipes[randomIndex]

			c.JSON(200, gin.H{
				"recipe": randomRecipe,
			})
		} else {
			c.JSON(200, gin.H{
				"message": "No recipes in the database",
			})
		}
	})

	r.Run()
}
