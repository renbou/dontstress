package models

type Task struct {
	LabId     string `dynamo:"labid" json:"lab_id"`
	Num       int    `dynamo:"num" json:"num"`
	Name      string `dynamo:"name" json:"name"`
	Validator string `dynamo:"validator" json:"validator"`
	Generator string `dynamo:"generator" json:"generator"`
}
