package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"SophiaNemes/models"
)

type CarRunRepository interface {
	Create(ctx context.Context, carRun *models.CarRun) error
	Update(ctx context.Context, carRun *models.CarRun) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	List(ctx context.Context) ([]*models.CarRun, error)
}

type MongoCarRunRepository struct {
	collection *mongo.Collection
}

func NewMongoCarRunRepository(db *mongo.Database) *MongoCarRunRepository {
	return &MongoCarRunRepository{
		collection: db.Collection("car_runs"),
	}
}

func (r *MongoCarRunRepository) Create(ctx context.Context, carRun *models.CarRun) error {
	_, err := r.collection.InsertOne(ctx, carRun)
	if err != nil {
		return fmt.Errorf("failed to create car run: %w", err)
	}
	return nil
}

func (r *MongoCarRunRepository) Update(ctx context.Context, carRun *models.CarRun) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": carRun.ID}, bson.M{"$set": carRun})
	if err != nil {
		return fmt.Errorf("failed to update car run: %w", err)
	}
	return nil
}

func (r *MongoCarRunRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete car run: %w", err)
	}
	return nil
}

func (r *MongoCarRunRepository) List(ctx context.Context) ([]*models.CarRun, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to list car runs: %w", err)
	}
	defer cursor.Close(ctx)

	var carRuns []*models.CarRun
	if err = cursor.All(ctx, &carRuns); err != nil {
		return nil, fmt.Errorf("failed to decode car runs: %w", err)
	}

	return carRuns, nil
}