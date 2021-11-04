package interfaces

import "github.com/renbou/dontstress/lambda-api/handlers/models"

type TaskDao interface {
	Create(task *models.Task) error
	Delete(task *models.Task) error
	Update(task *models.Task) error
	GetCount(labId string) (int, error)
	GetAll(labId string) ([]models.Task, error)
}
