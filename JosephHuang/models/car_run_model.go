package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CarRun struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	DateUploaded time.Time          `json:"date_uploaded"`
	Location     string             `json:"location"`
	CarModel     string             `json:"car_model"`
	EventType    string             `json:"event_type"`
	Notes        string             `json:"notes"`
	File         File               `json:"file"`
}

type File struct {
	AWSBucket string `json:"aws_bucket"`
	FilePath  string `json:"file_path"`
	FileName  string `json:"file_name"`
}

func NewCarRun(location string, carModel string, eventType string, notes string, file File) CarRun {
	return CarRun{
		ID:           primitive.NilObjectID,
		DateUploaded: time.Now().UTC(),
		Location:     location,
		CarModel:     carModel,
		EventType:    eventType,
		Notes:        notes,
		File:         file,
	}
}
