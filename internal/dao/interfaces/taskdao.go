package interfaces

import "github.com/renbou/dontstress/internal/models"

type TaskDao interface {
	Create(task *models.Task) error
	Delete(task *models.Task) error
	Update(task *models.Task) error
	GetCount(labId string) (int, error)
	GetAll(labId string) ([]models.Task, error)
}
