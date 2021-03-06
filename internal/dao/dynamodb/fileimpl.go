package dynamodb

import (
	"github.com/renbou/dontstress/internal/models"
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

func (dao *FileImpl) Update(fileId string, connectionId string) error {
	db := getDB()
	table := db.Table(FilesDynamoName)
	update := table.Update("id", fileId)
	update.Set("connectionId", connectionId)
	return update.Run()
}
