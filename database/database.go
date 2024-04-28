package database

import (
	// stdlib imports
	"database/sql"

	// third-pary imports
	// zerolog logger
	"github.com/rs/zerolog"
	// import SQLite3 imports so "sql.Open" method can access the SQLite driver
	_ "github.com/mattn/go-sqlite3"

	// internal imports
	"github.com/gouthamkrishnakv/zerocounter/database/migrations"
	"github.com/gouthamkrishnakv/zerocounter/logging"
	utils_fs "github.com/gouthamkrishnakv/zerocounter/utils/filesystem"
)

// -- constants --

// DefaultDatabaseFile is the name of the default database file
const DefaultDatabaseFile = "zc.db"

// -- variables --

// databasePath holds the database path being "set"
var databasePath = ""

// logger for database module
var logger *zerolog.Logger = nil

// -- structs --

// Database struct defines and holds an SQL database connection
type Database struct {
	driver *sql.DB
}

var database *Database = nil

// -- functions --

// Initialize initializes the database for the service on launch
func Initialize() error {
	// setup logger
	newLogger := logging.L().
		// add timestamp, and stack printing, set "module" as database
		With().Str("module", "database").Logger()
	logger = &newLogger
	// search for DatabasePath in XDG directories
	filePath, searchErr := utils_fs.SearchDataFile(DefaultDatabaseFile)
	if searchErr != nil {
		var setupErr error
		logger.Warn().Str("relative_db_path", DefaultDatabaseFile).Err(searchErr).Msg("database not found, running database setup")
		filePath, setupErr = setup()
		if setupErr != nil {
			return setupErr
		}
		logger.Info().Str("db_path", filePath).Msg("database setup successful")
	}

	return loadDatabase(filePath)
}

func loadDatabase(filePath string) error {
	db, openErr := sql.Open("sqlite3", filePath)
	// if no error, set database
	if openErr != nil {
		return openErr
	}
	logger.Debug().Str("file_path", filePath).Msg("database file opened")
	database = &Database{
		db,
	}
	return nil
}

func setup() (string, error) {
	// generate file-path for data files, and Generate file-path for data files, and then send back the generated
	// standard file-path.
	dbPath, dbFileErr := utils_fs.DataFile(DefaultDatabaseFile)
	if dbFileErr != nil {
		return "", dbFileErr
	}
	logger.Debug().Str("db_path", dbPath).Msg("xdg-style data path generated and folders created")

	// do setup migrations, and then record any errors
	if setupMigrErr := migrations.Run(dbPath); setupMigrErr != nil {
		return "", setupMigrErr
	}
	logger.Debug().Str("db_path", dbPath).Msg("database migration in provided database path successful")
	return dbPath, nil
}
