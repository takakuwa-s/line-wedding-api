package dto

import (
	"fmt"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/takakuwa-s/line-wedding-api/conf"
	"google.golang.org/api/option"
)

type Firestore struct {
	*firestore.Client
}

// Newコンストラクタ
func NewFirestore() *Firestore {
	// Use a service account
	sa := option.WithCredentialsFile("./resource/common/wedding-343414-3910a4e5d7e9.json")
	app, err := firebase.NewApp(conf.Ctx, nil, sa)
	if err != nil {
		panic(fmt.Sprintf("Failed to create the firebase app instance; err = %v", err))
	}

	client, err := app.Firestore(conf.Ctx)
	if err != nil {
		panic(fmt.Sprintf("Failed to create the firestore instance; err = %v", err))
	}
	// defer client.Close()
	return &Firestore{client}
}