package main

import (
	"context"
	"fmt"
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
	Message string `json:"message"`
	Unlock  string `json:"unlock"`
}

func PostMessage(c *gin.Context, client firestore.Client, authClient auth.Client, ctx context.Context) {

	var message MessageCapsule

	if err := c.BindJSON(&message); err != nil {
		return
	}

	fmt.Println(c.Request.Header)

	authToken, err := authClient.VerifyIDToken(context.Background(), c.Request.Header["Token"][0])
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

func GetMessageSummaries(c *gin.Context, client firestore.Client, authClient auth.Client, ctx context.Context) {

	authToken, err := authClient.VerifyIDToken(context.Background(), c.Request.Header["Token"][0])
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	uid := authToken.UID

	query := client.Collection("messages").Where("user", "==", uid)

	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		log.Fatalf("Failed to retrieve documents: %v", err)
	}

	for _, doc := range docs {
		fmt.Printf("Document ID: %s\nData: %v\n", doc.Ref.ID, doc.Data())
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

	router.GET("/messages", func(c *gin.Context) {
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

		GetMessageSummaries(c, firestoreClient, authClient, ctx)
	})

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
