package database

import (
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	err := godotenv.Load("../../.env")
	if err != nil {
		t.Fatal("Error loading .env file")
	}

	os.Setenv("LOGGER_DISABLED", "true")

	db := _NewDatabaseConnection()
	dbInstance := db.Connect()

	assert.NotNil(t, dbInstance)

	defer dbInstance.Close()
}

func TestGetDatabase(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("CREATE TABLE IF NOT EXISTS tasks").WillReturnResult(sqlmock.NewResult(1, 1))

	sqlDB := db
	dbInstance = sqlDB

	result := GetDatabase()

	assert.NotNil(t, result)
	assert.Equal(t, sqlDB, result)

	dbInstance = nil

	result = GetDatabase()

	assert.NotNil(t, result)
	assert.NotEqual(t, sqlDB, result)
}

func TestClose(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mock.ExpectClose()

	sqlDB := db
	dbInstance = sqlDB

	Close()

	assert.Nil(t, dbInstance)
}
