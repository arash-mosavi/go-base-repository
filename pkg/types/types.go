package types

import (
	mongoUOW "github.com/arash-mosavi/mongo-unit-of-work-system/pkg/persistence"
	postgresDomain "github.com/arash-mosavi/postgrs-unit-of-work-system/pkg/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MongoID represents MongoDB ObjectID
type MongoID = primitive.ObjectID

// PostgresID represents PostgreSQL integer ID
type PostgresID = int

// MongoEntity defines the constraint for MongoDB entities
type MongoEntity interface {
	mongoUOW.ModelConstraint
}

// PostgresEntity defines the constraint for PostgreSQL entities
type PostgresEntity interface {
	postgresDomain.BaseModel
}

// Identifier represents query identifiers for both database types
type Identifier interface {
	// MongoDB methods
	ToBSON() map[string]interface{}
	Has(field string) bool

	// PostgreSQL methods
	ToMap() map[string]interface{}

	// Common builder methods
	Equal(field string, value interface{}) Identifier
	GreaterThan(field string, value interface{}) Identifier
	LessThan(field string, value interface{}) Identifier
	Between(field string, min, max interface{}) Identifier
	Like(field string, pattern string) Identifier
	In(field string, values []interface{}) Identifier
}

// SortDirection represents sort direction
type SortDirection string

const (
	SortAsc  SortDirection = "ASC"
	SortDesc SortDirection = "DESC"
)

// SortMap defines sorting parameters
type SortMap map[string]SortDirection

// QueryParams defines query parameters for pagination and filtering
type QueryParams[T any] struct {
	Filter  T
	Limit   int
	Offset  int
	Sort    SortMap
	Include []string
}
