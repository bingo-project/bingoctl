// ABOUTME: Tests for the migrate package
// ABOUTME: Verifies migrator functionality including batch number calculation

package migrate

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	return db
}

func TestGetBatch_ReturnsMaxBatchPlusOne(t *testing.T) {
	// Setup: create a test database with migration records where
	// the record with highest ID does NOT have the highest batch
	db := setupTestDB(t)

	migrator := NewMigrator(db)

	// Insert records where id=2 has batch=1 but id=1 has batch=3
	// This simulates a scenario where migrations were rolled back and re-run
	db.Exec("INSERT INTO bingo_migration (id, migration, batch) VALUES (1, 'migration_1', 3)")
	db.Exec("INSERT INTO bingo_migration (id, migration, batch) VALUES (2, 'migration_2', 1)")

	// getBatch should return max(batch) + 1 = 3 + 1 = 4
	// NOT id DESC first batch + 1 = 1 + 1 = 2
	batch := migrator.getBatch()

	if batch != 4 {
		t.Errorf("getBatch() = %d, want 4 (max batch 3 + 1)", batch)
	}
}

func TestGetBatch_ReturnsOneWhenNoMigrations(t *testing.T) {
	db := setupTestDB(t)
	migrator := NewMigrator(db)

	batch := migrator.getBatch()

	if batch != 1 {
		t.Errorf("getBatch() = %d, want 1 for empty table", batch)
	}
}

func TestRollback_RollsBackMaxBatch(t *testing.T) {
	// Setup: create records where id=1 has batch=3 but id=2 has batch=1
	// Rollback should roll back batch=3 (the max), not batch=1 (the latest id)
	db := setupTestDB(t)
	migrator := NewMigrator(db)

	db.Exec("INSERT INTO bingo_migration (id, migration, batch) VALUES (1, 'migration_1', 3)")
	db.Exec("INSERT INTO bingo_migration (id, migration, batch) VALUES (2, 'migration_2', 1)")

	migrator.Rollback()

	// After rollback, batch=3 record should be deleted, batch=1 should remain
	var remaining []Migration
	db.Find(&remaining)

	if len(remaining) != 1 {
		t.Fatalf("expected 1 remaining migration, got %d", len(remaining))
	}

	if remaining[0].Batch != 1 {
		t.Errorf("expected remaining migration to have batch=1, got batch=%d", remaining[0].Batch)
	}

	if remaining[0].Migration != "migration_2" {
		t.Errorf("expected remaining migration to be 'migration_2', got '%s'", remaining[0].Migration)
	}
}
