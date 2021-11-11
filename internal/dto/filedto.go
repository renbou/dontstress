package dto

type FileDTO struct {
	Lang string `dynamo:"lang" json:"lang"`
	Data string `dynamo:"data" json:"data"`
}
