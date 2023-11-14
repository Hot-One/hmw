package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"app/config"
	"app/storage"
)

type store struct {
	db       *sql.DB
	category *categoryRepo
	product  *productRepo
}

func NewConnectionPostgres(cfg *config.Config) (storage.StorageI, error) {

	connet := fmt.Sprintf(
		"host=%s user=%s dbname=%s password=%s port=%d sslmode=disable",
		cfg.PostgresHost,
		cfg.PostgresUser,
		cfg.PostgresDatabase,
		cfg.PostgresPassword,
		cfg.PostgresPort,
	)

	sqlDb, err := sql.Open("postgres", connet)
	if err != nil {
		return nil, err
	}

	if err := sqlDb.Ping(); err != nil {
		return nil, err
	}

	return &store{
		db: sqlDb,
	}, nil
}

func (s *store) Close() {
	s.db.Close()
}

func (s *store) Category() storage.CategoryRepoI {

	if s.category == nil {
		s.category = NewCategoryRepo(s.db)
	}

	return s.category
}

func (s *store) Product() storage.ProductRepoI {

	if s.product == nil {
		s.product = NewProductRepo(s.db)
	}

	return s.product
}
