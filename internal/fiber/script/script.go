package script

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"jjsdapi/internal/fiber/result"
	"jjsdapi/internal/fiber/validate"
	"jjsdapi/internal/repository/dao"
	"jjsdapi/internal/repository/model"
	"time"
)

func Boot(ctx *fiber.Ctx) (e error) {

	pwd := time.Now().Format("20060102")

	if ctx.Params("pwd") != pwd {
		return result.Json(ctx, nil, errors.New("check"))
	}

	var form struct {
		Data  string `json:"data"`
		Scene string `json:"scene" validate:"required"`
	}
	if err := validate.CheckParams(ctx, &form); err != nil {
		return result.Json(ctx, nil, err)
	}

	var data = make(map[string]interface{})
	switch form.Scene {
	case "syncAdminUser":
		data, e = syncAdminUser()
	}

	return result.Json(ctx, data, e)

}

func syncAdminUser() (data map[string]interface{}, err error) {
	client, err := gorm.Open(mysql.Open("root:Ws123456@tcp(120.79.141.163:3306)/jjsd?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{
		Logger: logger.Discard,
	})
	if err != nil {
		return nil, err
	}

	var list []*model.AdminUser
	client.Table("sd_admin").Where("is_del = ?", 0).Find(&list)
	b, _ := json.Marshal(list)

	dao.AdminUser.WithContext(context.Background()).Create(list...)
	fmt.Println(string(b))

	return
}
