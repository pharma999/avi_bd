package database

import (
	"fmt"

	"avi_bd/internal/models"

	"gorm.io/gorm"
)

// RunMigrations runs all database migrations
func RunMigrations(db *gorm.DB) error {
	// Auto migrate all models
	return db.AutoMigrate(
		&models.User{},
		&models.Wallet{},
		&models.GameRound{},
		&models.Bet{},
		&models.Transaction{},
	)
}

// CreateIndexes creates necessary database indexes
func CreateIndexes(db *gorm.DB) error {
	// Indexes are defined in the models, but you can add additional ones here if needed

	// Composite index for bets: user_id and round_id
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_bets_user_round ON bets(user_id, round_id)").Error; err != nil {
		return fmt.Errorf("failed to create bet index: %w", err)
	}

	// Index for active bets
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_bets_result_round ON bets(result, round_id)").Error; err != nil {
		return fmt.Errorf("failed to create result index: %w", err)
	}

	// Index for transaction lookup
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_transactions_user_type ON transactions(user_id, type)").Error; err != nil {
		return fmt.Errorf("failed to create transaction index: %w", err)
	}

	return nil
}

// DropAllTables drops all tables (use with caution!)
func DropAllTables(db *gorm.DB) error {
	return db.Migrator().DropTable(
		&models.Transaction{},
		&models.Bet{},
		&models.GameRound{},
		&models.Wallet{},
		&models.User{},
	)
}
