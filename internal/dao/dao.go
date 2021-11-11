package dao

import (
	"github.com/renbou/dontstress/internal/dao/dynamodb"
	"github.com/renbou/dontstress/internal/dao/interfaces"
)

func AdminDao() interfaces.AdminDao {
	return &dynamodb.AdminImpl{}
}

func FileDao() interfaces.FileDao {
	return &dynamodb.FileImpl{}
}

func LabDao() interfaces.LabDao {
	return &dynamodb.LabImpl{}
}

func TaskDao() interfaces.TaskDao {
	return &dynamodb.TaskImpl{}
}

func TestrunDao() interfaces.TestrunDao {
	return &dynamodb.TestrunImpl{}
}
