package models

type Task struct {
	LabId string `dynamo:"labid"`
	Num   int    `dynamo:"num"`
	Name  string `dynamo:"name"`

	//Lang      string `dynamo:"lang"`
	//Validator string `dynamo:"validator"`
	//Generator string `dynamo:"generator"`
}
