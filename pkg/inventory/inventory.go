package inventory

import "github.com/jackc/pgx/v4/pgxpool"

type Service struct {
	DBPool *pgxpool.Pool
}

func NewService(dbpool *pgxpool.Pool) *Service {
	return &Service{DBPool: dbpool}
}
