package application

import (
	httpadapter "backend/internal/backend/adapters/http-adapter"
	postgresadapter "backend/internal/backend/adapters/postgres-adapter"
)

type App struct {
	Config Config
}

type Config struct {
	HttpAdapter     *httpadapter.Adapter
	PostgresAdapter *postgresadapter.Adapter
}

func New(config Config) *App {
	return &App{
		Config: config,
	}
}

func (app *App) Run() error {
	err := app.Config.HttpAdapter.Run()
	if err != nil {
		return err
	}
	return nil
}

func (app *App) Close() error {
	err := app.Config.PostgresAdapter.Close()
	if err != nil {
		return err
	}
	return nil
}
