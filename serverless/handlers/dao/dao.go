package dao

import (
	"github.com/renbou/dontstress/serverless/handlers/dao/dynamodb"
	"github.com/renbou/dontstress/serverless/handlers/dao/interfaces"
)

func FileDao() interfaces.FileDao {
	return &dynamodb.FileImpl{}
}

func LabDao() interfaces.LabDao {
	return &dynamodb.LabImpl{}
}

func TaskDao() interfaces.TaskDao {
	return &dynamodb.TaskImpl{}
}
