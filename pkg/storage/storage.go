package storage

import "context"

type DB interface {
	Connect(ctx context.Context, dsn string) error
	Close(ctx context.Context)
}
