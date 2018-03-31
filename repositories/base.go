package repositories

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/labstack/gommon/log"
)

type Repository struct {
	*gorm.DB
}

func (r *Repository) Find(model interface{}, filter, value string) bool {
	if r.DB.Where(fmt.Sprintf("%s = ?", filter), value).First(model).RecordNotFound() {
		return false
	}

	return true
}

func (r *Repository) Save(model interface{}) bool {
	err := r.DB.Save(model).Error
	if err != nil {
		log.Errorf("error: %v", err)
		return false
	}

	return true
}