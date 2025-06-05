package main

import (
	"context"
	"time"

	"github.com/arash-mosavi/go-base-repository/pkg/factory"
	"github.com/arash-mosavi/go-base-repository/pkg/identifier"
	"github.com/arash-mosavi/go-base-repository/pkg/interfaces"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gorm.io/gorm"
)

// MongoUser represents a MongoDB user entity
type MongoUser struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Email     string             `bson:"email" json:"email"`
	Slug      string             `bson:"slug" json:"slug"`
	Age       int                `bson:"age" json:"age"`
	Active    bool               `bson:"active" json:"active"`
	CreatedAt time.Time          `bson:"createdAt" json:"created_at"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updated_at"`
	DeletedAt *time.Time         `bson:"deletedAt,omitempty" json:"deleted_at,omitempty"`
}

// Implement MongoDB BaseModel interface
func (u *MongoUser) GetID() primitive.ObjectID   { return u.ID }
func (u *MongoUser) SetID(id primitive.ObjectID) { u.ID = id }
func (u *MongoUser) GetSlug() string             { return u.Slug }
func (u *MongoUser) SetSlug(slug string)         { u.Slug = slug }
func (u *MongoUser) GetName() string             { return u.Name }
func (u *MongoUser) SetName(name string)         { u.Name = name }
func (u *MongoUser) GetCreatedAt() time.Time     { return u.CreatedAt }
func (u *MongoUser) GetUpdatedAt() time.Time     { return u.UpdatedAt }
func (u *MongoUser) GetDeletedAt() *time.Time    { return u.DeletedAt }
func (u *MongoUser) SetCreatedAt(t time.Time)    { u.CreatedAt = t }
func (u *MongoUser) SetUpdatedAt(t time.Time)    { u.UpdatedAt = t }
func (u *MongoUser) SetDeletedAt(t *time.Time)   { u.DeletedAt = t }
func (u *MongoUser) IsDeleted() bool             { return u.DeletedAt != nil }

