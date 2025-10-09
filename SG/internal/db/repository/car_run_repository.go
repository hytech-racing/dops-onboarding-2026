package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"main/models"
)

// CarRunRepository defines the interface
type CarRunRepository interface {
	Create(ctx context.Context, carRun *models.CarRun) error
	Update(ctx context.Context, carRun *models.CarRun) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	List(ctx context.Context) ([]*models.CarRun, error)
}

// MongoCarRunRepository implements CarRunRepository
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
		return err
	}
	return nil
}

func (r *MongoCarRunRepository) Update(ctx context.Context, carRun *models.CarRun) error {
	// You can enhance this later with $set logic
	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": carRun.ID}, carRun)
	return err
}

func (r *MongoCarRunRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *MongoCarRunRepository) List(ctx context.Context) ([]*models.CarRun, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var runs []*models.CarRun
	if err = cursor.All(ctx, &runs); err != nil {
		return nil, err
	}
	return runs, nil
}
