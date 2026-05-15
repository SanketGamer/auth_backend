package database

import (
	"fmt"
	"user_management_system/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

type Config struct{
	Host string
	Port string
	Password string
	User string
	DBName string
	SSLMode string
}

func NewConnection(cfg *Config) error {
    dsn := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=UTC",
        cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
    )
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    //log.Fatal(db)
    if err != nil {
        return err
    }
    
    DB = db
    //create tables automatically
    
	err = DB.AutoMigrate(&models.User{})
    if err != nil {
        return fmt.Errorf("failed to migrate database: %w", err)
    }
	fmt.Println("Database connected successfully")
    return nil
}