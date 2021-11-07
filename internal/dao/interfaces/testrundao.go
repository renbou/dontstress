package interfaces

import "github.com/renbou/dontstress/internal/models"

type TestrunDao interface {
	Create(testRun *models.Run) error
	Delete(testRun *models.Run) error
	Update(testRun *models.Run) error
	GetById(id string) (*models.Run, error)
}
