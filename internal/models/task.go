package models

import "github.com/renbou/dontstress/internal/dto"

type Task struct {
	LabId     string `dynamo:"labid" json:"labid"`
	Num       int    `dynamo:"num" json:"num"`
	Name      string `dynamo:"name" json:"name" validate:"required"`
	Validator string `dynamo:"validator" json:"validator"`
	Generator string `dynamo:"generator" json:"generator"`
}

func (task *Task) ToDTO() *dto.TaskDTO {
	return &dto.TaskDTO{Id: task.Num, Name: task.Name}
}
