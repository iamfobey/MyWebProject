package application

import "database/sql"

type Application struct {
	Database *sql.DB
}

func NewApplication() *Application {
	db := setupDatabase()
	return &Application{Database: db}
}
