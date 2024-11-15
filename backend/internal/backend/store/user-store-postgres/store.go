package user_store_postgres

import (
	"backend/internal/backend/domain/user"
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type Store struct {
	Config Config
}

type Config struct {
	Database      *sql.DB
	AdminLogin    string
	AdminPassword string
}

func New(config Config) *Store {
	_, err := config.Database.Exec(`CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		login TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL
	);`)

	if err != nil {
		log.Fatal(err)
	}

	adminHashedPassword, err := bcrypt.GenerateFromPassword([]byte(config.AdminPassword), bcrypt.MinCost+4)

	_, err = config.Database.Exec(`INSERT INTO users (login, password) VALUES ($1, $2) ON CONFLICT (login) DO NOTHING`,
		config.AdminLogin, adminHashedPassword)
	if err != nil {
		log.Fatal(err)
	}

	return &Store{
		Config: config,
	}
}

func (store *Store) CreateUser(user *user.User) error {
	query := `INSERT INTO users (login, password) VALUES ($1, $2)`
	_, err := store.Config.Database.Exec(query, user.Credentials.Login, user.Credentials.Password)
	return err
}

func (store *Store) GetUserByLogin(login string) (*user.User, error) {
	query := `SELECT login, password FROM users WHERE login=$1`

	var usr user.User

	err := store.Config.Database.QueryRow(query, login).Scan(&usr.Credentials.Login, &usr.Credentials.Password)

	if err != nil {
		return nil, err
	}

	return &usr, nil
}
