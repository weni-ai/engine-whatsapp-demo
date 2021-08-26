package services

import (
	"github.com/weni/whatsapp-router/models"
	"github.com/weni/whatsapp-router/repositories"
)

type ChannelService interface {
	FindChannel(*models.Channel) (*models.Channel, error)
	FindChannelById(string) (*models.Channel, error)
	FindChannelByToken(string) (*models.Channel, error)
}

type DefaultChannelService struct {
	repo repositories.ChannelRepository
}

func (s DefaultChannelService) FindChannel(req *models.Channel) (*models.Channel, error) {
	return nil, nil
}

func (s DefaultChannelService) FindChannelById(req string) (*models.Channel, error) {
	ch, err := s.repo.FindById(req)
	if err != nil {
		return nil, err
	}
	return ch, nil
}

func (s DefaultChannelService) FindChannelByToken(req string) (*models.Channel, error) {
	ch, err := s.repo.FindByToken(req)
	if err != nil {
		return nil, err
	}
	return ch, nil
}

func NewChannelService(repo repositories.ChannelRepository) DefaultChannelService {
	return DefaultChannelService{repo}
}
