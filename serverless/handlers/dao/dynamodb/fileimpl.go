package dynamodb

import (
	"github.com/renbou/dontstress/serverless/handlers/models"
)

type FileImpl struct {
}

func (dao *FileImpl) Create(file *models.File) error {
	db := getDB()
	table := db.Table(FilesDynamoName)
	return table.Put(file).Run()
}

func (dao *FileImpl) Delete(file *models.File) error {
	db := getDB()
	table := db.Table(FilesDynamoName)
	return table.Delete("id", file.Id).Run()
}
