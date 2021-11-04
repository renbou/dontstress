package interfaces

import "github.com/renbou/dontstress/serverless/handlers/models"

type TaskDao interface {
	Create(task *models.Task) error
	Delete(task *models.Task) error
	Update(task *models.Task) error
	GetAll(labId string) ([]models.Task, error)
}
