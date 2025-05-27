package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	migrationpostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	"github.com/mcrosignani/uala_challenge/tweets/internal/config"
)

func NewPostgresDB(cfg config.Config) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.PostgresDB.Host,
		cfg.PostgresDB.Port,
		cfg.PostgresDB.User,
		cfg.PostgresDB.Password,
		cfg.PostgresDB.DBName,
	)

	var db *sql.DB
	var err error
	for i := 0; i < 10; i++ {
		db, err = sql.Open("postgres", connStr)
		if err == nil {
			err = db.Ping()
		}

		if err == nil {
			log.Println("Conectado a PostgreSQL!")
			break
		}

		log.Printf("Esperando PostgreSQL (%d/10)...", i+1)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatalf("No se pudo conectar a la base de datos: %v", err)
		return nil, err
	}

	return db, nil
}

func RunDBMigrations(db *sql.DB, filesPath string) error {
	driver, err := migrationpostgres.WithInstance(db, &migrationpostgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", filesPath),
		"postgres",
		driver,
	)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
