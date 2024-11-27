package firestore

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func Setup() *firestore.Client {
	env := os.Getenv("ENV")
	ctx := context.Background()
	var sa option.ClientOption
	if env == "PRODUCTION" {
		sa = option.WithCredentialsFile("../../prod-firebase-service-account.json")
	} else {
		sa = option.WithCredentialsFile("../../staging-firebase-service-account.json")
	}

	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	return client
}
