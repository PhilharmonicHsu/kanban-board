package models

type Field struct {
	ID    int64  `bson:"_id" json:"id"`
	Name  string `bson:"name" json:"name"`
	Group int64  `bson:"group" json:"group"`
	Value string `bson:"value" json:"value"`
}
