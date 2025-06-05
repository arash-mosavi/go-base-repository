# Go Base Repository SDK

A comprehensive Base Repository SDK in Go that uses composition (not inheritance) to wrap either MongoDB or PostgreSQL Unit of Work systems. Built with Go 1.18+ generics for type safety and following clean architecture principles.

## Features

- **Composition-based Architecture**: Uses composition instead of inheritance for better flexibility
- **Dependency Injection**: BaseRepository is injected into concrete repositories
- **Full Transaction Support**: Complete transaction/session handling for both databases
- **Type Safety**: Go 1.18+ generics ensure compile-time type safety
- **Clean Architecture**: Follows clean architecture principles with idiomatic Go code
- **Dual Database Support**: Works with both MongoDB and PostgreSQL Unit of Work systems
- **Separate Type Hierarchies**: Maintains type safety with different ID types (ObjectID vs int)

## Architecture

The SDK is designed with separate type-safe interfaces for MongoDB and PostgreSQL to handle their fundamental differences:

- **MongoDB**: Uses `primitive.ObjectID` for primary keys
- **PostgreSQL**: Uses `int` for primary keys

This approach ensures type safety while maintaining the benefits of a unified base repository pattern.

## Installation

```bash
go get github.com/arash-mosavi/go-base-repository
```

### Dependencies

The SDK depends on these Unit of Work systems:

```bash
go get github.com/arash-mosavi/mongo-unit-of-work-system
go get github.com/arash-mosavi/postgrs-unit-of-work-system
```

## Quick Start

### MongoDB Usage

```go
package main

import (
    "context"
    "time"
    
    "go.mongodb.org/mongo-driver/bson/primitive"
    mongoUOW "github.com/arash-mosavi/mongo-unit-of-work-system/pkg/persistence"
    
    "github.com/arash-mosavi/go-base-repository/pkg/domain"
    "github.com/arash-mosavi/go-base-repository/pkg/repository"
)

// Create a concrete repository using composition
type UserRepository struct {
    *repository.MongoBaseRepository[domain.MongoUser]
}

func NewUserRepository(uow mongoUOW.IUnitOfWork[domain.MongoUser]) *UserRepository {
    return &UserRepository{
        MongoBaseRepository: repository.NewMongoBaseRepository(uow), // Dependency Injection
    }
}

func main() {
    // Initialize MongoDB UoW (your implementation)
    mongoUoW := mongoUOW.NewUnitOfWork[domain.MongoUser](client, "database", "users")
    
    // Create repository with injected base repository
    userRepo := NewUserRepository(mongoUoW)
    
    ctx := context.Background()
    
    // Use base repository functionality
    user := domain.MongoUser{
        ID:        primitive.NewObjectID(),
        Name:      "John Doe",
        Email:     "john@example.com",
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    
    // Create user
    createdUser, err := userRepo.Create(ctx, user)
    if err != nil {
        panic(err)
    }
    
    // Find user by ID
    foundUser, err := userRepo.FindByID(ctx, createdUser.GetID())
    if err != nil {
        panic(err)
    }
    
    // Find all users with pagination
    params := domain.QueryParams{
        Page:          1,
        PageSize:      10,
        SortBy:        "created_at",
        SortDirection: domain.SortDirectionDesc,
    }
    users, err := userRepo.FindAll(ctx, params)
    if err != nil {
        panic(err)
    }
}
```

### PostgreSQL Usage

```go
package main

import (
    "context"
    "time"
    
    postgresUOW "github.com/arash-mosavi/postgrs-unit-of-work-system/pkg/persistence"
    
    "github.com/arash-mosavi/go-base-repository/pkg/domain"
    "github.com/arash-mosavi/go-base-repository/pkg/repository"
)

// Create a concrete repository using composition
type UserRepository struct {
    *repository.PostgresBaseRepository[domain.PostgresUser]
}

func NewUserRepository(uow postgresUOW.IUnitOfWork[domain.PostgresUser]) *UserRepository {
    return &UserRepository{
        PostgresBaseRepository: repository.NewPostgresBaseRepository(uow), // Dependency Injection
    }
}

func main() {
    // Initialize PostgreSQL UoW (your implementation)
    postgresUoW := postgresUOW.NewUnitOfWork[domain.PostgresUser](db, "users")
    
    // Create repository with injected base repository
    userRepo := NewUserRepository(postgresUoW)
    
    ctx := context.Background()
    
    // Use base repository functionality
    user := domain.PostgresUser{
        Name:      "Jane Smith",
        Email:     "jane@example.com",
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    
    // Create user
    createdUser, err := userRepo.Create(ctx, user)
    if err != nil {
        panic(err)
    }
    
    // Find user by ID
    foundUser, err := userRepo.FindByID(ctx, createdUser.GetID())
    if err != nil {
        panic(err)
    }
}
```

## Model Constraints

### MongoDB Models

Your MongoDB models must implement `MongoModelConstraint`:

