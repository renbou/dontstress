package models

import "github.com/renbou/dontstress/internal/dto"

type File struct {
	Id   string `dynamo:"id" json:"id"`
	Lang string `dynamo:"lang" json:"lang"`
}

func (file *File) ToDTO(data string) *dto.FileDTO {
	return &dto.FileDTO{Lang: file.Lang, Data: data}
}
