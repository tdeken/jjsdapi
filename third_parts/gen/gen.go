package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io"
	"log"
	"os"
	"strings"
)

// dataMap mapping relationship
var dataMap = map[string]func(gorm.ColumnType) (dataType string){
	// int mapping
	"bigint": func(columnType gorm.ColumnType) (dataType string) {
		ct, _ := columnType.ColumnType()
		if strings.Index(ct, "unsigned") > -1 {
			return "uint64"
		} else {
			return "int64"
		}
	},

	// bool mapping
	"tinyint": func(columnType gorm.ColumnType) (dataType string) {
		ct, _ := columnType.ColumnType()
		if strings.HasPrefix(ct, "tinyint(1)") {
			return "bool"
		}
		return "int8"
	},
}

var url = "root:123456@tcp(127.0.0.1:3306)/bazaar?charset=utf8mb4&parseTime=True&loc=Local"

type Queries interface {
	// GetByID
	// SELECT * FROM @@table WHERE id=@id
	GetByID(id int) (*gen.T, error) // returns struct and error
}

func main() {
	log.SetOutput(io.Discard)
	gormDB, _ := gorm.Open(mysql.Open(url), &gorm.Config{
		Logger: logger.Discard,
	})

	path, _ := os.Getwd()

	gn := gen.NewGenerator(gen.Config{
		OutPath:      fmt.Sprintf("%s/internal/repository/dao", path),
		ModelPkgPath: "model",
		Mode:         gen.WithDefaultQuery, // generate mode

		// if you want the nullable field generation property to be pointer type, set FieldNullable true
		FieldNullable: true,
		// if you want to assign field which has a default value in the `Create` API, set FieldCoverable true, reference: https://gorm.io/docs/create.html#Default-Values
		FieldCoverable: false,
		// if you want to generate field with unsigned integer type, set FieldSignable true
		FieldSignable: true,
		// if you want to generate index tags from database, set FieldWithIndexTag true
		FieldWithIndexTag: true,
		// if you want to generate type tags from database, set FieldWithTypeTag true
		FieldWithTypeTag: true,
	})

	gn.UseDB(gormDB) // reuse your gorm db

	// specify diy mapping relationship
	gn.WithDataTypeMap(dataMap)

	gn.ApplyInterface(func(Queries) {},
		gn.GenerateModel("admin_user"),
		gn.GenerateModel("customer_address"),
		gn.GenerateModel("customer"),
	)

	// Generate the code
	gn.Execute()
}
