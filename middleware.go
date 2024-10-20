package main

import (
	"context"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
)

func FirestoreMiddleware(client firestore.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("firestoreConn", client)
		c.Next()
	}
}

func AuthMiddleware(auth auth.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("authConn", auth)
		c.Next()
	}
}

func CtxMiddleware(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("ctx", ctx)
		c.Next()
	}
}
