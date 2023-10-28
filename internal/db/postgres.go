package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"mail_system/internal/config"

	pgx "github.com/jackc/pgx/v5/pgxpool"
)

// "context"

type PostgresDB struct {
	connPool *pgx.Pool
	cfg      *config.ConfigDb
}

func NewDb(ctx context.Context) (db *PostgresDB, err error) {
	db = new(PostgresDB)

	db.cfg = &config.ConfigDb{
		Host:         os.Getenv("DB_HOST"),
		Port:         os.Getenv("DB_PORT"),
		Pass:         os.Getenv("DB_PASSWORD"),
		User:         os.Getenv("DB_USER"),
		DbName:       os.Getenv("DB_NAME"),
		SSLMode:      os.Getenv("DB_SSL_MODE"),
		MaxPoolConns: os.Getenv("DB_MAX_CONN_POOLS"),
	}

	url := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=%s pool_max_conns=%s",
		db.cfg.User,
		db.cfg.Pass,
		db.cfg.Host,
		db.cfg.Port,
		db.cfg.DbName,
		db.cfg.SSLMode,
		db.cfg.MaxPoolConns,
	)

	db.connPool, err = pgx.New(ctx, url)

	if err != nil {
		log.Fatalf("Failed connection to Postgres: %s", err)
		return db, err
	}

	return db, err
}

// func (db *PostgresDB) CreateTables(ctx context.Context) (err error) {
// 	context, close := context.WithTimeout(ctx, 3*time.Second)
// 	defer close()

// 	db.connPool.Exec(context, )

// 	return err
// }
