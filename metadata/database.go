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
				id,
				job_id,
				size,
				path,
				content_type
			FROM artifacts_metadata
			WHERE job_id = $1
				AND path = $2`,
		},
		SelectAllMetadataForJobID: &dbStatement{
			SQL: `
			SELECT
				id,
				job_id,
				size,
				path,
				content_type
			FROM artifacts_metadata
			WHERE job_id = $1`,
		},
	}
)

type dbStatement struct {
	SQL string
	St  *sql.Stmt
}

type dbStatements struct {
	InsertMetadata            *dbStatement
	SelectMetadata            *dbStatement
	SelectAllMetadataForJobID *dbStatement
}

// Database holds on to metadata. Wow!
type Database struct {
	url  string
	conn *sql.DB
	st   *dbStatements
	log  *logrus.Logger

	migrations           map[string][]string
	unpreparedStatements *dbStatements
}

// NewDatabase creates a (postgres) *Database from a database URL string
func NewDatabase(url string, log *logrus.Logger) (*Database, error) {
	db := &Database{
		url:                  url,
		st:                   &dbStatements{},
		log:                  log,
		migrations:           defaultMigrations,
		unpreparedStatements: defaultStatements,
	}

	if db.log == nil {
		db.log = logrus.New()
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
	statement := db.st.InsertMetadata.St
	if statement == nil {
		return errNoStatement
	}

	res, err := statement.Exec(m.JobID, m.Size, m.Path, m.ContentType)
	logFields := logrus.Fields{
		"job_id":       m.JobID,
		"size":         m.Size,
		"path":         m.Path,
		"content_type": m.ContentType,
		"result":       res,
	}
	if err != nil {
		logFields["err"] = err
		db.log.WithFields(logFields).Error("error saving metadata")
	} else {
		db.log.WithFields(logFields).Info("saving metadata")
	}

	return err
}

// Lookup looks up some metadata
func (db *Database) Lookup(jobID, path string) (*Metadata, error) {
	statement := db.st.SelectMetadata.St
	if statement == nil {
		return nil, errNoStatement
	}

	md := &Metadata{}

	err := statement.QueryRow(jobID, path).Scan(&md.ID, &md.JobID, &md.Size, &md.Path, &md.ContentType)
	if err != nil {
		return nil, errNoMetadata
	}

	return md, nil
}

// LookupAll returns all metadata associated with a given job ID
func (db *Database) LookupAll(jobID string) ([]*Metadata, error) {
	statement := db.st.SelectAllMetadataForJobID.St
	if statement == nil {
		return nil, errNoStatement
	}

	mds := []*Metadata{}
	rows, err := statement.Query(jobID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		md := &Metadata{}
		err = rows.Scan(&md.ID, &md.JobID, &md.Size, &md.Path, &md.ContentType)
		if err != nil {
			return nil, err
		}
		mds = append(mds, md)
	}

	db.log.WithFields(logrus.Fields{
		"job_id":       jobID,
		"all_metadata": mds,
	}).Info("returning metadata for job")

	return mds, nil
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
		db.log.WithFields(logrus.Fields{
			"err": err,
		}).Error("failed to prepare insert metadata statement")
		return err
	}

	pr.InsertMetadata.St = insertStmt

	selectStmt, err := db.conn.Prepare(pr.SelectMetadata.SQL)
	if err != nil {
		db.log.WithFields(logrus.Fields{
			"err": err,
		}).Error("failed to prepare select metadata statement")
		return err
	}

	pr.SelectMetadata.St = selectStmt

	selectAllStmt, err := db.conn.Prepare(pr.SelectAllMetadataForJobID.SQL)
	if err != nil {
		db.log.WithFields(logrus.Fields{
			"err": err,
		}).Error("failed to prepare select all metadata statement")
		return err
	}

	pr.SelectAllMetadataForJobID.St = selectAllStmt

	db.st = pr
	return nil
}
