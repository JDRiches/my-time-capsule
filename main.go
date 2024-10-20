package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type MessageCapsule struct {
	IdToken string `json:"idToken"`
	Message string `json:"message"`
	Unlock  string `json:"unlock"`
}

func PostMessage(c *gin.Context, client firestore.Client, authClient auth.Client, ctx context.Context) {

	var message MessageCapsule

	if err := c.BindJSON(&message); err != nil {
		return
	}

	authToken, err := authClient.VerifyIDToken(context.Background(), message.IdToken)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	unlockTime, _ := time.Parse("02-01-2006 15:04:05", message.Unlock)

	_, _, postErr := client.Collection("messages").Add(ctx, map[string]interface{}{
		"user":     authToken.UID,
		"message":  message.Message,
		"created":  time.Now(),
		"unlocked": unlockTime,
	})
	if postErr != nil {
		log.Fatalf("Failed adding message: %v", err)
	}

}

func main() {

	// Env File Stuff
	envFile, _ := godotenv.Read(".env")
	projectID := envFile["project"]

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

	router.Use(AuthMiddleware(*authClient))
	router.Use(FirestoreMiddleware(*client))
	router.Use(CtxMiddleware(ctx))

	router.POST("/create", func(c *gin.Context) {

		authClient, ok := c.MustGet("authConn").(auth.Client)
		if !ok {
			//handle error
		}

		firestoreClient, ok := c.MustGet("firestoreConn").(firestore.Client)
		if !ok {
			//handle error
		}

		ctx, ok := c.MustGet("ctx").(context.Context)
		if !ok {
			//handle error
		}

		PostMessage(c, firestoreClient, authClient, ctx)
	})

	router.Run("localhost:8080")

}
