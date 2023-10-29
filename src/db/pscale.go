package db

import (
	"context"
	"log"
	"nak-auth/services"
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
			db.AutoMigrate(&services.Code{}, &services.User{}, &services.Client{})
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
			db.AutoMigrate(&services.Code{}, &services.User{}, &services.Client{})
			return nil
		},
	})
	return db, err
}
