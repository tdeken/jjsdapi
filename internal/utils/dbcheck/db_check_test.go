package dbcheck

import (
	"gorm.io/gorm"
	"testing"
)

func TestDbError(t *testing.T) {

	var err = gorm.ErrRecordNotFound

	t.Log(DbError(err))

	t.Log(err)
}
