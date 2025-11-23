package database

import (
	"moonbeam/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestInitWithSQLite(t *testing.T) {
	// For testing, we'll use SQLite in-memory
	// Test the migration part
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.User{})
	require.NoError(t, err)

	// Verify table exists
	var user models.User
	err = db.First(&user).Error
	assert.Error(t, err) // Should error because table is empty, not because table doesn't exist
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestUserModel(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.User{})
	require.NoError(t, err)

	// Create a user
	user := models.User{
		Name:  "Test User",
		Email: "test@example.com",
	}
	err = db.Create(&user).Error
	require.NoError(t, err)
	assert.NotZero(t, user.ID)
	assert.NotZero(t, user.CreatedAt)
	assert.NotZero(t, user.UpdatedAt)

	// Retrieve user
	var retrievedUser models.User
	err = db.First(&retrievedUser, user.ID).Error
	require.NoError(t, err)
	assert.Equal(t, user.Name, retrievedUser.Name)
	assert.Equal(t, user.Email, retrievedUser.Email)

	// Update user
	retrievedUser.Name = "Updated Name"
	err = db.Save(&retrievedUser).Error
	require.NoError(t, err)

	// Verify update
	var updatedUser models.User
	err = db.First(&updatedUser, user.ID).Error
	require.NoError(t, err)
	assert.Equal(t, "Updated Name", updatedUser.Name)

	// Delete user (soft delete)
	err = db.Delete(&updatedUser).Error
	require.NoError(t, err)

	// Verify soft delete
	var deletedUser models.User
	err = db.First(&deletedUser, user.ID).Error
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)

	// But should exist with Unscoped
	err = db.Unscoped().First(&deletedUser, user.ID).Error
	require.NoError(t, err)
	assert.NotZero(t, deletedUser.DeletedAt)
}

func TestUserUniqueEmail(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.User{})
	require.NoError(t, err)

	// Create first user
	user1 := models.User{
		Name:  "User 1",
		Email: "unique@example.com",
	}
	err = db.Create(&user1).Error
	require.NoError(t, err)

	// Try to create user with same email
	user2 := models.User{
		Name:  "User 2",
		Email: "unique@example.com",
	}
	err = db.Create(&user2).Error
	assert.Error(t, err) // Should fail due to unique constraint
}
