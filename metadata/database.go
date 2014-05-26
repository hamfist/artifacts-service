package metadata

import (
	"database/sql"

	"github.com/lib/pq"
)

// Database holds on to metadata. Wow!
type Database struct {
	url  string
	conn *sql.DB
	st   map[string]*sql.Stmt
}

// NewDatabase creates a (postgres) *Database from a database URL string
func NewDatabase(url string) (*Database, error) {
	db := &Database{
		url: url,
		st: map[string]*sql.Stmt{
			"insert_metadata": nil,
		},
	}

	parsedURL, err := pq.ParseURL(url)
	if err != nil {
		return nil, err
	}

	conn, err := sql.Open("postgres", parsedURL)
	if err != nil {
		return nil, err
	}

	db.conn = conn

	err = db.init()
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Save saves some metadata
func (db *Database) Save(m *Metadata) error {
	return nil
}

// Lookup looks up some metadata
func (db *Database) Lookup(slug, jobID, path string) (*Metadata, error) {
	return nil, errNoMetadata
}

func (db *Database) init() error {
	err := newPGSchemaEnsurer(db.conn).EnsureSchema()
	if err != nil {
		return err
	}

	if err = db.establishConnection(); err != nil {
		return err
	}

	if err = db.prepareStatements(); err != nil {
		return err
	}

	return nil
}

func (db *Database) establishConnection() error {
	conn, err := sql.Open("postgres", db.url)
	if err != nil {
		return err
	}
	db.conn = conn

	return nil
}

func (db *Database) prepareStatements() error {
	for key, statementSQL := range statements {
		stmt, err := db.conn.Prepare(statementSQL)

		if err != nil {
			return err
		}

		db.st[key] = stmt
	}

	return nil
}
