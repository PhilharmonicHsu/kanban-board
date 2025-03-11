package models

import "time"

type User struct {
	ID        int32     `bson:"_id" json:"id"`
	Account   string    `bson:"account" json:"account"`
	Password  string    `bson:"password" json:"password"`
	Name      string    `bson:"name" json:"name"`
	Avatar    *string   `bson:"avatar" json:"avatar"` // 有可能為null時，就使用指標
	CreateAt  time.Time `bson:"created_at" json:"create_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

type UserResponse struct {
	ID        int32   `bson:"_id" json:"id"`
	Account   string  `bson:"account" json:"account"`
	Password  string  `bson:"password" json:"password"`
	Name      string  `bson:"name" json:"name"`
	Avatar    *string `bson:"avatar" json:"avatar"`
	CreateAt  string  `bson:"created_at" json:"create_at"`
	UpdatedAt string  `bson:"updated_at" json:"updated_at"`
}
