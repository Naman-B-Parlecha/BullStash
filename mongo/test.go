package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/Naman-B-Parlecha/BullStash/util"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestConnection(mongoURI string) error {
	if mongoURI == "" {
		return fmt.Errorf("MongoDB URI cannot be empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("connection failed: %v", err)
	}
	defer client.Disconnect(ctx)

	err = client.Ping(ctx, nil)
	if err != nil {
		return fmt.Errorf("ping failed: %v", err)
	}

	successMsg := fmt.Sprintf("MongoDB connection established successfully to: %s", clientOptions.GetURI())
	util.CallWebHook(successMsg, false)
	return nil
}
