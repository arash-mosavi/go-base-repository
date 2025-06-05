package interfaces

import (
	"context"

	"github.com/arash-mosavi/go-base-repository/pkg/types"
)

// MongoBaseRepository defines the base repository interface for MongoDB entities
type MongoBaseRepository[T types.MongoEntity] interface {
	// Core CRUD operations
	FindOneById(ctx context.Context, id types.MongoID) (T, error)
	FindOne(ctx context.Context, filter types.Identifier) (T, error)
	FindAll(ctx context.Context, filter types.Identifier) ([]T, error)
	FindAllWithPagination(ctx context.Context, params types.QueryParams[T]) ([]T, int64, error)

	Insert(ctx context.Context, entity T) (T, error)
	Update(ctx context.Context, filter types.Identifier, entity T) (T, error)
	Delete(ctx context.Context, filter types.Identifier) error

	// Bulk operations
	BulkInsert(ctx context.Context, entities []T) ([]T, error)
	BulkUpdate(ctx context.Context, entities []T) ([]T, error)
	BulkDelete(ctx context.Context, filters []types.Identifier) error

	// Soft delete operations
	SoftDelete(ctx context.Context, filter types.Identifier) (T, error)
	HardDelete(ctx context.Context, filter types.Identifier) (T, error)
	BulkSoftDelete(ctx context.Context, filters []types.Identifier) error
	BulkHardDelete(ctx context.Context, filters []types.Identifier) error

	// Trashed data management
	GetTrashed(ctx context.Context) ([]T, error)
	GetTrashedWithPagination(ctx context.Context, params types.QueryParams[T]) ([]T, int64, error)
	Restore(ctx context.Context, filter types.Identifier) (T, error)
	RestoreAll(ctx context.Context) error

	// Transaction management
	BeginTransaction(ctx context.Context) error
	CommitTransaction(ctx context.Context) error
	RollbackTransaction(ctx context.Context) error
}

// PostgresBaseRepository defines the base repository interface for PostgreSQL entities
type PostgresBaseRepository[T types.PostgresEntity] interface {
	// Core CRUD operations
	FindOneById(ctx context.Context, id types.PostgresID) (T, error)
	FindOne(ctx context.Context, filter types.Identifier) (T, error)
	FindAll(ctx context.Context, filter types.Identifier) ([]T, error)
	FindAllWithPagination(ctx context.Context, params types.QueryParams[T]) ([]T, int64, error)

	Insert(ctx context.Context, entity T) (T, error)
	Update(ctx context.Context, filter types.Identifier, entity T) (T, error)
	Delete(ctx context.Context, filter types.Identifier) error

	// Bulk operations
	BulkInsert(ctx context.Context, entities []T) ([]T, error)
	BulkUpdate(ctx context.Context, entities []T) ([]T, error)
	BulkDelete(ctx context.Context, filters []types.Identifier) error

	// Soft delete operations
	SoftDelete(ctx context.Context, filter types.Identifier) (T, error)
	HardDelete(ctx context.Context, filter types.Identifier) (T, error)
	BulkSoftDelete(ctx context.Context, filters []types.Identifier) error
	BulkHardDelete(ctx context.Context, filters []types.Identifier) error

	// Trashed data management
	GetTrashed(ctx context.Context) ([]T, error)
	GetTrashedWithPagination(ctx context.Context, params types.QueryParams[T]) ([]T, int64, error)
	Restore(ctx context.Context, filter types.Identifier) (T, error)
	RestoreAll(ctx context.Context) error

	// Transaction management
	BeginTransaction(ctx context.Context) error
	CommitTransaction(ctx context.Context) error
	RollbackTransaction(ctx context.Context) error
}
