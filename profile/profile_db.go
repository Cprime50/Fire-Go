package profile

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/cprime50/fire-go/db"
	"github.com/google/uuid"
)

type Profile struct {
	Id        string    `json:"id"`
	UserId    string    `json:"user_id"`
	Email     string    `json:"email"`
	UserName  string    `json:"username"`
	Bio       string    `json:"bio"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var (
	ErrQuoteNotFound             = errors.New("quote not found")
	ErrProfileNotFound           = errors.New("profile not found")
	ErrDuplicateEntry            = errors.New("duplicate entry")
	ErrForeignKeyViolation       = errors.New("foreign key violation")
	ErrUniqueConstraintViolation = errors.New("unique constraint violation")
)

func createProfile(p *Profile) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return fmt.Errorf("uuid.NewRandom: %w", err)
	}
	_, err = db.Db.Exec(
		"INSERT INTO profiles (id, user_id, email, username, bio, created_at) VALUES ($1, $2, $3, $4, $5, $6)",
		id.String(),
		p.UserId,
		p.Email,
		p.UserName,
		p.Bio,
		time.Now(),
	)
	if err != nil {
		return fmt.Errorf("CreateProfile error: %w", err)
	}
	return nil
}

// GetProfileByUserId retrieves a user profile by user ID.
func getProfileByUserId(userId string) (*Profile, error) {
	profile := &Profile{}
	var createdAt, updatedAt time.Time
	err := db.Db.QueryRow("SELECT id, user_id, email, username, bio, created_at, updated_at FROM profiles WHERE user_id = $1", userId).
		Scan(&profile.Id, &profile.UserId, &profile.Email, &profile.UserName, &profile.Bio, &createdAt, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrProfileNotFound
		}
		return nil, fmt.Errorf("GetProfileByUserId: %w", err)
	}
	profile.CreatedAt = createdAt
	profile.UpdatedAt = updatedAt
	return profile, nil
}

func updateProdile(p *Profile) error {
	result, err := db.Db.Exec(
		"UPDATE profiles SET bio = $1, username = $2, updated_at = $3 WHERE user_id = $4",
		p.Bio,
		p.UserName,
		time.Now(),
		p.UserId,
	)
	if err != nil {
		return fmt.Errorf("UpdateProfile error: %w", err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrProfileNotFound
	}
	return nil
}

// DeleteProfile deletes a user profile by user ID.
func deleteProfile(userId string) error {
	_, err := db.Db.Exec("DELETE FROM profiles WHERE user_id = $1", userId)
	if err != nil {
		return fmt.Errorf("DeleteProfile error: %w", err)
	}
	return nil
}

// GetAllProfiles retrieves all user profiles.
func getAllProfiles() ([]*Profile, error) {
	rows, err := db.Db.Query("SELECT * FROM profiles")
	if err != nil {
		return nil, fmt.Errorf("GetAllProfiles error: %w", err)
	}
	defer rows.Close()

	var profiles []*Profile
	for rows.Next() {
		profile := &Profile{}
		var createdAt, updatedAt time.Time
		if err := rows.Scan(&profile.Id, &profile.UserId, &profile.Email, &profile.UserName, &profile.Bio, &createdAt, &updatedAt); err != nil {
			return nil, fmt.Errorf("GetAllProfiles rows.Scan: %w", err)
		}
		profile.CreatedAt = createdAt
		profile.UpdatedAt = updatedAt
		profiles = append(profiles, profile)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("GetAllProfiles rows.Err: %w", err)
	}

	if len(profiles) == 0 {
		return nil, ErrProfileNotFound
	}

	return profiles, nil
}
