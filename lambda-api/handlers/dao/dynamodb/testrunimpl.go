package dynamodb

import (
	"github.com/renbou/dontstress/lambda-api/handlers/models"
)

type TestrunImpl struct {
}

func (dao *TestrunImpl) Create(testrun *models.Run) error {
	db := getDB()
	table := db.Table(TestrunsDynamoName)
	return table.Put(testrun).Run()
}

func (dao *TestrunImpl) Delete(testrun *models.Run) error {
	db := getDB()
	table := db.Table(TestrunsDynamoName)
	return table.Delete("id", testrun.Id).Run()
}

func (dao *TestrunImpl) Update(testrun *models.Run) error {
	db := getDB()
	table := db.Table(TestrunsDynamoName)
	update := table.Update("id", testrun.Id)
	// TODO: update more fields
	update.Set("status", testrun.Status)
	return update.Run()
}

func (dao *TestrunImpl) GetById(id string) (*models.Run, error) {
	db := getDB()
	table := db.Table(TestrunsDynamoName)
	var result models.Run
	err := table.Get("id", id).One(&result)
	return &result, err
}
