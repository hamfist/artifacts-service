package store

import (
	"database/sql"

	// register the pq driver
	_ "github.com/lib/pq"

	"github.com/meatballhat/artifacts-service/artifact"
)

// PostgreSQLStore is a storer specific to PostgreSQL.  Big surprise!
type PostgreSQLStore struct {
	conn *sql.DB
}

// NewPostgreSQLStore creates a *PostgreSQLStore from a database URL string
func NewPostgreSQLStore(url string) (*PostgreSQLStore, error) {
	conn, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	return &PostgreSQLStore{
		conn: conn,
	}, nil
}

// Store saves the stuff
func (pg *PostgreSQLStore) Store(a *artifact.Artifact) error {
	return nil
}
