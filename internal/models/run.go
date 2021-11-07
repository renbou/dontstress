package models

import "github.com/renbou/dontstress/internal/dto"

type Run struct {
	Id     string `dynamo:"id" json:"id"`
	Labid  string `dynamo:"labid" json:"labid"`
	Taskid int    `dynamo:"taskid" json:"taskid"`
	Fileid string `dynamo:"fileid" json:"fileid"`
	Status string `dynamo:"status" json:"status"`
	Tests  []struct {
		Result  string `dynamo:"results" json:"result"`
		Message string `dynamo:"message" json:"message"`
		Info    struct {
			Test     string `dynamo:"test" json:"test"`
			Expected string `dynamo:"expected" json:"expected"`
			Got      string `dynamo:"got" json:"got"`
		} `dynamo:"info" json:"info"`
	} `dynamo:"tests" json:"tests"`
}

func (run *Run) ToDTO() *dto.RunDTO {
	return &dto.RunDTO{Status: run.Status, Tests: run.Tests}
}
