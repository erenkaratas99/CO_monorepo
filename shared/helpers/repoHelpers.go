package helpers

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func ObjectExists(col *mongo.Collection, objectId string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	filter := bson.M{"_id": objectId}
	projection := options.FindOne().SetProjection(bson.M{"_id": 1})
	result := col.FindOne(ctx, filter, projection)
	if err := result.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func CountDocuments(col *mongo.Collection, ctx context.Context, filter *bson.D) int {
	totalCount64, _ := col.CountDocuments(ctx, filter)
	return int(totalCount64)
}
