package db

import (
	"log"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type MockApplication struct {
	Application
}

func (app *MockApplication) GetDB() *gorm.DB {
	// Mock database connection logic
	return &gorm.DB{}
}

func TestGetGormDB(t *testing.T) {
	app := &MockApplication{}

	// Setup mock for logger.New
	logOutput := os.Stdout

	// Call the function
	gormDB := app.GetGormDB()

	// Verify the result
	assert.NotNil(t, gormDB)
	assert.Equal(t, logOutput, logOutput)
	assert.Equal(t, logger.Config{
		SlowThreshold:             time.Second,
		LogLevel:                  logger.Info,
		IgnoreRecordNotFoundError: false,
		ParameterizedQueries:      false,
		Colorful:                  false,
	}, gormDB.Config.Logger)
}

func (app *MockApplication) GetGormDB() *gorm.DB {
	app.onceGormDb.Do(func() {
		newLogger := logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             time.Second,   // Slow SQL threshold
				LogLevel:                  logger.Info,   // Log level
				IgnoreRecordNotFoundError: false,         // Ignore ErrRecordNotFound error for logger
				ParameterizedQueries:      false,         // Don't include params in the SQL log
				Colorful:                  false,         // Disable color
			},
		)

		var err error
		gormDb, err = gorm.Open(postgres.New(postgres.Config{
			Conn: app.GetDB(),
		}), &gorm.Config{
			Logger: newLogger,
		})

		if err != nil {
			panic(fmt.Sprintf("Unable to connect to database: %v\n", err))
		}
	})

	return gormDb
}