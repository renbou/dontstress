package interfaces

import "github.com/renbou/dontstress/serverless/handlers/models"

type LabDao interface {
	Create(lab *models.Lab) error
	Delete(lab *models.Lab) error
	GetById(labId string) (models.Lab, error)
	GetAll() ([]models.Lab, error)
}
