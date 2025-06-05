package factory

import (
	"github.com/arash-mosavi/go-base-repository/pkg/interfaces"
	"github.com/arash-mosavi/go-base-repository/pkg/mongo"
	"github.com/arash-mosavi/go-base-repository/pkg/postgres"
	"github.com/arash-mosavi/go-base-repository/pkg/types"
	mongoFactory "github.com/arash-mosavi/mongo-unit-of-work-system/pkg/mongodb"
	postgresFactory "github.com/arash-mosavi/postgrs-unit-of-work-system/pkg/postgres"
)

// MongoConfig wraps MongoDB configuration
type MongoConfig struct {
	*mongoFactory.Config
}

// PostgresConfig wraps PostgreSQL configuration
type PostgresConfig struct {
	*postgresFactory.Config
}

// NewMongoConfig creates a new MongoDB configuration
func NewMongoConfig() *MongoConfig {
	return &MongoConfig{
		Config: mongoFactory.NewConfig(),
	}
}

// NewPostgresConfig creates a new PostgreSQL configuration
func NewPostgresConfig() *PostgresConfig {
	return &PostgresConfig{
		Config: postgresFactory.NewConfig(),
	}
}

// NewMongoBaseRepository creates a new MongoDB base repository
func NewMongoBaseRepository[T types.MongoEntity](config *MongoConfig) (interfaces.MongoBaseRepository[T], error) {
	factory, err := mongoFactory.NewFactory[T](config.Config)
	if err != nil {
		return nil, err
	}
	return mongo.NewBaseRepository[T](factory), nil
}

// NewPostgresBaseRepository creates a new PostgreSQL base repository
func NewPostgresBaseRepository[T types.PostgresEntity](config *PostgresConfig) (interfaces.PostgresBaseRepository[T], error) {
	factory := postgresFactory.NewUnitOfWorkFactory[T](config.Config)
	return postgres.NewBaseRepository[T](factory), nil
}
