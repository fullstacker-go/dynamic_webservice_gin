package main

import (
	"context"
	"log"
	"time"

	"github.com/fullstacker-go/dynamic_webservice_gin/model"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func main() {
	r := gin.Default()

	// Create a Client and execute a ListDatabases operation.
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	// databases, err := client.ListDatabaseNames(ctx, bson.M{})

	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(databases)
	r.POST("/stats", func(c *gin.Context) {
		var websites []model.Webstats

		c.Bind(&websites)
		size := make(chan int)
		for i, _ := range websites {
			url := "https://www." + websites[i].Domain
			go model.ResponseSize(url, size)
			websites[i].Domain = url
			websites[i].ResponseSize = <-size
		}
		collection := client.Database("webstats").Collection("performance")
		ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
		for _, website := range websites {
			result, _ := collection.InsertOne(ctx, website)
			c.String(200, "Inserted id's is %s", result.InsertedID)
		}

		// if err != nil {
		// 	log.Fatal(err)
		// }

		//c.JSON(200, &websites)

	})
	r.GET("/getstats", func(c *gin.Context) {
		var webstats []model.Webstats
		collection := client.Database("webstats").Collection("performance")
		ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

		results, _ := collection.Find(ctx, bson.M{})

		err = results.All(ctx, &webstats)
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(200, webstats)

	})
	r.Run(":3000")
}
