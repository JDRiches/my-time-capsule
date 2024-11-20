package main

import (
	"context"

	"firebase.google.com/go/auth"
)

func GetRequestUID(idToken string, authClient auth.Client) (string, error) {
	authToken, err := authClient.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		return "", err
	}

	return authToken.UID, nil
}
