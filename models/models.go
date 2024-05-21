package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id"`
	User_name *string            `json:"user_name" validate:"required,min=2,max=100"`
	Password  *string            `json:"Password" validate:"required,min=6"`
	Email     *string            `json:"email" validate:"email,required"`
	// Token            *string            `json:"token"`
	// Refresh_token    *string            `json:"refresh_token"`
	// Created_at       time.Time          `json:"created_at"`
	// Updated_at       time.Time          `json:"updated_at"`
	User_id          string    `json:"user_id"`
	Codeforce_handle *string   `json:"codeforce_handle" validate:"required"`
	Subscribedblogs  *[]string `json:"subscribedblogs"`
}
type BlogEntry struct {
	Id                      int      `json:"id"`
	OriginalLocale          string   `json:"originalLocale"`
	CreationTimeSeconds     int      `json:"creationTimeSeconds"`
	AuthorHandle            string   `json:"authorHandle"`
	Title                   string   `json:"title"`
	Content                 string   `json:"content"`
	Locale                  string   `json:"locale"`
	ModificationTimeSeconds int      `json:"modificationTimeSeconds"`
	AllowViewHistory        bool     `json:"allowViewHistory"`
	Tags                    []string `json:"tags"`
	Rating                  int      `json:"rating"`
}

type Comment struct {
	Id                  int    `json:"id" bson:"id"`
	CreationTimeSeconds int    `json:"creationTimeSeconds" bson:"creationTimeSeconds"`
	CommentatorHandle   string `json:"commentatorHandle" bson:"commentatorHandle"`
	Locale              string `json:"locale" bson:"locale"`
	Text                string `json:"text" bson:"text"`
	ParentCommentId     int    `json:"parentCommentId" bson:"parentCommentId"`
	Rating              int    `json:"rating" bson:"rating"`
}

type RecentAction struct {
	TimeSeconds int64      `json:"timeSeconds" bson:"timeSeconds"`
	BlogEntry   *BlogEntry `json:"blogEntry" bson:"blogEntry"`
	Comment     *Comment   `json:"comment" bson:"comment"`
}
