package models

import "github.com/renbou/dontstress/internal/dto"

type Lab struct {
	Id   string `dynamo:"id" json:"id"`
	Name string `dynamo:"name" json:"name" validate:"required"`
}

func (lab *Lab) ToDTO() *dto.LabDTO {
	return &dto.LabDTO{Id: lab.Id, Name: lab.Name}
}
