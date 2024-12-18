package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
)

// Struct for post message requests
type PostMessageCapsule struct {
	Message    string `json:"message"`
	UnlockDate string `json:"unlock_date"`
}

// Str
type Capsule struct {
	Id          string    `json:"id", firestore:"-"`
	CreatedDate time.Time `json:"created", firestore:"created"`
	UnlockDate  time.Time `json:"unlock_date", firestore:"unlock_date"`
	Unlocked    bool      `json:"unlocked", firestore:"unlocked"`
}

// Handler function to create a capsule
func PostCapsule(c *gin.Context, client firestore.Client, authClient auth.Client) {

	ctx := context.Background()

	var message PostMessageCapsule
	if err := c.BindJSON(&message); err != nil {
		return
	}

	fmt.Println(message)

	authToken, err := authClient.VerifyIDToken(context.Background(), c.Request.Header["Token"][0])
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	unlockDate, err := time.Parse("02-01-2006 15:04:05", message.UnlockDate)
	if err != nil {
		log.Fatalf("Failed adding message: %v", err)
	}

	//Map representing the capsule to be posted
	capsule := map[string]interface{}{
		"user":        authToken.UID,
		"message":     message.Message,
		"created":     time.Now(),
		"unlock_date": unlockDate,
		"unlocked":    false,
	}

	//Add capsule to database using map
	_, _, postErr := client.Collection("messages").Add(ctx, capsule)
	if postErr != nil {
		log.Fatalf("Failed adding message: %v", err)
	}

	//return reposnse with capsule
	c.JSON(http.StatusOK, capsule)
}

// Handler function to return all capsules owned by the requestor
// Does not show the message of the capsule
func GetCapsules(c *gin.Context, client firestore.Client, authClient auth.Client) {

	ctx := context.Background()

	uid, err := GetRequestUID(c.Request.Header["Token"][0], authClient)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	query := client.Collection("messages").Where("user", "==", uid)

	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		log.Fatalf("Failed to retrieve documents: %v", err)
	}

	var capsule Capsule
	capsules := []Capsule{}

	for _, doc := range docs {

		data := doc.Data()
		doc.DataTo(&capsule)
		capsule.Id = doc.Ref.ID
		// Some manual marshalling due to lack of implicit time.Time conversion
		capsule.CreatedDate = data["created"].(time.Time)
		capsule.UnlockDate = data["unlock_date"].(time.Time)
		capsules = append(capsules, capsule)

	}

	c.JSON(http.StatusOK, gin.H{"capsules": capsules})

}

// Handler function to show all details of the capsule
// Will only return a succesfull response if the capsule is opened
func GetCapsuleDetail(c *gin.Context, client firestore.Client, authClient auth.Client) {

	ctx := context.Background()

	// Get User who sent request
	uid, err := GetRequestUID(c.Request.Header["Token"][0], authClient)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	capulseID := c.Query("id")

	query := client.Collection("messages").Doc(capulseID)

	doc, err := query.Get(ctx)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	capsule := doc.Data()

	if capsule["user"] != uid {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if capsule["unlocked"] != true {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "capsule locked"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": capsule})

}

// Handler function to delete a capsule given an ID
func DeleteCapsule(c *gin.Context, client firestore.Client, authClient auth.Client) {

	ctx := context.Background()

	// Get User who sent request
	uid, err := GetRequestUID(c.Request.Header["Token"][0], authClient)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	capulseID := c.Query("id")

	query := client.Collection("messages").Doc(capulseID)

	doc, err := query.Get(ctx)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	capsule := doc.Data()

	if capsule["user"] != uid {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	_, delErr := client.Collection("messages").Doc(capulseID).Delete(ctx)

	if delErr != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "capsule deleted"})

}

// Handler function to open a capsule
// Will only succeed if the opening time if before the current time.
func OpenCapsule(c *gin.Context, client firestore.Client, authClient auth.Client) {

	ctx := context.Background()

	// Get User who sent request
	uid, err := GetRequestUID(c.Request.Header["Token"][0], authClient)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	capulseID := c.Query("id")

	query := client.Collection("messages").Doc(capulseID)

	doc, err := query.Get(ctx)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	capsule := doc.Data()

	if capsule["user"] != uid {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if capsule["unlocked"].(bool) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "capsule already opened"})
		return
	}

	// Check if capsule can be opened
	if capsule["unlock_date"].(time.Time).Before(time.Now()) {
		_, err := client.Collection("messages").Doc(capulseID).Update(ctx, []firestore.Update{
			{
				Path:  "unlocked",
				Value: true,
			},
		})
		if err != nil {
			// Handle any errors in an appropriate way, such as returning them.
			log.Printf("An error has occurred: %s", err)
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "too early to open capsule"})
		return
	}

	capsule["unlocked"] = true
	c.JSON(http.StatusOK, gin.H{"message": capsule})

}
