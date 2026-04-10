// package database

// import (
// 	"fmt"

// 	"gorm.io/driver/postgres"
// 	"gorm.io/gorm"
// )

// // InitDatabase initializes the database connection
// func InitDatabase(connectionString string) (*gorm.DB, error) {
// 	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to connect to database: %w", err)
// 	}

// 	return db, nil
// }

// // GetDatabase returns a database instance (for dependency injection)
// func GetDatabase(connectionString string) (*gorm.DB, error) {
// 	return InitDatabase(connectionString)
// }

// // CloseDatabase closes the database connection
// func CloseDatabase(db *gorm.DB) error {
// 	sqlDB, err := db.DB()
// 	if err != nil {
// 		return err
// 	}
// 	return sqlDB.Close()
// }

package database

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDatabase(connectionString string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("gorm open failed: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	return db, nil
}
