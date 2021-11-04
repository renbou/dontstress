package models

type File struct {
	Id   string `dynamo:"id" json:"id"`
	Lang string `dynamo:"lang" json:"lang"`
}
