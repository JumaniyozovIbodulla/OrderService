package postgres

import (
	"context"
	"fmt"
	"order/config"
	"order/storage"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Store struct {
	db    *pgxpool.Pool
	order storage.OrderRepo
	orderProduct storage.OrderProductsRepo
	orderNotes storage.OrderNotesRepo
}

func NewPostgres(ctx context.Context, cfg config.Config) (storage.IStorage, error) {
	config, err := pgxpool.ParseConfig(fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresDatabase,
	))
	if err != nil {
		return nil, err
	}

	config.MaxConns = cfg.PostgresMaxConnections

	pool, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	return &Store{
		db: pool,
	}, nil
}

func (s *Store) CloseDB() {
	s.db.Close()
}

func (l *Store) Log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {
	args := make([]interface{}, 0, len(data)+2)
	args = append(args, level, msg)

	for k, v := range data {
		args = append(args, fmt.Sprintf("%s=%v", k, v))
	}

	fmt.Println(args...)
}

func (s *Store) Order() storage.OrderRepo {
	if s.order == nil {
		s.order = NewOrderRepo(s.db)
	}
	return s.order
}

func (s *Store) OrderProduct() storage.OrderProductsRepo {
	if s.orderProduct == nil {
		s.orderProduct = NewOrderProductsRepo(s.db)
	}
	return s.orderProduct
}

func (s *Store) OrderNotes() storage.OrderNotesRepo {
	if s.orderNotes == nil {
		s.orderNotes = NewOrderNotesRepo(s.db)
	}
	return s.orderNotes
}