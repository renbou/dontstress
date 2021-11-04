package dynamodb

import (
	"github.com/renbou/dontstress/internal/models"
)

type LabImpl struct {
}

func (dao *LabImpl) Create(lab *models.Lab) error {
	db := getDB()
	table := db.Table(LabsDynamoName)
	return table.Put(lab).Run()
}

func (dao *LabImpl) Delete(lab *models.Lab) error {
	db := getDB()
	table := db.Table(LabsDynamoName)
	return table.Delete("id", lab.Id).Run()
}

func (dao *LabImpl) GetById(labId string) (*models.Lab, error) {
	db := getDB()
	table := db.Table(LabsDynamoName)
	var result models.Lab
	err := table.Get("id", labId).One(&result)
	return &result, err
}

func (dao *LabImpl) GetAll() ([]models.Lab, error) {
	db := getDB()
	table := db.Table(LabsDynamoName)
	var results []models.Lab
	err := table.Scan().All(&results)
	return results, err
}
