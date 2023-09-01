package main

import (
	"github.com/gin-gonic/gin"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		phrases := []string{"Hello world!", "Hi there!", "Greetings!", "Salutations!", "Hey!"}
		randomIndex := rand.Intn(len(phrases))

		c.JSON(200, gin.H{
			"message": phrases[randomIndex],
		})
	})
	r.Run()
}