```go
type MongoUser struct {
    ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Name      string             `bson:"name" json:"name"`
    Email     string             `bson:"email" json:"email"`
    CreatedAt time.Time          `bson:"created_at" json:"created_at"`
    UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

func (u MongoUser) GetID() primitive.ObjectID { return u.ID }
func (u *MongoUser) SetID(id primitive.ObjectID) { u.ID = id }
func (u MongoUser) GetCreatedAt() time.Time { return u.CreatedAt }
func (u *MongoUser) SetCreatedAt(t time.Time) { u.CreatedAt = t }
func (u MongoUser) GetUpdatedAt() time.Time { return u.UpdatedAt }
func (u *MongoUser) SetUpdatedAt(t time.Time) { u.UpdatedAt = t }
```

### PostgreSQL Models

Your PostgreSQL models must implement `PostgresModelConstraint`:

```go
type PostgresUser struct {
    ID        int       `json:"id" gorm:"primaryKey"`
    Name      string    `json:"name" gorm:"column:name"`
    Email     string    `json:"email" gorm:"column:email"`
    CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
    UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (u PostgresUser) GetID() int { return u.ID }
func (u *PostgresUser) SetID(id int) { u.ID = id }
func (u PostgresUser) GetCreatedAt() time.Time { return u.CreatedAt }
func (u *PostgresUser) SetCreatedAt(t time.Time) { u.CreatedAt = t }
func (u PostgresUser) GetUpdatedAt() time.Time { return u.UpdatedAt }
func (u *PostgresUser) SetUpdatedAt(t time.Time) { u.UpdatedAt = t }

func (PostgresUser) TableName() string { return "users" }
```

## Transaction Handling

Both base repositories support full transaction handling:

```go
// MongoDB Transaction Example
ctx := context.Background()

// Begin transaction
txRepo, err := userRepo.BeginTransaction(ctx)
if err != nil {
    return err
}

// Perform operations
user1, err := txRepo.Create(ctx, newUser1)
if err != nil {
    txRepo.RollbackTransaction(ctx)
    return err
}

user2, err := txRepo.Create(ctx, newUser2)
if err != nil {
    txRepo.RollbackTransaction(ctx)
    return err
}

// Commit transaction
err = txRepo.CommitTransaction(ctx)
if err != nil {
    return err
}
```

## Factory Pattern

Use the factory pattern for managing multiple database types:

```go
import "github.com/arash-mosavi/go-base-repository/pkg/factory"

// Create configuration
config := &config.DatabaseConfig{
    Type: "mongodb", // or "postgres"
    Host: "localhost",
    Port: 27017,
    Name: "example_db",
}

// Create repository manager
manager, err := factory.NewRepositoryManager(config)
if err != nil {
    return err
}

// Create repositories based on configuration
userRepo := manager.GetUserRepository()
```

## API Reference

### MongoBaseRepository[T]

**Methods:**
- `Create(ctx, entity) (T, error)`
- `FindByID(ctx, id) (T, error)`
- `FindAll(ctx, params) ([]T, error)`
- `Update(ctx, id, entity) (T, error)`
- `Delete(ctx, id) error`
- `DeleteMultiple(ctx, ids) error`
- `FindWithFilter(ctx, filter) ([]T, error)`
- `Count(ctx) (int64, error)`
- `Exists(ctx, id) (bool, error)`
- `BeginTransaction(ctx) (MongoRepository[T], error)`
- `CommitTransaction(ctx) error`
- `RollbackTransaction(ctx) error`

### PostgresBaseRepository[T]

**Methods:**
- `Create(ctx, entity) (T, error)`
- `FindByID(ctx, id) (T, error)`
- `FindAll(ctx, params) ([]T, error)`
- `Update(ctx, id, entity) (T, error)`
- `Delete(ctx, id) error`
- `DeleteMultiple(ctx, ids) error`
- `FindWithFilter(ctx, filter) ([]T, error)`
- `Count(ctx) (int64, error)`
- `Exists(ctx, id) (bool, error)`
- `BeginTransaction(ctx) (PostgresRepository[T], error)`
- `CommitTransaction(ctx) error`
- `RollbackTransaction(ctx) error`

## Testing

Run the tests:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test -cover ./...
```

## Examples

See the `examples/` directory for comprehensive usage examples:

- `examples/usage_example.go` - Complete usage examples for both databases
- `pkg/repository/repository_test.go` - Unit tests with mock implementations

## Design Principles

1. **Composition over Inheritance**: Uses composition to build repositories
2. **Dependency Injection**: BaseRepository is injected into concrete repositories
3. **Type Safety**: Leverages Go 1.18+ generics for compile-time type checking
4. **Separation of Concerns**: Separate interfaces for MongoDB and PostgreSQL
5. **Clean Architecture**: Domain models are separate from persistence concerns
6. **Idiomatic Go**: Follows Go best practices and conventions

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For questions and support, please open an issue on GitHub.
