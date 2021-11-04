package dynamodb

import (
	"github.com/renbou/dontstress/serverless/handlers/models"
)

type TaskImpl struct {
}

func (dao TaskImpl) Create(task models.Task) error {
	db := getDB()
	table := db.Table(TasksDynamoName)
	return table.Put(task).Run()
}

func (dao TaskImpl) Delete(task models.Task) error {
	db := getDB()
	table := db.Table(TasksDynamoName)
	return table.Delete("labid", task.LabId).Range("num", task.Num).Run()
}

func (dao TaskImpl) GetAll(labId string) ([]models.Task, error) {
	db := getDB()
	table := db.Table(TasksDynamoName)
	var results []models.Task
	err := table.Get("labid", labId).All(&results)
	return results, err
}
