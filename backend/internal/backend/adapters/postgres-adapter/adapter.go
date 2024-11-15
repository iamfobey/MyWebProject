package postgres_adapter

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Adapter struct {
	Database *sql.DB
	Config   Config
}

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DbName   string
}

func New(config Config) *Adapter {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.DbName,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Minute * 10)

	timeout := time.After(5 * time.Minute)
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			log.Fatal("Не удалось подключиться к серверу PostgreSQL в течение 5 минут")
		case <-ticker.C:
			if err := db.Ping(); err == nil {
				log.Println("Успешное подключение к серверу PostgreSQL")
				return &Adapter{
					Database: db,
					Config:   config,
				}
			} else {
				log.Println("Ожидание подключения к серверу PostgreSQL...")
			}
		}
	}
}

func (adapter *Adapter) Close() error {
	err := adapter.Database.Close()
	if err != nil {
		return err
	}
	return nil
}
