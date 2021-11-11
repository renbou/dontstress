package dynamodb

import "github.com/renbou/dontstress/internal/models"

type AdminImpl struct {
}

func (dao *AdminImpl) Get(admin *models.Admin) (*models.Admin, error) {
	db := getDB()
	table := db.Table(AdminsDynamoName)
	var result models.Admin
	err := table.Get("id", admin.Id).One(&result)
	return &result, err
}
