package dto

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

type Firestore struct {
	*firestore.Client
	Ctx context.Context
}

// Newコンストラクタ
func NewFirestore() *Firestore {
	// Use a service account
	ctx := context.Background()
	sa := option.WithCredentialsFile("./resource/common/wedding-343414-3910a4e5d7e9.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		panic(fmt.Sprintf("Failed to create the firebase app instance; err = %v", err))
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		panic(fmt.Sprintf("Failed to create the firestore instance; err = %v", err))
	}
	// defer client.Close()
	return &Firestore{client, ctx}
}