package services

import (
	"context"

	"github.com/weni/whatsapp-router/logger"
	"github.com/weni/whatsapp-router/metric"
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
	CreateChannelDefault(*models.Channel) (*models.Channel, error)
}

type DefaultChannelService struct {
	repo    repositories.ChannelRepository
	Metrics *metric.Service
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
	token := utils.GenToken()
	channel.Token = token
	err := s.repo.Insert(&channel)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	channelCreationMetric := metric.NewChannelCreation(channel.UUID)
	s.Metrics.SaveChannelCreation(channelCreationMetric)
	return &pb.ChannelResponse{
		Token: channel.Token,
	}, nil
}

func (s DefaultChannelService) CreateChannelDefault(channel *models.Channel) (*models.Channel, error) {
	err := s.repo.Insert(channel)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	channelCreationMetric := metric.NewChannelCreation(channel.UUID)
	s.Metrics.SaveChannelCreation(channelCreationMetric)
	return channel, nil
}

func NewChannelService(repo repositories.ChannelRepository, metricService *metric.Service) DefaultChannelService {
	return DefaultChannelService{repo, metricService}
}
