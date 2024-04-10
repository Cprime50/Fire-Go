package profile

import (
	"errors"
	"log"
	"os"
	"testing"

	"github.com/cprime50/fire-go/db"
)

func TestMain(m *testing.M) {

	log.Println("Running tests...")
	Db, err := db.ConnectTest()
	if err != nil {
		log.Fatal(err)
	}
	err = db.Migrate(Db)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Db.Close()

	os.Exit(m.Run())
}

func TestCreateProfile(t *testing.T) {
	clearProfiles()

	// Test case 1: Insert a valid profile
	profile := &Profile{
		UserId:   "test1",
		Email:    "test1@email.com",
		UserName: "Username1",
		Bio:      "test bio 1",
	}
	err := createProfile(profile)
	if err != nil {
		t.Errorf("createProfile error: %v", err)
	}

	// Test case 2: Insert a second valid profile
	profile2 := &Profile{
		UserId:   "test2",
		Email:    "test2@email.com",
		UserName: "Username2",
		Bio:      "test bio 2",
	}
	err = createProfile(profile2)
	if err != nil {
		t.Errorf("createProfile error: %v", err)
	}

	// Test case 3: Insert a profile that already exists
	err = createProfile(profile)
	if err == nil {
		t.Errorf("creating duplicate profile error: %v", err)
	}
}

func TestGetProfileByUserId(t *testing.T) {
	clearProfiles()

	// Test case 1: Select a profile by id
	profile := &Profile{
		UserId:   "test1",
		Email:    "test1@email.com",
		UserName: "Username1",
		Bio:      "test bio 1",
	}
	err := createProfile(profile)
	if err != nil {
		t.Errorf("createProfile error: %v", err)
	}
	gottenProfile, err := getProfileByUserId(profile.UserId)
	if err != nil {
		t.Errorf("getProfileByUserId error: %v", err)
	}
	if gottenProfile.UserId != profile.UserId {
		t.Errorf("getProfileByUserId error: not equal")
	}

	// Test case 2: Select a profile by id that does not exist
	_, err = getProfileByUserId("not_exist")
	if !errors.Is(err, ErrProfileNotFound) {
		t.Errorf("getProfileByUserId error: %v", err)
	}
}

func TestUpdateProfile(t *testing.T) {
	clearProfiles()

	// Test case 1: Update a valid profile
	profile := &Profile{
		UserId:   "test1",
		Email:    "test1@email.com",
		UserName: "Username1",
		Bio:      "test bio 1",
	}
	err := createProfile(profile)
	if err != nil {
		t.Errorf("createProfile error: %v", err)
	}
	newProfile := &Profile{
		UserId:   profile.UserId,
		UserName: "New Username",
		Bio:      "New Bio",
	}
	err = updateProfile(newProfile)
	if err != nil {
		t.Errorf("updateProfile error: %v", err)
	}
	updatedProfile, _ := getProfileByUserId(profile.UserId)
	if updatedProfile.UserName != newProfile.UserName || updatedProfile.Bio != newProfile.Bio {
		t.Errorf("updateProfile error: not equal")
	}

	// Test case 2: Update a profile that does not exist
	newProfile = &Profile{
		UserId: "not_exist",
	}
	err = updateProfile(newProfile)
	if !errors.Is(err, ErrProfileNotFound) {
		t.Errorf("updateProfile error: %v", err)
	}
}

func TestDeleteProfile(t *testing.T) {
	clearProfiles()

	// Test case 1: Delete an existing profile
	profile := &Profile{
		UserId:   "test1",
		Email:    "test1@email.com",
		UserName: "Username1",
		Bio:      "test bio 1",
	}
	err := createProfile(profile)
	if err != nil {
		t.Errorf("createProfile error: %v", err)
	}
	err = deleteProfile(profile.UserId)
	if err != nil {
		t.Errorf("deleteProfile error: %v", err)
	}

	// Test case 2: Delete a profile that does not exist
	err = deleteProfile("not_exist")
	if err != nil {
		t.Error("Error, deleting non existent profile error")
	}
}

func TestGetAllProfiles(t *testing.T) {
	clearProfiles()

	// Insert profiles
	profile1 := &Profile{
		UserId:   "test1",
		Email:    "test1@email.com",
		UserName: "Username1",
		Bio:      "test bio 1",
	}
	_ = createProfile(profile1)
	profile2 := &Profile{
		UserId:   "test2",
		Email:    "test2@email.com",
		UserName: "Username2",
		Bio:      "test bio 2",
	}
	_ = createProfile(profile2)

	// Get profiles
	gottenProfiles, err := getAllProfiles()
	if err != nil {
		t.Errorf("getAllProfiles error: %v", err)
	}
	if len(gottenProfiles) != 2 {
		t.Errorf("getAllProfiles error: expected 2 profiles, got %d", len(gottenProfiles))
	}
}

func clearProfiles() {
	_, err := db.Db.Exec("DELETE FROM profiles")
	if err != nil {
		log.Fatal(err)
	}
}
