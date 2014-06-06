package metadata

import (
	"database/sql"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/lib/pq"
)

var (
	errNoStatement    = fmt.Errorf("no prepared statement")
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

	defaultStatements = &dbStatements{
		InsertMetadata: &dbStatement{
			SQL: `
			INSERT INTO artifacts_metadata (
				job_id,
				size,
				path,
				content_type
			) VALUES ($1, $2, $3, $4)`,
		},
		SelectMetadata: &dbStatement{
			SQL: `
			SELECT
				job_id,
				size,
				path,
				content_type
			FROM artifacts_metadata
			WHERE job_id = $1
				AND path = $2`,
		},
	}
)

type dbStatement struct {
	SQL string
	St  *sql.Stmt
}

type dbStatements struct {
	InsertMetadata *dbStatement
	SelectMetadata *dbStatement
}

// Database holds on to metadata. Wow!
type Database struct {
	url  string
	conn *sql.DB
	st   *dbStatements

	migrations           map[string][]string
	unpreparedStatements *dbStatements
}

// NewDatabase creates a (postgres) *Database from a database URL string
func NewDatabase(url string) (*Database, error) {
	db := &Database{
		url:                  url,
		st:                   &dbStatements{},
		migrations:           defaultMigrations,
		unpreparedStatements: defaultStatements,
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
func (db *Database) Save(m *Metadata) (sql.Result, error) {
	statement := db.st.InsertMetadata.St
	if statement == nil {
		return nil, errNoStatement
	}

	return statement.Exec(m.JobID, m.Size, m.Path, m.ContentType)
}

// Lookup looks up some metadata
func (db *Database) Lookup(jobID, path string) (*Metadata, error) {
	statement := db.st.SelectMetadata.St
	if statement == nil {
		return nil, errNoStatement
	}

	md := &Metadata{}

	err := statement.QueryRow(jobID, path).Scan(&md.JobID, &md.Size, &md.Path, &md.ContentType)
	if err != nil {
		return nil, errNoMetadata
	}

	return md, nil
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
	pr := db.unpreparedStatements

	insertStmt, err := db.conn.Prepare(pr.InsertMetadata.SQL)
	if err != nil {
		return err
	}

	pr.InsertMetadata.St = insertStmt

	selectStmt, err := db.conn.Prepare(pr.SelectMetadata.SQL)
	if err != nil {
		return err
	}

	pr.SelectMetadata.St = selectStmt

	db.st = pr
	return nil
}
