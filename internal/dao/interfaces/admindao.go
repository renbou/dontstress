package interfaces

import "github.com/renbou/dontstress/internal/models"

type AdminDao interface {
	Get(file *models.Admin) (*models.Admin, error)
}
