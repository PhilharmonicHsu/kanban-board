package models

import "time"

type RoughCard struct {
	ID        int64     `bson:"_id" json:"id"`
	Title     string    `bson:"title" json:"title"`
	SpendTime string    `bson:"spend_time" json:"spend_time"`
	StartedAt time.Time `bson:"started_at" json:"started_at"`
	Type      int32     `bson:"type" json:"type"`
	Fields    []Field   `bson:"fields" json:"fields"`
	Members   []User    `bson:"members" json:"members"`
}

type Card struct {
	ID        int64     `bson:"_id" json:"id"`
	Title     string    `bson:"title" json:"title"`
	SpendTime string    `bson:"spend_time" json:"spend_time"`
	StartedAt time.Time `bson:"started_at" json:"started_at"`
	Type      int32     `bson:"type" json:"type"`
}
