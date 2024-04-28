package create_table

import "database/sql"

// TODO: implement table versioning, and move this code to
//
//	`schema/create_table/vX`, where X = 1, 2, 3, etc.

var CreateTables = map[string]func(*sql.DB) error{
	"variables": CreateTableVariables,
	"accounts":  CreatTableAccounts,
}

// -- constants --

// CreateVariablesTableQuery holds the query for creating `variables` table
const CreateVariablesTableQuery = `CREATE TABLE IF NOT EXISTS variables (
	name TEXT,
	version TEXT,
	value TEXT,
	PRIMARY KEY (name, version)
);`

// CreateAccountsTableQuery holds the query for creating `accounts` table
const CreateAccountsTableQuery = `CREATE TABLE IF NOT EXISTS accounts (
	email VARCHAR(320) PRIMARY KEY,
	label VARCHAR(320) NOT NULL,
	smtp_server VARCHAR(320) NOT NULL,
	smtp_port INT NOT NULL,
	smtp_username VARCHAR(320) NOT NULL,
	smtp_password TEXT NOT NULL,
	imap_server VARCHAR(320) NOT NULL,
	imap_port INT NOT NULL,
	imap_username VARCHAR(320) NOT NULL,
	imap_password TEXT NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP
)`

// -- functions --

// CreateTableVariables creates the `variables` table
func CreateTableVariables(db *sql.DB) error {
	// TODO: add table versioning
	_, err := db.Exec(CreateVariablesTableQuery)
	return err
}

// CreatTableAccounts creates the `accounts` table
func CreatTableAccounts(db *sql.DB) error {
	// TODO: add table versioning
	_, err := db.Exec(CreateAccountsTableQuery)
	return err
}
