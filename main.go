package main

import (
	"context"
	"log"
	"my-time-capsule/handler"
	"os"

	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
)

func main() {

	// Load Environment Variables
	projectID := os.Getenv("PROJECT_ID")

	// Setting up Firestore database
	ctx := context.Background()
	conf := &firebase.Config{ProjectID: projectID}
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalln(err)
	}

	// Set up authentication with Firebase
	authClient, err := app.Auth(context.Background())
	if err != nil {
		panic(err)
	}

	// Set up client for Firestore database access
	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	// Set up Gin Router
	router := gin.Default()

	router.GET("/messages", func(c *gin.Context) {
		handler.GetCapsules(c, *client, *authClient)
	})

	router.POST("/create", func(c *gin.Context) {
		handler.PostCapsule(c, *client, *authClient)
	})

	router.POST("/delete", func(c *gin.Context) {
		handler.DeleteCapsule(c, *client, *authClient)
	})

	router.POST("/open", func(c *gin.Context) {
		handler.OpenCapsule(c, *client, *authClient)
	})

	router.GET("/detail", func(c *gin.Context) {
		handler.GetCapsuleDetail(c, *client, *authClient)
	})

	router.Run()

}
