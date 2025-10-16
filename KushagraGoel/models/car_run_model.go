package models

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Car_Run struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Date_uploaded string             `bson:"date_uploaded" json:"date_uploaded"`
	Location      string             `bson:"location" json:"location"`
	Car_model     string             `bson:"car_model" json:"car_model"`
	Event_type    string             `bson:"event_type" json:"event_type"`
	Notes         string             `bson:"notes" json:"notes"`
	File          FileInfo           `bson:"file" json:"file"`
}

type FileInfo struct {
	Aws_bucket string `bson:"aws_bucket" json:"aws_bucket"`
	File_path  string `bson:"file_path" json:"file_path"`
	File_name  string `bson:"file_name" json:"file_name"`
}

func NewCar(location, car_model, event_type, notes, aws_bucket, file_path, file_name string) (*Car_Run, error) {
	if location == "" {
		return nil, errors.New("location is required")
	}
	if car_model == "" {
		return nil, errors.New("car model is required")
	}

	return &Car_Run{
		ID:            primitive.NewObjectID(),
		Date_uploaded: time.Now().Format(time.RFC3339),
		Location:      location,
		Car_model:     car_model,
		Event_type:    event_type,
		Notes:         notes,
		File: FileInfo{
			Aws_bucket: aws_bucket,
			File_path:  file_path,
			File_name:  file_name,
		},
	}, nil
}
