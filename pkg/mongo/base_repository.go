package mongo

import (
	"context"

	"github.com/arash-mosavi/go-base-repository/pkg/identifier"
	"github.com/arash-mosavi/go-base-repository/pkg/interfaces"
	"github.com/arash-mosavi/go-base-repository/pkg/types"
	mongoDomain "github.com/arash-mosavi/mongo-unit-of-work-system/pkg/domain"
	mongoIdentifier "github.com/arash-mosavi/mongo-unit-of-work-system/pkg/identifier"
	mongoUOW "github.com/arash-mosavi/mongo-unit-of-work-system/pkg/persistence"
)

// BaseRepository implements the MongoDB base repository using composition
type BaseRepository[T types.MongoEntity] struct {
	factory mongoUOW.IUnitOfWorkFactory[T]
}

// NewBaseRepository creates a new MongoDB base repository
func NewBaseRepository[T types.MongoEntity](factory mongoUOW.IUnitOfWorkFactory[T]) interfaces.MongoBaseRepository[T] {
	return &BaseRepository[T]{
		factory: factory,
	}
}

// FindOneById finds an entity by its MongoDB ObjectID
func (r *BaseRepository[T]) FindOneById(ctx context.Context, id types.MongoID) (T, error) {
	uow := r.factory.CreateWithContext(ctx)
	return uow.FindOneById(ctx, id)
}

// FindOne finds a single entity using identifier
func (r *BaseRepository[T]) FindOne(ctx context.Context, filter types.Identifier) (T, error) {
	uow := r.factory.CreateWithContext(ctx)
	unifiedFilter := filter.(*identifier.UnifiedIdentifier)
	return uow.FindOneByIdentifier(ctx, unifiedFilter.GetMongoIdentifier())
}

// FindAll finds all entities matching the identifier
func (r *BaseRepository[T]) FindAll(ctx context.Context, filter types.Identifier) ([]T, error) {
	uow := r.factory.CreateWithContext(ctx)
	return uow.FindAll(ctx)
}

// FindAllWithPagination finds entities with pagination
func (r *BaseRepository[T]) FindAllWithPagination(ctx context.Context, params types.QueryParams[T]) ([]T, int64, error) {
	uow := r.factory.CreateWithContext(ctx)

	queryParams := mongoDomain.QueryParams[T]{
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
	return uow.Update(ctx, unifiedFilter.GetMongoIdentifier(), entity)
}

// Delete removes an entity
func (r *BaseRepository[T]) Delete(ctx context.Context, filter types.Identifier) error {
	uow := r.factory.CreateWithContext(ctx)
	unifiedFilter := filter.(*identifier.UnifiedIdentifier)
	return uow.Delete(ctx, unifiedFilter.GetMongoIdentifier())
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

	// Convert unified filters to native MongoDB identifiers
	mongoFilters := make([]mongoIdentifier.IIdentifier, len(filters))
	for i, filter := range filters {
		unifiedFilter := filter.(*identifier.UnifiedIdentifier)
		mongoFilters[i] = unifiedFilter.GetMongoIdentifier()
	}

	return uow.BulkHardDelete(ctx, mongoFilters)
}

// SoftDelete marks an entity as deleted
func (r *BaseRepository[T]) SoftDelete(ctx context.Context, filter types.Identifier) (T, error) {
	uow := r.factory.CreateWithContext(ctx)
	unifiedFilter := filter.(*identifier.UnifiedIdentifier)
	return uow.SoftDelete(ctx, unifiedFilter.GetMongoIdentifier())
}

// HardDelete permanently removes an entity
func (r *BaseRepository[T]) HardDelete(ctx context.Context, filter types.Identifier) (T, error) {
	uow := r.factory.CreateWithContext(ctx)
	unifiedFilter := filter.(*identifier.UnifiedIdentifier)
	return uow.HardDelete(ctx, unifiedFilter.GetMongoIdentifier())
}

// BulkSoftDelete marks multiple entities as deleted
func (r *BaseRepository[T]) BulkSoftDelete(ctx context.Context, filters []types.Identifier) error {
	uow := r.factory.CreateWithContext(ctx)

	mongoFilters := make([]mongoIdentifier.IIdentifier, len(filters))
	for i, filter := range filters {
		unifiedFilter := filter.(*identifier.UnifiedIdentifier)
		mongoFilters[i] = unifiedFilter.GetMongoIdentifier()
	}

	return uow.BulkSoftDelete(ctx, mongoFilters)
}

// BulkHardDelete permanently removes multiple entities
func (r *BaseRepository[T]) BulkHardDelete(ctx context.Context, filters []types.Identifier) error {
	uow := r.factory.CreateWithContext(ctx)

	mongoFilters := make([]mongoIdentifier.IIdentifier, len(filters))
	for i, filter := range filters {
		unifiedFilter := filter.(*identifier.UnifiedIdentifier)
		mongoFilters[i] = unifiedFilter.GetMongoIdentifier()
	}

	return uow.BulkHardDelete(ctx, mongoFilters)
}

// GetTrashed retrieves all soft-deleted entities
func (r *BaseRepository[T]) GetTrashed(ctx context.Context) ([]T, error) {
	uow := r.factory.CreateWithContext(ctx)
	return uow.GetTrashed(ctx)
}

// GetTrashedWithPagination retrieves soft-deleted entities with pagination
func (r *BaseRepository[T]) GetTrashedWithPagination(ctx context.Context, params types.QueryParams[T]) ([]T, int64, error) {
	uow := r.factory.CreateWithContext(ctx)

	queryParams := mongoDomain.QueryParams[T]{
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
	return uow.Restore(ctx, unifiedFilter.GetMongoIdentifier())
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

// convertSortMap converts unified sort map to MongoDB sort map
func convertSortMap(sort types.SortMap) mongoDomain.SortMap {
	if sort == nil {
		return nil
	}

	mongoSort := make(mongoDomain.SortMap)
	for field, direction := range sort {
		switch direction {
		case types.SortAsc:
			mongoSort[field] = mongoDomain.SortAsc
		case types.SortDesc:
			mongoSort[field] = mongoDomain.SortDesc
		}
	}

	return mongoSort
}
