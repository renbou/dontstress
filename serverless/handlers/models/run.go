package models

type Run struct {
	Id     string `dynamo:"id"`
	Labid  string `dynamo:"labid"`
	Taskid int    `dynamo:"taskid"`
	Fileid string `dynamo:"fileid"`
	Status string `dynamo:"status"`
	Tests  []struct {
		Result  string `dynamo:"results"`
		Message string
		Info    struct {
			Test     string `dynamo:"test"`
			Expected string `dynamo:"expected"`
			Got      string `dynamo:"got"`
		} `dynamo:"info"`
	}
}
