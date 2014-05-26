package metadata

import (
	"database/sql"
)

type pgSchemaEnsurer struct {
	db *sql.DB
}

var (
	migrations = map[string][]string{
		"20140525125633": {`
	  CREATE TABLE IF NOT EXISTS artifacts_metadata (
		id serial PRIMARY KEY,
		owner character varying(128) NOT NULL,
		repo character varying(128) NOT NULL,
		build_id character varying(32) NOT NULL,
		build_number character varying(32) NOT NULL,
		job_id character varying(32) NOT NULL,
		job_number character varying(32) NOT NULL,
		path character varying(32) NOT NULL
	  );
	  `,
		},
	}

	statements = map[string]string{
		"insert_metadata": `
		INSERT INTO artifacts_metadata (
			owner,
			repo,
			build_id,
			build_number,
			job_id,
			job_number,
			path
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`,
	}
)

func newPGSchemaEnsurer(db *sql.DB) *pgSchemaEnsurer {
	return &pgSchemaEnsurer{db}
}

func (pg *pgSchemaEnsurer) EnsureSchema() error {
	if err := pg.ensureMigrationsTable(); err != nil {
		return err
	}
	return pg.runMigrations()
}

func (pg *pgSchemaEnsurer) ensureMigrationsTable() error {
	_, err := pg.db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (version character varying(255) NOT NULL);`)
	return err
}

func (pg *pgSchemaEnsurer) runMigrations() error {
	for schemaVersion, sqls := range migrations {
		if pg.containsMigration(schemaVersion) {
			continue
		}

		if err := pg.migrateTo(schemaVersion, sqls); err != nil {
			return err
		}
	}
	return nil
}

func (pg *pgSchemaEnsurer) containsMigration(schemaVersion string) bool {
	var count int
	if err := pg.db.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE version = $1", schemaVersion).Scan(&count); err != nil {
		return false
	}

	return count == 1
}

func (pg *pgSchemaEnsurer) migrateTo(schemaVersion string, sqls []string) error {
	var (
		tx  *sql.Tx
		err error
	)

	if tx, err = pg.db.Begin(); err != nil {
		return err
	}

	for _, sql := range sqls {
		if _, err = tx.Exec(sql); err != nil {
			tx.Rollback()
			return err
		}
	}
	if _, err = tx.Exec("INSERT INTO schema_migrations VALUES ($1)", schemaVersion); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
