package postgres

import (
	"context"

	"github.com/arash-mosavi/go-base-repository/pkg/identifier"
	"github.com/arash-mosavi/go-base-repository/pkg/interfaces"
	"github.com/arash-mosavi/go-base-repository/pkg/types"
	postgresDomain "github.com/arash-mosavi/postgrs-unit-of-work-system/pkg/domain"
	postgresIdentifier "github.com/arash-mosavi/postgrs-unit-of-work-system/pkg/identifier"
	postgresUOW "github.com/arash-mosavi/postgrs-unit-of-work-system/pkg/persistence"
)

// BaseRepository implements the PostgreSQL base repository using composition
type BaseRepository[T types.PostgresEntity] struct {
	factory postgresUOW.IUnitOfWorkFactory[T]
}

// NewBaseRepository creates a new PostgreSQL base repository
func NewBaseRepository[T types.PostgresEntity](factory postgresUOW.IUnitOfWorkFactory[T]) interfaces.PostgresBaseRepository[T] {
	return &BaseRepository[T]{
		factory: factory,
	}
}

// FindOneById finds an entity by its PostgreSQL integer ID
func (r *BaseRepository[T]) FindOneById(ctx context.Context, id types.PostgresID) (T, error) {
	uow := r.factory.CreateWithContext(ctx)
	return uow.FindOneById(ctx, id)
}

// FindOne finds a single entity using identifier
func (r *BaseRepository[T]) FindOne(ctx context.Context, filter types.Identifier) (T, error) {
	uow := r.factory.CreateWithContext(ctx)
	unifiedFilter := filter.(*identifier.UnifiedIdentifier)
	return uow.FindOneByIdentifier(ctx, unifiedFilter.GetPostgresIdentifier())
}

// FindAll finds all entities matching the identifier
func (r *BaseRepository[T]) FindAll(ctx context.Context, filter types.Identifier) ([]T, error) {
	uow := r.factory.CreateWithContext(ctx)
	return uow.FindAll(ctx)
}

// FindAllWithPagination finds entities with pagination
func (r *BaseRepository[T]) FindAllWithPagination(ctx context.Context, params types.QueryParams[T]) ([]T, int64, error) {
	uow := r.factory.CreateWithContext(ctx)

	queryParams := postgresDomain.QueryParams[T]{
		Filter: params.Filter,
		Limit:  params.Limit,
		Offset: params.Offset,
		Sort:   convertSortMap(params.Sort),
	}

	entities, count, err := uow.FindAllWithPagination(ctx, queryParams)
	return entities, int64(count), err
}

// Insert creates a new entity
func (r *BaseRepository[T]) Insert(ctx context.Context, entity T) (T, error) {
	uow := r.factory.CreateWithContext(ctx)
	return uow.Insert(ctx, entity)
}

// Update modifies an existing entity
func (r *BaseRepository[T]) Update(ctx context.Context, filter types.Identifier, entity T) (T, error) {
	uow := r.factory.CreateWithContext(ctx)
	unifiedFilter := filter.(*identifier.UnifiedIdentifier)
	return uow.Update(ctx, unifiedFilter.GetPostgresIdentifier(), entity)
}

// Delete removes an entity
func (r *BaseRepository[T]) Delete(ctx context.Context, filter types.Identifier) error {
	uow := r.factory.CreateWithContext(ctx)
	unifiedFilter := filter.(*identifier.UnifiedIdentifier)
	return uow.Delete(ctx, unifiedFilter.GetPostgresIdentifier())
}

// BulkInsert creates multiple entities
func (r *BaseRepository[T]) BulkInsert(ctx context.Context, entities []T) ([]T, error) {
	uow := r.factory.CreateWithContext(ctx)
	return uow.BulkInsert(ctx, entities)
}

// BulkUpdate modifies multiple entities
func (r *BaseRepository[T]) BulkUpdate(ctx context.Context, entities []T) ([]T, error) {
	uow := r.factory.CreateWithContext(ctx)
	return uow.BulkUpdate(ctx, entities)
}

// BulkDelete removes multiple entities
func (r *BaseRepository[T]) BulkDelete(ctx context.Context, filters []types.Identifier) error {
	uow := r.factory.CreateWithContext(ctx)

	postgresFilters := make([]postgresIdentifier.IIdentifier, len(filters))
	for i, filter := range filters {
		unifiedFilter := filter.(*identifier.UnifiedIdentifier)
		postgresFilters[i] = unifiedFilter.GetPostgresIdentifier()
	}

	return uow.BulkHardDelete(ctx, postgresFilters)
}

