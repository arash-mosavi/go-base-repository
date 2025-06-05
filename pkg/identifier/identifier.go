package identifier

import (
	"github.com/arash-mosavi/go-base-repository/pkg/types"
	mongoIdentifier "github.com/arash-mosavi/mongo-unit-of-work-system/pkg/identifier"
	postgresIdentifier "github.com/arash-mosavi/postgrs-unit-of-work-system/pkg/identifier"
)

// UnifiedIdentifier provides a unified interface for both MongoDB and PostgreSQL identifiers
type UnifiedIdentifier struct {
	mongoID    mongoIdentifier.IIdentifier
	postgresID postgresIdentifier.IIdentifier
}

// NewMongoIdentifier creates a new MongoDB identifier
func NewMongoIdentifier() *UnifiedIdentifier {
	return &UnifiedIdentifier{
		mongoID: mongoIdentifier.New(),
	}
}

// NewPostgresIdentifier creates a new PostgreSQL identifier
func NewPostgresIdentifier() *UnifiedIdentifier {
	return &UnifiedIdentifier{
		postgresID: postgresIdentifier.NewIdentifier(),
	}
}

// Equal adds an equality condition
func (u *UnifiedIdentifier) Equal(field string, value interface{}) types.Identifier {
	if u.mongoID != nil {
		u.mongoID = u.mongoID.Equal(field, value)
	}
	if u.postgresID != nil {
		u.postgresID = u.postgresID.Equal(field, value)
	}
	return u
}

// GreaterThan adds a greater than condition
func (u *UnifiedIdentifier) GreaterThan(field string, value interface{}) types.Identifier {
	if u.mongoID != nil {
		u.mongoID = u.mongoID.GreaterThan(field, value)
	}
	if u.postgresID != nil {
		u.postgresID = u.postgresID.GreaterThan(field, value)
	}
	return u
}

// LessThan adds a less than condition
func (u *UnifiedIdentifier) LessThan(field string, value interface{}) types.Identifier {
	if u.mongoID != nil {
		u.mongoID = u.mongoID.LessThan(field, value)
	}
	if u.postgresID != nil {
		u.postgresID = u.postgresID.LessThan(field, value)
	}
	return u
}

// Between adds a between condition
func (u *UnifiedIdentifier) Between(field string, min, max interface{}) types.Identifier {
	if u.mongoID != nil {
		u.mongoID = u.mongoID.Between(field, min, max)
	}
	if u.postgresID != nil {
		u.postgresID = u.postgresID.Between(field, min, max)
	}
	return u
}

// Like adds a like condition
func (u *UnifiedIdentifier) Like(field string, pattern string) types.Identifier {
	if u.mongoID != nil {
		u.mongoID = u.mongoID.Like(field, pattern)
	}
	if u.postgresID != nil {
		u.postgresID = u.postgresID.Like(field, pattern)
	}
	return u
}

// In adds an in condition
func (u *UnifiedIdentifier) In(field string, values []interface{}) types.Identifier {
	if u.mongoID != nil {
		u.mongoID = u.mongoID.In(field, values)
	}
	if u.postgresID != nil {
		u.postgresID = u.postgresID.In(field, values)
	}
	return u
}

// ToBSON returns the MongoDB BSON representation
func (u *UnifiedIdentifier) ToBSON() map[string]interface{} {
	if u.mongoID != nil {
		return u.mongoID.ToBSON()
	}
	return make(map[string]interface{})
}

// ToMap returns the PostgreSQL map representation
func (u *UnifiedIdentifier) ToMap() map[string]interface{} {
	if u.postgresID != nil {
		return u.postgresID.ToMap()
	}
	return make(map[string]interface{})
}

// Has checks if a field exists in the identifier
func (u *UnifiedIdentifier) Has(field string) bool {
	if u.mongoID != nil {
		return u.mongoID.Has(field)
	}
	if u.postgresID != nil {
		// Note: PostgreSQL identifier doesn't have Has method, so we check the map
		_, exists := u.postgresID.ToMap()[field]
		return exists
	}
	return false
}

// GetMongoIdentifier returns the underlying MongoDB identifier
func (u *UnifiedIdentifier) GetMongoIdentifier() mongoIdentifier.IIdentifier {
	return u.mongoID
}

// GetPostgresIdentifier returns the underlying PostgreSQL identifier
func (u *UnifiedIdentifier) GetPostgresIdentifier() postgresIdentifier.IIdentifier {
	return u.postgresID
}
