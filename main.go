package main

import (
	"context"
	"log"
	"os"

	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
)

func main() {

	// Env File Stuff
	projectID := os.Getenv("PROJECT_ID")

	// Setting up Firestore database
	ctx := context.Background()
	conf := &firebase.Config{ProjectID: projectID}
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalln(err)
	}

	//Set up auth

	authClient, err := app.Auth(context.Background())
	if err != nil {
		panic(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	// Gin API stuff
	router := gin.Default()

	router.GET("/messages", func(c *gin.Context) {
		GetCapsules(c, *client, *authClient)
	})

	router.POST("/create", func(c *gin.Context) {
		PostCapsule(c, *client, *authClient)
	})

	router.POST("/delete", func(c *gin.Context) {
		DeleteCapsule(c, *client, *authClient)
	})

	router.POST("/open", func(c *gin.Context) {
		OpenCapsule(c, *client, *authClient)
	})

	router.GET("/detail", func(c *gin.Context) {
		GetCapsuleDetail(c, *client, *authClient)
	})

	router.Run()

}
