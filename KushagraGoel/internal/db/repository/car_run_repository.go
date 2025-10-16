package repository

import (
	"KushagraGoel/models"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CarRunRepository interface {
	Create(ctx context.Context, carRun *models.Car_Run) error
	Update(ctx context.Context, carRun *models.Car_Run) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	List(ctx context.Context) ([]*models.Car_Run, error)
}

type MongoCarRepository struct {
	collection *mongo.Collection
}

func NewMongoCarRepository(db *mongo.Database) *MongoCarRepository {
	return &MongoCarRepository{
		collection: db.Collection("cars"),
	}
}

func (r *MongoCarRepository) Create(ctx context.Context, carRun *models.Car_Run) error {
	_, err := r.collection.InsertOne(ctx, carRun)
	if err != nil {
		return fmt.Errorf("failed to create car: %w", err)
	}
	return nil
}

func (r *MongoCarRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete car: %w", err)
	}
	return nil
}

func (r *MongoCarRepository) List(ctx context.Context) ([]*models.Car_Run, error) {
	filter := bson.D{}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var car []*models.Car_Run

	if err = cursor.All(ctx, &car); err != nil {
		return nil, err
	}

	return car, nil
}

func (r *MongoCarRepository) Update(ctx context.Context, carRun *models.Car_Run) error {
	if carRun.ID.IsZero() {
		return fmt.Errorf("missing car ID")
	}	

	filter := bson.M{"_id": carRun.ID}
    update := bson.M{
        "$set": bson.M{
            "date_uploaded": carRun.Date_uploaded,
            "location":      carRun.Location,
            "car_model":     carRun.Car_model,
            "event_type":    carRun.Event_type,
            "notes":         carRun.Notes,
            "file": bson.M{
                "aws_bucket": carRun.File.Aws_bucket,
                "file_path":  carRun.File.File_path,
                "file_name":  carRun.File.File_name,
            },
        },
    }

	result, err := r.collection.UpdateOne(ctx, filter, update)

	if err != nil {
		return fmt.Errorf("failed to update car: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("no car found with ID %s", carRun.ID.Hex())
	}

	return nil
}