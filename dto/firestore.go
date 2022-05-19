package dto

import (
	"fmt"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go"
	"github.com/takakuwa-s/line-wedding-api/conf"
	"google.golang.org/api/option"
)

type Firestore struct {
	Firestore *firestore.Client
	Bucket *storage.BucketHandle
}

// Newコンストラクタ
func NewFirestore() *Firestore {
	config := &firebase.Config{
    StorageBucket: "wedding-343414.appspot.com",
	}

	// Use a service account
	sa := option.WithCredentialsFile("./resource/common/wedding-343414-3910a4e5d7e9.json")
	app, err := firebase.NewApp(conf.Ctx, config, sa)
	if err != nil {
		panic(fmt.Sprintf("Failed to create the firebase app instance; err = %v", err))
	}

	firestore, err := app.Firestore(conf.Ctx)
	if err != nil {
		panic(fmt.Sprintf("Failed to create the firestore; err = %v", err))
	}
	storage, err := app.Storage(conf.Ctx)
	if err != nil {
		panic(fmt.Sprintf("Failed to create the firebase storage; err = %v", err))
	}
	bucket, err := storage.DefaultBucket()
	if err != nil {
		panic(fmt.Sprintf("Failed to create the firebase storage bucket; err = %v", err))
	}
	return &Firestore{Firestore: firestore, Bucket: bucket}
}