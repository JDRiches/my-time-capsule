package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
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

	router.Use(AuthMiddleware(*authClient))
	router.Use(FirestoreMiddleware(*client))
	router.Use(CtxMiddleware(ctx))

	router.GET("/messages", func(c *gin.Context) {
		authClient, ok := c.MustGet("authConn").(auth.Client)
		if !ok {
			fmt.Println(ok)
		}

		firestoreClient, ok := c.MustGet("firestoreConn").(firestore.Client)
		if !ok {
			fmt.Println(ok)
		}

		ctx, ok := c.MustGet("ctx").(context.Context)
		if !ok {
			fmt.Println(ok)
		}

		GetCapsules(c, firestoreClient, authClient, ctx)
	})

	router.POST("/create", func(c *gin.Context) {

		authClient, ok := c.MustGet("authConn").(auth.Client)
		if !ok {
			fmt.Println(ok)
		}

		firestoreClient, ok := c.MustGet("firestoreConn").(firestore.Client)
		if !ok {
			fmt.Println(ok)
		}

		ctx, ok := c.MustGet("ctx").(context.Context)
		if !ok {
			fmt.Println(ok)
		}

		PostCapsule(c, firestoreClient, authClient, ctx)
	})

	router.POST("/delete", func(c *gin.Context) {
		authClient, ok := c.MustGet("authConn").(auth.Client)
		if !ok {
			fmt.Println(ok)
		}

		firestoreClient, ok := c.MustGet("firestoreConn").(firestore.Client)
		if !ok {
			fmt.Println(ok)
		}

		ctx, ok := c.MustGet("ctx").(context.Context)
		if !ok {
			fmt.Println(ok)
		}

		DeleteCapsule(c, firestoreClient, authClient, ctx)
	})

	router.POST("/open", func(c *gin.Context) {

		authClient, ok := c.MustGet("authConn").(auth.Client)
		if !ok {
			fmt.Println(ok)
		}

		firestoreClient, ok := c.MustGet("firestoreConn").(firestore.Client)
		if !ok {
			fmt.Println(ok)
		}

		ctx, ok := c.MustGet("ctx").(context.Context)
		if !ok {
			fmt.Println(ok)
		}

		OpenCapsule(c, firestoreClient, authClient, ctx)
	})

	router.GET("/detail", func(c *gin.Context) {

		authClient, ok := c.MustGet("authConn").(auth.Client)
		if !ok {
			fmt.Println(ok)
		}

		firestoreClient, ok := c.MustGet("firestoreConn").(firestore.Client)
		if !ok {
			fmt.Println(ok)
		}

		ctx, ok := c.MustGet("ctx").(context.Context)
		if !ok {
			fmt.Println(ok)
		}

		GetCapsuleDetail(c, firestoreClient, authClient, ctx)
	})

	router.Run()

}
