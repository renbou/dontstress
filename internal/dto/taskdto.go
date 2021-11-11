package dto

type TaskDTO struct {
	Id   int    `dynamo:"id" json:"id"`
	Name string `dynamo:"name" json:"name"`
}
