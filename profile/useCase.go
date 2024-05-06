package profile

import (
	"errors"
	"log"
	"time"
)

type ProfileService interface {
	CreateProfile(userID, email string) (*ProfileResponse, error)
	UpdateProfile(userID, bio, username string) (*ProfileResponse, error)
	DeleteProfile(userID string, role string, profileId string) error
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
		}
		log.Printf("Error updating profile: %v", err)
		return nil, ErrUpdateProfile

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

func (s *ProfileServiceImpl) DeleteProfile(userID string, role string, profileId string) error {
	err := deleteProfile(userID)
	if err != nil {
		if errors.Is(err, ErrProfileNotFound) {
			log.Printf("Error deleting profile: %v", err)
			return ErrProfileNotFound
		}
		log.Printf("Error deleting profile: %v", err)
		return ErrDeletingProfile

	}
	if role != "admin" {
		if userID != profileId {
			log.Printf("DeleteProfile: Error User with id %s and role %s not allowed access to delete user with id %s", userID, role, profileId)
			return ErrNotAuthorized
		}
	}

	log.Printf("DeleteProfile: Profile deleted successfully for userID %s", userID)
	return nil
}

func (s *ProfileServiceImpl) GetProfile(userID string) (*Profile, error) {
	profile, err := getProfileByUserId(userID)
	if err != nil {
		if errors.Is(err, ErrProfileNotFound) {
			log.Printf("GetProfile: Profile not found for userID %s", userID)
			return nil, ErrProfileNotFound
		}
		log.Printf("GetProfile: Error retrieving profile for userID %s: %v", userID, err)
		return nil, ErrGettingProfile
	}

	return profile, nil
}

func (s *ProfileServiceImpl) GetAllProfiles() ([]*Profile, error) {
	profiles, err := getAllProfiles()
	if err != nil {
		if errors.Is(err, ErrProfileNotFound) {
			log.Print("GetAllProfiles: No profiles found")
			return nil, ErrProfileNotFound
		}
		log.Printf("GetAllProfiles: Database error: %v", err)
		return nil, ErrGettingProfile
	}

	return profiles, nil
}
