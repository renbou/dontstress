package models

type Run struct {
	Id     string `dynamo:"id" json:"-"`
	Labid  string `dynamo:"labid" json:"-"`
	Taskid int    `dynamo:"taskid" json:"-"`
	Fileid string `dynamo:"fileid" json:"-"`
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
