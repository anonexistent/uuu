package lib

import (
	"log"
	"os"
	"time"

	"gorm.io/gorm/logger"
)

var MyLogger = logger.New(
	log.New(os.Stderr, "\r\n", log.LstdFlags),
	logger.Config{
		SlowThreshold:             time.Millisecond,
		LogLevel:                  logger.Info,
		IgnoreRecordNotFoundError: true,
		Colorful:                  true,
	},
)

type User struct {
	Id    uint
	Name  string
	Email string
}

func (User) TableName() string {
	return "users"
}
