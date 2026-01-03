package database

import (
	"context"
	"fmt"

	"github.com/yurin-kami/PackChann/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewConnection(ctx context.Context, cfg models.DBConfig) (*gorm.DB, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto Migrate the schema
	err = db.AutoMigrate(&models.User{}, &models.Pack{}, &models.Notice{}, &models.UserToken{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
