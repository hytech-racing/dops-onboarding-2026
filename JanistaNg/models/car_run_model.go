package models

import (
    "time"

    "go.mongodb.org/mongo-driver/bson/primitive"
)

type CarRun struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    DateUploaded time.Time          `bson:"date_uploaded" json:"date_uploaded"`
    Location     string             `bson:"location" json:"location"`
    CarModel     string             `bson:"car_model" json:"car_model"`
    EventType    string             `bson:"event_type" json:"event_type"`
    Notes        string             `bson:"notes" json:"notes"`
    File         FileInfo           `bson:"file" json:"file"`
}
type FileInfo struct {
	AwsBucket string `bson:"aws_bucket" json:"aws_bucket"`
    FilePath  string `bson:"file_path" json:"file_path"`
    FileName  string `bson:"file_name" json:"file_name"`
}

func NewCarRun() (*CarRun, error) {
	return &CarRun{
		ID: primitive.NewObjectID(),
		DateUploaded: time.Now(),
	}, nil
}