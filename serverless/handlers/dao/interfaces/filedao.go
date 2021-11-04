package interfaces

import "github.com/renbou/dontstress/serverless/handlers/models"

type FileDao interface {
	Create(file *models.File) error
	Delete(file *models.File) error
}
