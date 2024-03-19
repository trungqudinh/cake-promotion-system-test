package mysql

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"cake/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var dbOnce sync.Once

type MySqlStorage struct {
	db *gorm.DB
}

func InitDatabase() *MySqlStorage {
	mysqlStorage := &MySqlStorage{}
	dbOnce.Do(
		func() {
			db, err := initGormDatabase()
			if err != nil {
				panic(fmt.Errorf("failed to connect database, error: %v", err))
			}
			mysqlStorage.db = db
		},
	)

	return mysqlStorage
}

func initGormDatabase() (*gorm.DB, error) {
	cfg := config.GetAppConfig().Database
	db, err := newGormDatabase()
	if err != nil {
		return nil, err
	}

	if cfg.EnableAutoMigrate {
		err = autoMigrate(db)
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}

func newGormDatabase() (*gorm.DB, error) {
	var (
		cfg = config.GetAppConfig().Database
		dsn string
	)
	switch cfg.DBType {
	case "mysql":
		dsn = config.GetAppConfig().Mysql.GetDSN()
	case "sqlite3":
		dsn = config.GetAppConfig().Sqlite.GetDSN()
		_ = os.MkdirAll(filepath.Dir(dsn), 0777)
	default:
		return nil, errors.New("unknown db: " + cfg.DBType)
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt: true,
	})

	if err != nil {
		return nil, err
	}

	if cfg.Debug {
		db = db.Debug()
	}

	sqlDb, _ := db.DB()
	sqlDb.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDb.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDb.SetConnMaxLifetime(time.Duration(cfg.MaxLifetime) * time.Second)

	return db, nil
}

func autoMigrate(db *gorm.DB) error {
	// if dbType := config.GetAppConfig().Database.DBType; strings.ToLower(dbType) == "mysql" {
	// 	db = db.Set("gorm:table_options", "ENGINE=InnoDB")
	// }
	return db.AutoMigrate(
		&User{},
	)
}
