package models

type Lab struct {
	Id   string `dynamo:"id"`
	Name string `dynamo:"name"`
}
