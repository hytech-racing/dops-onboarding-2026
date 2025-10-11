package repository

import (
	"context"
	"errors"
	"main/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CarRunRepository interface {
	Create(ctx context.Context, carRun *models.CarRun) error
	Update(ctx context.Context, carRun *models.CarRun, id primitive.ObjectID) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	List(ctx context.Context) ([]*models.CarRun, error)
}

type MongoCarRunRepository struct {
	collection *mongo.Collection
}

func NewMongoCarRepository(db *mongo.Database) (*MongoCarRunRepository, error) {
	if db == nil {
		return nil, errors.New("nil *mongo.Database")
	}
	return &MongoCarRunRepository{
		collection: db.Collection("vehicle_run"),
	}, nil
}

func (r *MongoCarRunRepository) Create(ctx context.Context, carRun *models.CarRun) error {
	_, err := r.collection.InsertOne(ctx, carRun)
	if err != nil {
		return errors.New("failed to create new car run: " + err.Error())
	}
	return nil

}

func (r *MongoCarRunRepository) Update(ctx context.Context, carRun *models.CarRun, id primitive.ObjectID) error {
	filter := bson.D{{Key: "_id", Value: id}}
	_, err := r.collection.ReplaceOne(ctx, filter, carRun)
	if err != nil {
		return errors.New("failed to replace car run: " + err.Error())
	}
	return nil
}

func (r *MongoCarRunRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.D{{Key: "_id", Value: id}}
	_, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return errors.New("failed to delete car run: " + err.Error())
	}
	return nil
}

func (r *MongoCarRunRepository) List(ctx context.Context) ([]*models.CarRun, error) {
	cursor, err := r.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, errors.New("failed to retrieve car run documents: " + err.Error())
	}
	defer cursor.Close(ctx)
	var results []*models.CarRun
	err = cursor.All(ctx, &results)
	if err != nil {
		return nil, errors.New("failed to read cursor: " + err.Error())
	}
	return results, nil
}
