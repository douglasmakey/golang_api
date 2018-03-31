package repositories

import (
	"github.com/jinzhu/gorm"

	"github.com/douglasmakey/backend_base/models"
)

type UserRepository struct {
	*Repository
}

func NewUserRepo(db *gorm.DB) *UserRepository {
	return &UserRepository{&Repository{db}}
}


func (ur *UserRepository) FindByCredentials(user *models.User) bool {
	if ur.DB.Where("email = ? and password = ?", user.Email, user.Password).First(&user).RecordNotFound() {
		return false
	}

	return true
}

