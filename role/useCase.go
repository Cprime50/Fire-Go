package role

import (
	"context"
	"log"

	"firebase.google.com/go/v4/auth"
	"github.com/cprime50/fire-go/middleware"
)

type AdminService interface {
	MakeAdmin(email string) error
	RemoveAdmin(email string) error
}

type AdminServiceImpl struct {
	client *auth.Client
}

func NewAdminService(client *auth.Client) *AdminServiceImpl {
	return &AdminServiceImpl{client: client}
}

func (s *AdminServiceImpl) MakeAdmin(email string) error {
	if err := middleware.AssignRole(context.Background(), s.client, email, "admin"); err != nil {
		log.Printf("Error assigning admin role: %v", err)
		return err
	}
	return nil
}

func (s *AdminServiceImpl) RemoveAdmin(email string) error {
	if err := middleware.AssignRole(context.Background(), s.client, email, "user"); err != nil {
		log.Printf("Error assigning user role: %v", err)
		return err
	}
	return nil
}
