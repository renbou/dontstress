package dynamodb

import (
	"github.com/renbou/dontstress/internal/models"
)

type TaskImpl struct {
}

func (dao *TaskImpl) Create(task *models.Task) error {
	db := getDB()
	table := db.Table(TasksDynamoName)
	return table.Put(task).Run()
}

func (dao *TaskImpl) Delete(task *models.Task) error {
	db := getDB()
	table := db.Table(TasksDynamoName)
	return table.Delete("labid", task.LabId).Range("num", task.Num).Run()
}

func (dao *TaskImpl) Update(task *models.Task) error {
	db := getDB()
	table := db.Table(TasksDynamoName)
	update := table.Update("labid", task.LabId).Range("num", task.Num)
	if task.Generator != "" {
		update.Set("generator", task.Generator)
	}
	if task.Validator != "" {
		update.Set("validator", task.Validator)
	}
	return update.Run()
}

func (dao *TaskImpl) GetCount(labId string) (int, error) {
	db := getDB()
	table := db.Table(TasksDynamoName)
	count, err := table.Get("labid", labId).Count()
	return int(count), err
}

func (dao *TaskImpl) GetAll(labId string) ([]models.Task, error) {
	db := getDB()
	table := db.Table(TasksDynamoName)
	var results []models.Task
	err := table.Get("labid", labId).All(&results)
	return results, err
}
