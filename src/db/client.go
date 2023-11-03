package db

import (
	"context"
	"log"
	"nak-auth/models"
	"os"

	"go.uber.org/fx"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewPScaleClient(lc fx.Lifecycle) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(os.Getenv("DSN")), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			db.AutoMigrate(&models.Code{}, &models.User{}, &models.Client{})
			return nil
		},
	})
	return db, err
}

func NewLibSqlClient(lc fx.Lifecycle) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("main.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			db.AutoMigrate(&models.Code{}, &models.User{}, &models.Client{})
			return nil
		},
	})
	return db, err
}
