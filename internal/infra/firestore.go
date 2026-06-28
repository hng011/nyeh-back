package infra

import (
	"context"
	"log"
	"nyeh-back/internal/core"

	"cloud.google.com/go/firestore"
)

// InitFirestore initializes the Firebase app and returns the Firestore client.
func InitFirestore(ctx context.Context) *firestore.Client {

	client, err := firestore.NewClientWithDatabase(ctx, core.Settings.GOOGLE_CLOUD_PROJECT, core.Settings.FIREBASE_DATABASE_ID)
	if err != nil {
		log.Fatalf("error initializing firestore client: %v\n", err)
	}

	return client
}