// SoftDelete marks an entity as deleted
func (r *BaseRepository[T]) SoftDelete(ctx context.Context, filter types.Identifier) (T, error) {
	uow := r.factory.CreateWithContext(ctx)
	unifiedFilter := filter.(*identifier.UnifiedIdentifier)
	return uow.SoftDelete(ctx, unifiedFilter.GetPostgresIdentifier())
}

// HardDelete permanently removes an entity
func (r *BaseRepository[T]) HardDelete(ctx context.Context, filter types.Identifier) (T, error) {
	uow := r.factory.CreateWithContext(ctx)
	unifiedFilter := filter.(*identifier.UnifiedIdentifier)
	return uow.HardDelete(ctx, unifiedFilter.GetPostgresIdentifier())
}

// BulkSoftDelete marks multiple entities as deleted
func (r *BaseRepository[T]) BulkSoftDelete(ctx context.Context, filters []types.Identifier) error {
	uow := r.factory.CreateWithContext(ctx)

	postgresFilters := make([]postgresIdentifier.IIdentifier, len(filters))
	for i, filter := range filters {
		unifiedFilter := filter.(*identifier.UnifiedIdentifier)
		postgresFilters[i] = unifiedFilter.GetPostgresIdentifier()
	}

	return uow.BulkSoftDelete(ctx, postgresFilters)
}

// BulkHardDelete permanently removes multiple entities
func (r *BaseRepository[T]) BulkHardDelete(ctx context.Context, filters []types.Identifier) error {
	uow := r.factory.CreateWithContext(ctx)

	postgresFilters := make([]postgresIdentifier.IIdentifier, len(filters))
	for i, filter := range filters {
		unifiedFilter := filter.(*identifier.UnifiedIdentifier)
		postgresFilters[i] = unifiedFilter.GetPostgresIdentifier()
	}

	return uow.BulkHardDelete(ctx, postgresFilters)
}

// GetTrashed retrieves all soft-deleted entities
func (r *BaseRepository[T]) GetTrashed(ctx context.Context) ([]T, error) {
	uow := r.factory.CreateWithContext(ctx)
	return uow.GetTrashed(ctx)
}

// GetTrashedWithPagination retrieves soft-deleted entities with pagination
func (r *BaseRepository[T]) GetTrashedWithPagination(ctx context.Context, params types.QueryParams[T]) ([]T, int64, error) {
	uow := r.factory.CreateWithContext(ctx)

	queryParams := postgresDomain.QueryParams[T]{
		Filter: params.Filter,
		Limit:  params.Limit,
		Offset: params.Offset,
		Sort:   convertSortMap(params.Sort),
	}

	entities, count, err := uow.GetTrashedWithPagination(ctx, queryParams)
	return entities, int64(count), err
}

// Restore recovers a soft-deleted entity
func (r *BaseRepository[T]) Restore(ctx context.Context, filter types.Identifier) (T, error) {
	uow := r.factory.CreateWithContext(ctx)
	unifiedFilter := filter.(*identifier.UnifiedIdentifier)
	return uow.Restore(ctx, unifiedFilter.GetPostgresIdentifier())
}

// RestoreAll recovers all soft-deleted entities
func (r *BaseRepository[T]) RestoreAll(ctx context.Context) error {
	uow := r.factory.CreateWithContext(ctx)
	return uow.RestoreAll(ctx)
}

// BeginTransaction starts a database transaction
func (r *BaseRepository[T]) BeginTransaction(ctx context.Context) error {
	uow := r.factory.CreateWithContext(ctx)
	return uow.BeginTransaction(ctx)
}

// CommitTransaction commits the current transaction
func (r *BaseRepository[T]) CommitTransaction(ctx context.Context) error {
	uow := r.factory.CreateWithContext(ctx)
	return uow.CommitTransaction(ctx)
}

// RollbackTransaction rolls back the current transaction
func (r *BaseRepository[T]) RollbackTransaction(ctx context.Context) error {
	uow := r.factory.CreateWithContext(ctx)
	uow.RollbackTransaction(ctx)
	return nil
}

// convertSortMap converts unified sort map to PostgreSQL sort map
func convertSortMap(sort types.SortMap) postgresDomain.SortMap {
	if sort == nil {
		return nil
	}

	postgresSort := make(postgresDomain.SortMap)
	for field, direction := range sort {
		switch direction {
		case types.SortAsc:
			postgresSort[field] = postgresDomain.SortAsc
		case types.SortDesc:
			postgresSort[field] = postgresDomain.SortDesc
		}
	}

	return postgresSort
}