// PostgresUser represents a PostgreSQL user entity
type PostgresUser struct {
	ID        int            `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	Email     string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Slug      string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"slug"`
	Age       int            `gorm:"type:int;not null" json:"age"`
	Active    bool           `gorm:"default:true" json:"active"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

// Implement PostgreSQL BaseModel interface
func (u *PostgresUser) GetID() int                    { return u.ID }
func (u *PostgresUser) GetSlug() string               { return u.Slug }
func (u *PostgresUser) SetSlug(slug string)           { u.Slug = slug }
func (u *PostgresUser) GetCreatedAt() time.Time       { return u.CreatedAt }
func (u *PostgresUser) GetUpdatedAt() time.Time       { return u.UpdatedAt }
func (u *PostgresUser) GetArchivedAt() gorm.DeletedAt { return u.DeletedAt }
func (u *PostgresUser) GetName() string               { return u.Name }

func (PostgresUser) TableName() string { return "users" }

// MongoUserRepository demonstrates composition-based repository
type MongoUserRepository struct {
	baseRepo interfaces.MongoBaseRepository[*MongoUser]
}

// NewMongoUserRepository creates a new MongoDB user repository
func NewMongoUserRepository(baseRepo interfaces.MongoBaseRepository[*MongoUser]) *MongoUserRepository {
	return &MongoUserRepository{
		baseRepo: baseRepo,
	}
}

// CreateUser creates a new user
func (r *MongoUserRepository) CreateUser(ctx context.Context, name, email string, age int) (*MongoUser, error) {
	user := &MongoUser{
		Name:   name,
		Email:  email,
		Age:    age,
		Active: true,
	}
	return r.baseRepo.Insert(ctx, user)
}

// FindUserByEmail finds a user by email
func (r *MongoUserRepository) FindUserByEmail(ctx context.Context, email string) (*MongoUser, error) {
	filter := identifier.NewMongoIdentifier().Equal("email", email)
	return r.baseRepo.FindOne(ctx, filter)
}

// FindActiveUsers finds all active users
func (r *MongoUserRepository) FindActiveUsers(ctx context.Context) ([]*MongoUser, error) {
	filter := identifier.NewMongoIdentifier().Equal("active", true)
	return r.baseRepo.FindAll(ctx, filter)
}

// PostgresUserRepository demonstrates composition-based repository
type PostgresUserRepository struct {
	baseRepo interfaces.PostgresBaseRepository[*PostgresUser]
}

// NewPostgresUserRepository creates a new PostgreSQL user repository
func NewPostgresUserRepository(baseRepo interfaces.PostgresBaseRepository[*PostgresUser]) *PostgresUserRepository {
	return &PostgresUserRepository{
		baseRepo: baseRepo,
	}
}

// CreateUser creates a new user
func (r *PostgresUserRepository) CreateUser(ctx context.Context, name, email string, age int) (*PostgresUser, error) {
	user := &PostgresUser{
		Name:   name,
		Email:  email,
		Age:    age,
		Active: true,
	}
	return r.baseRepo.Insert(ctx, user)
}

// FindUserByEmail finds a user by email
func (r *PostgresUserRepository) FindUserByEmail(ctx context.Context, email string) (*PostgresUser, error) {
	filter := identifier.NewPostgresIdentifier().Equal("email", email)
	return r.baseRepo.FindOne(ctx, filter)
}

// FindActiveUsers finds all active users
func (r *PostgresUserRepository) FindActiveUsers(ctx context.Context) ([]*PostgresUser, error) {
	filter := identifier.NewPostgresIdentifier().Equal("active", true)
	return r.baseRepo.FindAll(ctx, filter)
}

// ExampleUsage demonstrates how to use the base repository SDK
func ExampleUsage() error {
	ctx := context.Background()

	// MongoDB Example
	mongoConfig := factory.NewMongoConfig()
	mongoConfig.Host = "localhost"
	mongoConfig.Port = 27017
	mongoConfig.Database = "testdb"

	mongoBaseRepo, err := factory.NewMongoBaseRepository[*MongoUser](mongoConfig)
	if err != nil {
		return err
	}

	mongoUserRepo := NewMongoUserRepository(mongoBaseRepo)

	// Create a user
	user, err := mongoUserRepo.CreateUser(ctx, "John Doe", "john@example.com", 30)
	if err != nil {
		return err
	}
	_ = user

	// Find user by email
	foundUser, err := mongoUserRepo.FindUserByEmail(ctx, "john@example.com")
	if err != nil {
		return err
	}
	_ = foundUser

	// PostgreSQL Example
	postgresConfig := factory.NewPostgresConfig()
	postgresConfig.Host = "localhost"
	postgresConfig.Port = 5432
	postgresConfig.User = "postgres"
	postgresConfig.Password = "password"
	postgresConfig.Database = "testdb"

	postgresBaseRepo, err := factory.NewPostgresBaseRepository[*PostgresUser](postgresConfig)
	if err != nil {
		return err
	}

	postgresUserRepo := NewPostgresUserRepository(postgresBaseRepo)

	// Create a user
	pgUser, err := postgresUserRepo.CreateUser(ctx, "Jane Doe", "jane@example.com", 25)
	if err != nil {
		return err
	}
	_ = pgUser

	return nil
}

func main() {
	// This is a demonstration file showing how to use the Base Repository SDK
	// In a real application, you would initialize your database connections
	// and use the repositories as shown in the example functions above

	// For now, this just demonstrates that all types compile correctly
	ctx := context.Background()
	_ = ctx

	// Example usage would be:
	// if err := demonstrateMongoUsage(); err != nil {
	//     log.Fatal(err)
	// }
	// if err := demonstratePostgresUsage(); err != nil {
	//     log.Fatal(err)
	// }
}
