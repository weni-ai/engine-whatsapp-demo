package services

import (
	"context"

	"github.com/weni/whatsapp-router/models"
	"github.com/weni/whatsapp-router/repositories"
	"github.com/weni/whatsapp-router/servers/grpc/pb"
	"github.com/weni/whatsapp-router/utils"
)

type ChannelService interface {
	FindChannel(*models.Channel) (*models.Channel, error)
	FindChannelById(string) (*models.Channel, error)
	FindChannelByToken(string) (*models.Channel, error)
	CreateChannel(context.Context, *pb.ChannelRequest) (*pb.ChannelResponse, error)
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

func (s DefaultChannelService) CreateChannel(ctx context.Context, req *pb.ChannelRequest) (*pb.ChannelResponse, error) {
	var channel models.Channel
	channel.UUID = req.GetUuid()
	channel.Name = req.GetName()
	token := utils.GenToken(channel.Name)
	channel.Token = token
	s.repo.Insert(&channel)
	return &pb.ChannelResponse{
		Token: channel.Token,
	}, nil
}

func NewChannelService(repo repositories.ChannelRepository) DefaultChannelService {
	return DefaultChannelService{repo}
}
