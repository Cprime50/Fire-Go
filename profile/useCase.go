package profile

import (
	"errors"
	"log"
	"time"
)

var (
// ErrProfileNotFound = errors.New("profile not found")
)

type ProfileService interface {
	CreateProfile(userID, email string) (*ProfileResponse, error)
	UpdateProfile(userID, bio, username string) (*ProfileResponse, error)
	DeleteProfile(userID string) error
	GetProfile(userID string) (*Profile, error)
	GetAllProfiles() ([]Profile, error)
}

type ProfileServiceImpl struct{}

func (s *ProfileServiceImpl) CreateProfile(userID, email string) (*ProfileResponse, error) {
	existingProfile, err := getProfileByUserId(userID)
	if err == nil && existingProfile != nil {
		log.Printf("Profile already exists for user %s with email %s", userID, email)
		return nil, ErrProfileAlreadyExists
	} else if err != nil && !errors.Is(err, ErrProfileNotFound) {
		log.Printf("Error checking profile existence: %v", err)
		return nil, ErrCreateProfile
	}

	username, err := generateUsername(email)
	if err != nil {
		log.Printf("Error generating username: %v", err)
		return nil, ErrCreateProfile
	}
	err = createProfile(&Profile{
		UserId:   userID,
		Email:    email,
		UserName: username,
		Bio:      "",
	})
	if err != nil {
		log.Printf("Error creating profile: %v", err)
		return nil, ErrCreateProfile
	}

	createdProfile, err := getProfileByUserId(userID)
	if err != nil {
		log.Printf("Error retrieving created profile: %v", err)
		return nil, ErrCreateProfile
	}

	response := &ProfileResponse{
		Profile: createdProfile,
		Message: "Profile created successfully",
	}

	log.Printf("Profile created successfully for user %s with email %s", userID, email)
	return response, nil
}

func (s *ProfileServiceImpl) UpdateProfile(userID, bio, username string) (*ProfileResponse, error) {

	err := updateProfile(&Profile{
		UserId:    userID,
		Bio:       bio,
		UserName:  username,
		UpdatedAt: time.Now(),
	})
	if err != nil {
		if errors.Is(err, ErrProfileNotFound) {
			log.Printf("Error updating profile: %v", err)
			return nil, ErrProfileNotFound
		} else {
			log.Printf("Error updating profile: %v", err)
			return nil, ErrUpdateProfile
		}

	}

	updatedProfile, err := getProfileByUserId(userID)
	if err != nil {
		log.Printf("Error retrieving profile: %v", err)
		return nil, ErrProfileNotFound
	}

	response := &ProfileResponse{
		Profile: updatedProfile,
		Message: "Profile updated successfully",
	}

	log.Printf("Profile updated successfully for user %s", userID)
	return response, nil
}

func (s *ProfileServiceImpl) DeleteProfile(userID string) error {
	// Implement the logic to delete a profile
	return nil
}

func (s *ProfileServiceImpl) GetProfile(userID string) (*Profile, error) {
	// Implement the logic to get a profile
	return nil, ErrProfileNotFound
}

func (s *ProfileServiceImpl) GetAllProfiles() ([]Profile, error) {
	// Implement the logic to get all profiles
	return nil, ErrProfileNotFound
}
