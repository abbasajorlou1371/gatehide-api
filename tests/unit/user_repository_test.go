package unit

import (
	"testing"

	"github.com/gatehide/gatehide-api/internal/repositories"
	testutils "github.com/gatehide/gatehide-api/tests/utils"
)

func TestUserRepository_GetByEmail(t *testing.T) {
	testutils.SkipIfNoDB(t)

	db := testutils.SetupTestDB(t)
	defer testutils.CleanupTestDB(t, db)
	defer db.Close()

	userRepo := repositories.NewUserRepository(db)

	// Create a test user
	testUser := testutils.CreateTestUser(t, db, "test@example.com", "password123", "Test User")

	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{
			name:    "existing user",
			email:   "test@example.com",
			wantErr: false,
		},
		{
			name:    "non-existing user",
			email:   "nonexistent@example.com",
			wantErr: true,
		},
		{
			name:    "empty email",
			email:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := userRepo.GetByEmail(tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserRepository.GetByEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if user == nil {
					t.Error("UserRepository.GetByEmail() returned nil user")
					return
				}

				if user.Email != tt.email {
					t.Errorf("UserRepository.GetByEmail() user.Email = %v, want %v", user.Email, tt.email)
				}

				if user.ID != testUser.ID {
					t.Errorf("UserRepository.GetByEmail() user.ID = %v, want %v", user.ID, testUser.ID)
				}

				if user.Name != testUser.Name {
					t.Errorf("UserRepository.GetByEmail() user.Name = %v, want %v", user.Name, testUser.Name)
				}
			}
		})
	}
}

func TestUserRepository_UpdateLastLogin(t *testing.T) {
	testutils.SkipIfNoDB(t)

	db := testutils.SetupTestDB(t)
	defer testutils.CleanupTestDB(t, db)
	defer db.Close()

	userRepo := repositories.NewUserRepository(db)

	// Create a test user
	testUser := testutils.CreateTestUser(t, db, "test@example.com", "password123", "Test User")

	tests := []struct {
		name    string
		userID  int
		wantErr bool
	}{
		{
			name:    "existing user",
			userID:  testUser.ID,
			wantErr: false,
		},
		{
			name:    "non-existing user",
			userID:  99999,
			wantErr: false, // Update doesn't fail for non-existing users, just affects 0 rows
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := userRepo.UpdateLastLogin(tt.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserRepository.UpdateLastLogin() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAdminRepository_GetByEmail(t *testing.T) {
	testutils.SkipIfNoDB(t)

	db := testutils.SetupTestDB(t)
	defer testutils.CleanupTestDB(t, db)
	defer db.Close()

	adminRepo := repositories.NewAdminRepository(db)

	// Create a test admin
	testAdmin := testutils.CreateTestAdmin(t, db, "admin@example.com", "admin123", "Test Admin")

	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{
			name:    "existing admin",
			email:   "admin@example.com",
			wantErr: false,
		},
		{
			name:    "non-existing admin",
			email:   "nonexistent@example.com",
			wantErr: true,
		},
		{
			name:    "empty email",
			email:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			admin, err := adminRepo.GetByEmail(tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("AdminRepository.GetByEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if admin == nil {
					t.Error("AdminRepository.GetByEmail() returned nil admin")
					return
				}

				if admin.Email != tt.email {
					t.Errorf("AdminRepository.GetByEmail() admin.Email = %v, want %v", admin.Email, tt.email)
				}

				if admin.ID != testAdmin.ID {
					t.Errorf("AdminRepository.GetByEmail() admin.ID = %v, want %v", admin.ID, testAdmin.ID)
				}

				if admin.Name != testAdmin.Name {
					t.Errorf("AdminRepository.GetByEmail() admin.Name = %v, want %v", admin.Name, testAdmin.Name)
				}
			}
		})
	}
}

func TestAdminRepository_UpdateLastLogin(t *testing.T) {
	testutils.SkipIfNoDB(t)

	db := testutils.SetupTestDB(t)
	defer testutils.CleanupTestDB(t, db)
	defer db.Close()

	adminRepo := repositories.NewAdminRepository(db)

	// Create a test admin
	testAdmin := testutils.CreateTestAdmin(t, db, "admin@example.com", "admin123", "Test Admin")

	tests := []struct {
		name    string
		adminID int
		wantErr bool
	}{
		{
			name:    "existing admin",
			adminID: testAdmin.ID,
			wantErr: false,
		},
		{
			name:    "non-existing admin",
			adminID: 99999,
			wantErr: false, // Update doesn't fail for non-existing admins, just affects 0 rows
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adminRepo.UpdateLastLogin(tt.adminID)
			if (err != nil) != tt.wantErr {
				t.Errorf("AdminRepository.UpdateLastLogin() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserRepository_EmailUniqueness(t *testing.T) {
	testutils.SkipIfNoDB(t)

	db := testutils.SetupTestDB(t)
	defer testutils.CleanupTestDB(t, db)
	defer db.Close()

	userRepo := repositories.NewUserRepository(db)

	// Create first user
	testutils.CreateTestUser(t, db, "unique@example.com", "password123", "First User")

	// Try to create another user with the same email
	_, err := userRepo.GetByEmail("unique@example.com")
	if err != nil {
		t.Errorf("UserRepository.GetByEmail() failed for existing user: %v", err)
	}

	// Create a user with different email
	testutils.CreateTestUser(t, db, "different@example.com", "password123", "Second User")

	// Both users should be retrievable
	user1, err := userRepo.GetByEmail("unique@example.com")
	if err != nil {
		t.Errorf("UserRepository.GetByEmail() failed for first user: %v", err)
	}

	user2, err := userRepo.GetByEmail("different@example.com")
	if err != nil {
		t.Errorf("UserRepository.GetByEmail() failed for second user: %v", err)
	}

	if user1.ID == user2.ID {
		t.Error("Users with different emails should have different IDs")
	}
}

func TestAdminRepository_EmailUniqueness(t *testing.T) {
	testutils.SkipIfNoDB(t)

	db := testutils.SetupTestDB(t)
	defer testutils.CleanupTestDB(t, db)
	defer db.Close()

	adminRepo := repositories.NewAdminRepository(db)

	// Create first admin
	testutils.CreateTestAdmin(t, db, "admin@example.com", "admin123", "First Admin")

	// Try to create another admin with the same email
	_, err := adminRepo.GetByEmail("admin@example.com")
	if err != nil {
		t.Errorf("AdminRepository.GetByEmail() failed for existing admin: %v", err)
	}

	// Create an admin with different email
	testutils.CreateTestAdmin(t, db, "admin2@example.com", "admin123", "Second Admin")

	// Both admins should be retrievable
	admin1, err := adminRepo.GetByEmail("admin@example.com")
	if err != nil {
		t.Errorf("AdminRepository.GetByEmail() failed for first admin: %v", err)
	}

	admin2, err := adminRepo.GetByEmail("admin2@example.com")
	if err != nil {
		t.Errorf("AdminRepository.GetByEmail() failed for second admin: %v", err)
	}

	if admin1.ID == admin2.ID {
		t.Error("Admins with different emails should have different IDs")
	}
}
