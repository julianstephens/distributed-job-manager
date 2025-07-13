package database

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/julianstephens/distributed-task-scheduler/backend/internal/config"
	appLogger "github.com/julianstephens/distributed-task-scheduler/backend/pkg/logger"
)

var (
	DB   *gorm.DB
	err  error
	once sync.Once
)

type Database struct {
	*gorm.DB
}

var conf = config.GetConfig()

func GetDB() *gorm.DB {
	once.Do(func() {
		DB = setup()
		if err != nil {
			appLogger.Fatalf("unable to initialize database: %+v", err)
		}
	})
	return DB
}

func setup() *gorm.DB {
	var db *gorm.DB
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", conf.Database.Host, conf.Database.User, conf.Database.Password, conf.Database.DB, conf.Database.Port)

	dbLogger := logger.New(
		log.New(getWriter(), "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level (Silent, Error, Warn, Info)
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,       // Disable color
		},
	)

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: dbLogger,
	})

	if err != nil {
		return nil
	}

	return db
}

func getWriter() io.Writer {
	file, err := os.OpenFile("ls.db.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return os.Stdout
	} else {
		return file
	}
}
