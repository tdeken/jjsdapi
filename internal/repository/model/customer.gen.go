// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameCustomer = "customer"

// Customer 客户配送地址表
type Customer struct {
	ID        int64  `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true" json:"id"`
	CreatedAt int64  `gorm:"column:created_at;type:bigint;not null;index:idx_created_at,priority:1" json:"created_at"`
	UpdatedAt int64  `gorm:"column:updated_at;type:bigint;not null" json:"updated_at"`
	DeletedAt int64  `gorm:"column:deleted_at;type:bigint;not null" json:"deleted_at"`
	Name      string `gorm:"column:name;type:varchar(255);not null;comment:客户名称" json:"name"`  // 客户名称
	Phone     string `gorm:"column:phone;type:varchar(15);not null;comment:配送地址" json:"phone"` // 配送地址
}

// TableName Customer's table name
func (*Customer) TableName() string {
	return TableNameCustomer
}