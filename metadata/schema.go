package metadata

import (
	"database/sql"
	"github.com/Sirupsen/logrus"
)

type pgSchemaEnsurer struct {
	db  *sql.DB
	log *logrus.Logger

	migrations map[string][]string
}

func newPGSchemaEnsurer(db *sql.DB, migrations map[string][]string, log *logrus.Logger) *pgSchemaEnsurer {
	return &pgSchemaEnsurer{
		db:  db,
		log: log,

		migrations: migrations,
	}
}

func (pg *pgSchemaEnsurer) EnsureSchema() error {
	if err := pg.ensureMigrationsTable(); err != nil {
		return err
	}
	return pg.runMigrations()
}

func (pg *pgSchemaEnsurer) ensureMigrationsTable() error {
	pg.log.Info("ensuring schema_migrations table exists")
	_, err := pg.db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (version character varying(255) NOT NULL);`)
	return err
}

func (pg *pgSchemaEnsurer) runMigrations() error {
	for schemaVersion, sqls := range pg.migrations {
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
