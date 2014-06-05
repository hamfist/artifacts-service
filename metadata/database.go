package metadata

import (
	"database/sql"

	"github.com/Sirupsen/logrus"
	"github.com/lib/pq"
)

var (
	defaultMigrations = map[string][]string{
		"20140525125633": {`
	  CREATE TABLE IF NOT EXISTS artifacts_metadata (
		id serial PRIMARY KEY,
		job_id character varying(32) NOT NULL,
		size bigint NOT NULL,
		path character varying(1024) NOT NULL,
		content_type character varying(255) NOT NULL
	  );
	  `,
		},
	}

	defaultStatements = map[string]string{
		"insert_metadata": `
		INSERT INTO artifacts_metadata (
			job_id,
			size,
			path,
			content_type
		) VALUES ($1, $2, $3, $4)
	`,
	}
)

// Database holds on to metadata. Wow!
type Database struct {
	url  string
	conn *sql.DB
	st   map[string]*sql.Stmt

	migrations map[string][]string
	statements map[string]string
}

// NewDatabase creates a (postgres) *Database from a database URL string
func NewDatabase(url string) (*Database, error) {
	db := &Database{
		url: url,
		st: map[string]*sql.Stmt{
			"insert_metadata": nil,
		},
		migrations: defaultMigrations,
		statements: defaultStatements,
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

// Migrate runs any missing migrations for make great schema
func (db *Database) Migrate(log *logrus.Logger) error {
	return newPGSchemaEnsurer(db.conn, db.migrations, log).EnsureSchema()
}

// Init runs initialization stuff so that stuff works
func (db *Database) Init() error {
	if err := db.establishConnection(); err != nil {
		return err
	}

	if err := db.prepareStatements(); err != nil {
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
	for key, statementSQL := range db.statements {
		stmt, err := db.conn.Prepare(statementSQL)

		if err != nil {
			return err
		}

		db.st[key] = stmt
	}

	return nil
}
