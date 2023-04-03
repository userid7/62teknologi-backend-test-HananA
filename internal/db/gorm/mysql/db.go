package mysql

import (
	"backend-test/config"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewGormMysql(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", cfg.MYSQL.Username, cfg.MYSQL.Password, cfg.MYSQL.Host, cfg.MYSQL.Port, cfg.MYSQL.Dbname)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})

}
