package mysqlx

import (
	"Tiktok/configs"
	"Tiktok/log/zlog"
	"Tiktok/pkg/database"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Mysql struct {
}

// InitDataBases 初始化
func (m *Mysql) InitDataBase(config configs.Config) (*gorm.DB, error) {
	dsn := m.GetDataBaseDsn(config)
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		zlog.Panicf("MySQL无法连接数据库！: %v", err)
		return nil, err
	}
	zlog.Infof("MySQL连接数据库成功！")
	return db, nil
}
func (m *Mysql) GetDataBaseDsn(config configs.Config) string {
	return config.DB.MySQL.Dsn
}
func NewMySql() database.DataBase {
	return &Mysql{}
}
