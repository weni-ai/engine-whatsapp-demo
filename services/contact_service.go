package services

import (
	"github.com/weni/whatsapp-router/models"
	"github.com/weni/whatsapp-router/repositories"
)

type ContactService interface {
	FindContact(*models.Contact) (*models.Contact, error)
	CreateContact(*models.Contact) (*models.Contact, error)
}

type DefaultContactService struct {
	repo repositories.ContactRepository
}

func (s DefaultContactService) FindContact(req *models.Contact) (*models.Contact, error) {
	c, err := s.repo.FindOne(req)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (s DefaultContactService) CreateContact(req *models.Contact) (*models.Contact, error) {
	c := &models.Contact{
		URN:     req.URN,
		Name:    req.Name,
		Channel: req.Channel,
	}

	newContact, err := s.repo.Insert(c)
	if err != nil {
		return nil, err
	}
	return newContact, nil
}

func NewContactService(repo repositories.ContactRepository) DefaultContactService {
	return DefaultContactService{repo}
}
