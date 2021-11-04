package interfaces

import "github.com/renbou/dontstress/serverless/handlers/models"

type TaskDao interface {
	Create(task *models.Task) error
	Delete(task *models.Task) error
	GetAll(labId string) ([]models.Lab, error)
}
