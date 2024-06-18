package main

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"time"
)

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

func NewDBConfig() *DBConfig {
	return &DBConfig{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     15450,  // 15450
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DBName:   os.Getenv("POSTGRES_DB"),
	}
}

func ConnectDB() (*gorm.DB, error) {
	dbConfig := NewDBConfig()
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=require TimeZone=Asia/Kolkata",
		dbConfig.Host, dbConfig.User, dbConfig.Password, dbConfig.DBName, dbConfig.Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	//Drop existing tables (if any)
  // db.Migrator().DropTable(&Account{}, &AccountType{}, &ChartOfAccount{}, &AccountBalance{}, &JournalEntry{})
	// AutoMigrate will create the necessary tables based on the models
	// db.AutoMigrate(&Account{}, &AccountType{}, &ChartOfAccount{}, &AccountBalance{}, &JournalEntry{})

	// db.AutoMigrate(&Account{})
	// db.AutoMigrate(&JournalEntry{})
	// db.AutoMigrate(&ChartOfAccount{})	
	// db.AutoMigrate(&AccountBalance{})
	// db.AutoMigrate(&AccountType{})

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

