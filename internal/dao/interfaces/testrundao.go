package interfaces

import "github.com/renbou/dontstress/internal/models"

type TestrunDao interface {
	Create(testrun *models.Run) error
	Delete(testrun *models.Run) error
	Update(testrun *models.Run) error
	GetById(id string) (*models.Run, error)
}