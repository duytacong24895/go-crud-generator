package models

import (
	"time"

	"gorm.io/gorm"
)

const TableNameEmployee = "employee"

// Employee mapped from table <users>
type Employee struct {
	gorm.Model
	ID     int64     `gorm:"column:id;primaryKey" json:"id"`
	Name   string    `gorm:"column:name" json:"name"`
	Email  string    `gorm:"column:email" json:"email"`
	Dob    time.Time `gorm:"column:dob" json:"dob"`
	Age    int64     `gorm:"column:age" json:"age"`
	Phone  string    `gorm:"column:phone" json:"phone"`
	Mature bool      `gorm:"column:mature" json:"mature"`
}

// TableName Employee's table name
func (*Employee) TableName() string {
	return TableNameEmployee
}
