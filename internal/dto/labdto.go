package dto

type LabDTO struct {
	Id   string `dynamo:"id" json:"id"`
	Name string `dynamo:"name" json:"name"`
}
