package interfaces_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/arash-mosavi/go-base-repository/pkg/identifier"
	"github.com/arash-mosavi/go-base-repository/pkg/interfaces"
	"github.com/arash-mosavi/go-base-repository/pkg/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gorm.io/gorm"
)

// MockMongoEntity for testing MongoDB repositories
type MockMongoEntity struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Email     string             `bson:"email" json:"email"`
	Slug      string             `bson:"slug" json:"slug"`
	CreatedAt time.Time          `bson:"createdAt" json:"created_at"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updated_at"`
	DeletedAt *time.Time         `bson:"deletedAt,omitempty" json:"deleted_at,omitempty"`
}

// Implement MongoDB BaseModel interface
func (m *MockMongoEntity) GetID() primitive.ObjectID   { return m.ID }
func (m *MockMongoEntity) SetID(id primitive.ObjectID) { m.ID = id }
func (m *MockMongoEntity) GetSlug() string             { return m.Slug }
func (m *MockMongoEntity) SetSlug(slug string)         { m.Slug = slug }
func (m *MockMongoEntity) GetName() string             { return m.Name }
func (m *MockMongoEntity) SetName(name string)         { m.Name = name }
func (m *MockMongoEntity) GetCreatedAt() time.Time     { return m.CreatedAt }
func (m *MockMongoEntity) GetUpdatedAt() time.Time     { return m.UpdatedAt }
func (m *MockMongoEntity) GetDeletedAt() *time.Time    { return m.DeletedAt }
func (m *MockMongoEntity) SetCreatedAt(t time.Time)    { m.CreatedAt = t }
func (m *MockMongoEntity) SetUpdatedAt(t time.Time)    { m.UpdatedAt = t }
func (m *MockMongoEntity) SetDeletedAt(t *time.Time)   { m.DeletedAt = t }
func (m *MockMongoEntity) IsDeleted() bool             { return m.DeletedAt != nil }

