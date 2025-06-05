package identifier

import (
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BenchmarkUnifiedIdentifier_MongoDB benchmarks MongoDB identifier operations
func BenchmarkUnifiedIdentifier_MongoDB(b *testing.B) {
	mongoID := primitive.NewObjectID()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		id := NewMongoIdentifier().
			Equal("field", mongoID)
		_ = id.ToBSON()
	}
}

// BenchmarkUnifiedIdentifier_PostgreSQL benchmarks PostgreSQL identifier operations
func BenchmarkUnifiedIdentifier_PostgreSQL(b *testing.B) {
	pgID := int64(12345)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		id := NewPostgresIdentifier().
			Equal("field", pgID)
		_ = id.ToMap()
	}
}

// BenchmarkUnifiedIdentifier_Creation benchmarks identifier creation
func BenchmarkUnifiedIdentifier_Creation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewMongoIdentifier()
		_ = NewPostgresIdentifier()
	}
}

// BenchmarkUnifiedIdentifier_TypeConversion benchmarks type conversion
func BenchmarkUnifiedIdentifier_TypeConversion(b *testing.B) {
	mongoID := primitive.NewObjectID()
	pgID := int64(12345)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Test MongoDB conversion
		id1 := NewMongoIdentifier().Equal("field", mongoID)
		_ = id1.ToBSON()

		// Test PostgreSQL conversion
		id2 := NewPostgresIdentifier().Equal("field", pgID)
		_ = id2.ToMap()
	}
}
