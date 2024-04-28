package migrations

import (
	// stdlib imports
	"database/sql"

	// third-party imports
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog"

	// internal imports
	// TODO: decide on how to do SQLite encryption later
	// database_base "github.com/gouthamkrishnakv/zerocounter/database/base"
	migrations_schema "github.com/gouthamkrishnakv/zerocounter/database/schema"
	"github.com/gouthamkrishnakv/zerocounter/logging"
)

// -- variables --

var logger *zerolog.Logger = nil

// -- function --

// Run runs all the migrations
func Run(dbPath string) error {
	migrationsLogger := logging.L().With().Str("module", "migrations").Logger()
	logger = &migrationsLogger

	// open database
	db, dbErr := sql.Open("sqlite3", dbPath)
	if dbErr != nil {
		return dbErr
	}
	defer db.Close()
	logger.Debug().Str("db_path", dbPath).Msg("database opened for migration")

	// Schema Migrations
	if tableName, migrationErr := migrations_schema.Run(db); migrationErr != nil {
		logger.Error().Err(migrationErr).Str("table_name", tableName).Msg("migration failed")
		return migrationErr
	}
	logger.Debug().Str("migration_kind", "schema").Msg("schema migration successful")

	return nil
}