// MockPostgresEntity for testing PostgreSQL repositories
type MockPostgresEntity struct {
	ID        int            `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	Email     string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Slug      string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"slug"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

// Implement PostgreSQL BaseModel interface
func (m *MockPostgresEntity) GetID() int                    { return m.ID }
func (m *MockPostgresEntity) GetSlug() string               { return m.Slug }
func (m *MockPostgresEntity) SetSlug(slug string)           { m.Slug = slug }
func (m *MockPostgresEntity) GetCreatedAt() time.Time       { return m.CreatedAt }
func (m *MockPostgresEntity) GetUpdatedAt() time.Time       { return m.UpdatedAt }
func (m *MockPostgresEntity) GetArchivedAt() gorm.DeletedAt { return m.DeletedAt }
func (m *MockPostgresEntity) GetName() string               { return m.Name }

func (MockPostgresEntity) TableName() string { return "mock_entities" }

// MockMongoRepository implements MongoBaseRepository for testing
type MockMongoRepository struct {
	entities map[primitive.ObjectID]*MockMongoEntity
}

func NewMockMongoRepository() *MockMongoRepository {
	return &MockMongoRepository{
		entities: make(map[primitive.ObjectID]*MockMongoEntity),
	}
}

func (m *MockMongoRepository) FindOneById(ctx context.Context, id types.MongoID) (*MockMongoEntity, error) {
	entity, exists := m.entities[id]
	if !exists {
		return nil, errors.New("entity not found")
	}
	return entity, nil
}

func (m *MockMongoRepository) FindOne(ctx context.Context, filter types.Identifier) (*MockMongoEntity, error) {
	// Simple mock implementation
	for _, entity := range m.entities {
		return entity, nil // Return first entity for simplicity
	}
	return nil, errors.New("entity not found")
}

func (m *MockMongoRepository) FindAll(ctx context.Context, filter types.Identifier) ([]*MockMongoEntity, error) {
	var entities []*MockMongoEntity
	for _, entity := range m.entities {
		entities = append(entities, entity)
	}
	return entities, nil
}

func (m *MockMongoRepository) FindAllWithPagination(ctx context.Context, params types.QueryParams[*MockMongoEntity]) ([]*MockMongoEntity, int64, error) {
	entities, err := m.FindAll(ctx, nil)
	if err != nil {
		return nil, 0, err
	}
	return entities, int64(len(entities)), nil
}

func (m *MockMongoRepository) Insert(ctx context.Context, entity *MockMongoEntity) (*MockMongoEntity, error) {
	if entity.ID == primitive.NilObjectID {
		entity.ID = primitive.NewObjectID()
	}
	entity.CreatedAt = time.Now()
	entity.UpdatedAt = time.Now()
	m.entities[entity.ID] = entity
	return entity, nil
}

func (m *MockMongoRepository) Update(ctx context.Context, filter types.Identifier, entity *MockMongoEntity) (*MockMongoEntity, error) {
	entity.UpdatedAt = time.Now()
	m.entities[entity.ID] = entity
	return entity, nil
}

func (m *MockMongoRepository) Delete(ctx context.Context, filter types.Identifier) error {
	// Simple mock implementation - delete first entity
	for id := range m.entities {
		delete(m.entities, id)
		return nil
	}
	return errors.New("entity not found")
}

func (m *MockMongoRepository) BulkInsert(ctx context.Context, entities []*MockMongoEntity) ([]*MockMongoEntity, error) {
	for _, entity := range entities {
		if entity.ID == primitive.NilObjectID {
			entity.ID = primitive.NewObjectID()
		}
		entity.CreatedAt = time.Now()
		entity.UpdatedAt = time.Now()
		m.entities[entity.ID] = entity
	}
	return entities, nil
}

func (m *MockMongoRepository) BulkUpdate(ctx context.Context, entities []*MockMongoEntity) ([]*MockMongoEntity, error) {
	for _, entity := range entities {
		entity.UpdatedAt = time.Now()
		m.entities[entity.ID] = entity
	}
	return entities, nil
}

func (m *MockMongoRepository) BulkDelete(ctx context.Context, filters []types.Identifier) error {
	// Simple mock implementation
	for id := range m.entities {
		delete(m.entities, id)
	}
	return nil
}

func (m *MockMongoRepository) SoftDelete(ctx context.Context, filter types.Identifier) (*MockMongoEntity, error) {
	for _, entity := range m.entities {
		now := time.Now()
		entity.DeletedAt = &now
		return entity, nil
	}
	return nil, errors.New("entity not found")
}

func (m *MockMongoRepository) HardDelete(ctx context.Context, filter types.Identifier) (*MockMongoEntity, error) {
	for id, entity := range m.entities {
		delete(m.entities, id)
		return entity, nil
	}
	return nil, errors.New("entity not found")
}

func (m *MockMongoRepository) BulkSoftDelete(ctx context.Context, filters []types.Identifier) error {
	now := time.Now()
	for _, entity := range m.entities {
		entity.DeletedAt = &now
	}
	return nil
}

func (m *MockMongoRepository) BulkHardDelete(ctx context.Context, filters []types.Identifier) error {
	for id := range m.entities {
		delete(m.entities, id)
	}
	return nil
}

func (m *MockMongoRepository) GetTrashed(ctx context.Context) ([]*MockMongoEntity, error) {
	var trashed []*MockMongoEntity
	for _, entity := range m.entities {
		if entity.DeletedAt != nil {
			trashed = append(trashed, entity)
		}
	}
	return trashed, nil
}

func (m *MockMongoRepository) GetTrashedWithPagination(ctx context.Context, params types.QueryParams[*MockMongoEntity]) ([]*MockMongoEntity, int64, error) {
	trashed, err := m.GetTrashed(ctx)
	if err != nil {
		return nil, 0, err
	}
	return trashed, int64(len(trashed)), nil
}

func (m *MockMongoRepository) Restore(ctx context.Context, filter types.Identifier) (*MockMongoEntity, error) {
	for _, entity := range m.entities {
		if entity.DeletedAt != nil {
			entity.DeletedAt = nil
			return entity, nil
		}
	}
	return nil, errors.New("entity not found")
}

func (m *MockMongoRepository) RestoreAll(ctx context.Context) error {
	for _, entity := range m.entities {
		entity.DeletedAt = nil
	}
	return nil
}

func (m *MockMongoRepository) BeginTransaction(ctx context.Context) error {
	return nil // Mock implementation
}

func (m *MockMongoRepository) CommitTransaction(ctx context.Context) error {
	return nil // Mock implementation
}

func (m *MockMongoRepository) RollbackTransaction(ctx context.Context) error {
	return nil // Mock implementation
}

// Test functions
func TestMongoBaseRepository_Interface(t *testing.T) {
	var repo interfaces.MongoBaseRepository[*MockMongoEntity]
	repo = NewMockMongoRepository()

	if repo == nil {
		t.Error("MockMongoRepository should implement MongoBaseRepository interface")
	}
}

func TestMongoBaseRepository_CRUD(t *testing.T) {
	repo := NewMockMongoRepository()
	ctx := context.Background()

	// Test Insert
	entity := &MockMongoEntity{
		Name:  "Test Entity",
		Email: "test@example.com",
		Slug:  "test-entity",
	}

	created, err := repo.Insert(ctx, entity)
	if err != nil {
		t.Fatalf("Failed to insert entity: %v", err)
	}

	if created.ID == primitive.NilObjectID {
		t.Error("Entity ID should be set after insert")
	}

	// Test FindOneById
	found, err := repo.FindOneById(ctx, created.ID)
	if err != nil {
		t.Fatalf("Failed to find entity by ID: %v", err)
	}

	if found.Name != "Test Entity" {
		t.Errorf("Expected name 'Test Entity', got '%s'", found.Name)
	}

	// Test FindAll
	all, err := repo.FindAll(ctx, identifier.NewMongoIdentifier())
	if err != nil {
		t.Fatalf("Failed to find all entities: %v", err)
	}

	if len(all) != 1 {
		t.Errorf("Expected 1 entity, got %d", len(all))
	}

	// Test Update
	entity.Name = "Updated Entity"
	updated, err := repo.Update(ctx, identifier.NewMongoIdentifier(), entity)
	if err != nil {
		t.Fatalf("Failed to update entity: %v", err)
	}

	if updated.Name != "Updated Entity" {
		t.Errorf("Expected name 'Updated Entity', got '%s'", updated.Name)
	}

	// Test SoftDelete
	softDeleted, err := repo.SoftDelete(ctx, identifier.NewMongoIdentifier())
	if err != nil {
		t.Fatalf("Failed to soft delete entity: %v", err)
	}

	if softDeleted.DeletedAt == nil {
		t.Error("Entity should have DeletedAt set after soft delete")
	}
}

func TestMongoBaseRepository_BulkOperations(t *testing.T) {
	repo := NewMockMongoRepository()
	ctx := context.Background()

	entities := []*MockMongoEntity{
		{Name: "Entity 1", Email: "entity1@example.com", Slug: "entity-1"},
		{Name: "Entity 2", Email: "entity2@example.com", Slug: "entity-2"},
	}

	// Test BulkInsert
	created, err := repo.BulkInsert(ctx, entities)
	if err != nil {
		t.Fatalf("Failed to bulk insert entities: %v", err)
	}

	if len(created) != 2 {
		t.Errorf("Expected 2 entities, got %d", len(created))
	}

	// Test BulkUpdate
	for _, entity := range created {
		entity.Name = "Updated " + entity.Name
	}

	updated, err := repo.BulkUpdate(ctx, created)
	if err != nil {
		t.Fatalf("Failed to bulk update entities: %v", err)
	}

	if len(updated) != 2 {
		t.Errorf("Expected 2 entities, got %d", len(updated))
	}
}

func TestMongoBaseRepository_TransactionHandling(t *testing.T) {
	repo := NewMockMongoRepository()
	ctx := context.Background()

	// Test transaction methods
	err := repo.BeginTransaction(ctx)
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}

	err = repo.CommitTransaction(ctx)
	if err != nil {
		t.Fatalf("Failed to commit transaction: %v", err)
	}

	err = repo.BeginTransaction(ctx)
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}

	err = repo.RollbackTransaction(ctx)
	if err != nil {
		t.Fatalf("Failed to rollback transaction: %v", err)
	}
}
