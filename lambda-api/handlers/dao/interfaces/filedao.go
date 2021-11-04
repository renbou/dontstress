package interfaces

import "github.com/renbou/dontstress/lambda-api/handlers/models"

type FileDao interface {
	Create(file *models.File) error
	Delete(file *models.File) error
}
