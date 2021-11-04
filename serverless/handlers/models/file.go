package models

type File struct {
	Id   string `dynamo:"id"`
	Lang string `dynamo:"lang"`
}
