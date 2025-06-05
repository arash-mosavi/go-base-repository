package identifier_test

import (
	"testing"

	"github.com/arash-mosavi/go-base-repository/pkg/identifier"
	"github.com/arash-mosavi/go-base-repository/pkg/types"
)

func TestUnifiedIdentifier_MongoDB(t *testing.T) {
	mongoID := identifier.NewMongoIdentifier()

	// Test that it implements the Identifier interface
	var id types.Identifier = mongoID
	if id == nil {
		t.Error("MongoDB identifier should implement Identifier interface")
	}

	// Test builder methods
	mongoID = mongoID.Equal("name", "test").(*identifier.UnifiedIdentifier)
	mongoID = mongoID.GreaterThan("age", 18).(*identifier.UnifiedIdentifier)
	mongoID = mongoID.LessThan("age", 65).(*identifier.UnifiedIdentifier)
	mongoID = mongoID.Between("score", 50, 100).(*identifier.UnifiedIdentifier)
	mongoID = mongoID.Like("email", "%@example.com").(*identifier.UnifiedIdentifier)
	mongoID = mongoID.In("status", []interface{}{"active", "pending"}).(*identifier.UnifiedIdentifier)

	// Test MongoDB-specific methods
	bson := mongoID.ToBSON()
	if bson == nil {
		t.Error("ToBSON should return a valid map")
	}

	// Test Has method
	has := mongoID.Has("name")
	if !has {
		t.Error("Should have the 'name' field")
	}

	// Test GetMongoIdentifier
	mongoIdentifier := mongoID.GetMongoIdentifier()
	if mongoIdentifier == nil {
		t.Error("GetMongoIdentifier should return a valid identifier")
	}
}

func TestUnifiedIdentifier_PostgreSQL(t *testing.T) {
	postgresID := identifier.NewPostgresIdentifier()

	// Test that it implements the Identifier interface
	var id types.Identifier = postgresID
	if id == nil {
		t.Error("PostgreSQL identifier should implement Identifier interface")
	}

	// Test builder methods
	postgresID = postgresID.Equal("name", "test").(*identifier.UnifiedIdentifier)
	postgresID = postgresID.GreaterThan("age", 18).(*identifier.UnifiedIdentifier)
	postgresID = postgresID.LessThan("age", 65).(*identifier.UnifiedIdentifier)
	postgresID = postgresID.Between("score", 50, 100).(*identifier.UnifiedIdentifier)
	postgresID = postgresID.Like("email", "%@example.com").(*identifier.UnifiedIdentifier)
	postgresID = postgresID.In("status", []interface{}{"active", "pending"}).(*identifier.UnifiedIdentifier)

	// Test PostgreSQL-specific methods
	mapping := postgresID.ToMap()
	if mapping == nil {
		t.Error("ToMap should return a valid map")
	}

	// Test Has method
	has := postgresID.Has("name")
	if !has {
		t.Error("Should have the 'name' field")
	}

	// Test GetPostgresIdentifier
	postgresIdentifier := postgresID.GetPostgresIdentifier()
	if postgresIdentifier == nil {
		t.Error("GetPostgresIdentifier should return a valid identifier")
	}
}

func TestUnifiedIdentifier_Fluent(t *testing.T) {
	mongoID := identifier.NewMongoIdentifier()

	// Test fluent interface
	result := mongoID.
		Equal("name", "John").
		GreaterThan("age", 18).
		Like("email", "%@example.com")

	if result == nil {
		t.Error("Fluent interface should return a valid identifier")
	}

	// Ensure it's still the same type
	unifiedResult, ok := result.(*identifier.UnifiedIdentifier)
	if !ok {
		t.Error("Result should be of type UnifiedIdentifier")
	}

	if unifiedResult == nil {
		t.Error("Unified result should not be nil")
	}
}

func TestUnifiedIdentifier_TypeSafety(t *testing.T) {
	// Test that both identifiers can be used as the same interface
	var mongoID types.Identifier = identifier.NewMongoIdentifier()
	var postgresID types.Identifier = identifier.NewPostgresIdentifier()

	// Both should implement the same interface
	if mongoID == nil || postgresID == nil {
		t.Error("Both identifiers should implement the Identifier interface")
	}

	// Test that methods return the interface type
	mongoResult := mongoID.Equal("test", "value")
	postgresResult := postgresID.Equal("test", "value")

	if mongoResult == nil || postgresResult == nil {
		t.Error("Both identifiers should return valid results")
	}

	// Test method chaining with interface
	chainedMongo := mongoID.Equal("a", 1).GreaterThan("b", 2).LessThan("c", 3)
	chainedPostgres := postgresID.Equal("a", 1).GreaterThan("b", 2).LessThan("c", 3)

	if chainedMongo == nil || chainedPostgres == nil {
		t.Error("Method chaining should work for both identifiers")
	}
}
