package migrate

import (
	"fmt"
	"os"

	"github.com/bingo-project/component-base/cli/console"
	"github.com/mgutz/ansi"
	"gorm.io/gorm"
)

// Default table name for migration records
const DefaultTableName = "bingo_migration"

// migrationTableName stores the configured table name
var migrationTableName = DefaultTableName

func init() {
	// Allow overriding table name via environment variable
	if tableName := os.Getenv("BINGOCTL_MIGRATE_TABLE"); tableName != "" {
		migrationTableName = tableName
	}
}

// SetTableName sets the migration table name
func SetTableName(name string) {
	if name != "" {
		migrationTableName = name
	}
}

type Migrator struct {
	DB       *gorm.DB
	Migrator gorm.Migrator
}

type Migration struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement;"`
	Migration string `gorm:"type:varchar(255);not null;unique;"`
	Batch     int
}

// TableName returns the configured migration table name
func (Migration) TableName() string {
	return migrationTableName
}

// NewMigrator creates a new Migrator with default table name.
func NewMigrator(db *gorm.DB) *Migrator {
	return NewMigratorWithTable(db, DefaultTableName)
}

// NewMigratorWithTable creates a new Migrator with custom table name.
func NewMigratorWithTable(db *gorm.DB, tableName string) *Migrator {
	if tableName != "" {
		migrationTableName = tableName
	}

	migrator := &Migrator{
		DB:       db,
		Migrator: db.Migrator(),
	}

	migrator.createMigrationsTable()

	return migrator
}

func (migrator *Migrator) createMigrationsTable() {
	migration := Migration{}

	if !migrator.Migrator.HasTable(&migration) {
		_ = migrator.Migrator.CreateTable(&migration)
	}
}

func (migrator *Migrator) Up() {
	// Get batch
	batch := migrator.getBatch()

	var migrations []Migration
	migrator.DB.Find(&migrations)

	ran := false
	for _, migrationFile := range migrationFiles {
		if isNotMigrated(migrations, migrationFile) {
			migrator.runUpMigration(migrationFile, batch)
			ran = true
		}
	}

	if !ran {
		console.Info("Nothing to migrate.")
	}
}

func (migrator *Migrator) Rollback() {
	lastMigration := Migration{}
	migrator.DB.Order("id DESC").First(&lastMigration)

	var migrations []Migration
	migrator.DB.Where("batch = ?", lastMigration.Batch).Order("id DESC").Find(&migrations)

	if !migrator.rollbackMigrations(migrations) {
		console.Info("Nothing to rollback.")
	}
}

func (migrator *Migrator) rollbackMigrations(migrations []Migration) bool {
	ran := false

	for _, _migration := range migrations {
		fmt.Printf("%s %s\n", ansi.Color("Rolling back:", "yellow"), _migration.Migration)

		migrationFile := GetMigrationFile(_migration.Migration)
		if migrationFile.Down != nil {
			migrationFile.Down(migrator.DB.Migrator())
		}

		ran = true

		migrator.DB.Delete(&_migration)

		fmt.Printf("%s  %s\n", ansi.Color("Rolled back:", "green"), migrationFile.FileName)
	}

	return ran
}

func (migrator *Migrator) getBatch() int {
	batch := 1

	lastMigration := Migration{}
	migrator.DB.Order("id DESC").First(&lastMigration)

	if lastMigration.ID > 0 {
		batch = lastMigration.Batch + 1
	}

	return batch
}

func (migrator *Migrator) runUpMigration(migrationFile MigrationFile, batch int) {
	if migrationFile.Up != nil {
		fmt.Printf("%s %s\n", ansi.Color("Migrating:", "yellow"), migrationFile.FileName)

		migrationFile.Up(migrator.DB.Migrator())

		fmt.Printf("%s  %s\n", ansi.Color("Migrated:", "green"), migrationFile.FileName)
	}

	err := migrator.DB.Create(&Migration{Migration: migrationFile.FileName, Batch: batch}).Error
	console.ExitIf(err)
}

func (migrator *Migrator) Reset() {
	var migrations []Migration

	migrator.DB.Order("id DESC").Find(&migrations)

	if !migrator.rollbackMigrations(migrations) {
		console.Info("Nothing to rollback.")
	}
}

func (migrator *Migrator) Refresh() {
	migrator.Reset()

	migrator.Up()
}

func (migrator *Migrator) Fresh() {
	// Delete all tables
	err := migrator.DeleteAllTables()
	console.ExitIf(err)
	console.Info("Dropped all tables successfully.")

	// Migrate
	migrator.createMigrationsTable()
	console.Info("Migration table created successfully.")

	migrator.Up()
}

func isNotMigrated(migrations []Migration, migrationFile MigrationFile) bool {
	for _, migration := range migrations {
		if migration.Migration == migrationFile.FileName {
			return false
		}
	}

	return true
}

func (migrator *Migrator) DeleteAllTables() error {
	tables, err := migrator.DB.Migrator().GetTables()
	if err != nil {
		return err
	}

	for _, table := range tables {
		err := migrator.DB.Migrator().DropTable(table)
		if err != nil {
			continue
		}
	}

	return nil
}
