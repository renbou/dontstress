package interfaces

import "github.com/renbou/dontstress/internal/models"

type AdminDao interface {
	Get(admin *models.Admin) (*models.Admin, error)
}
