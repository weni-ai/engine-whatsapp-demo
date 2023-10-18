package services

import (
	"github.com/weni/whatsapp-router/models"
	"github.com/weni/whatsapp-router/repositories"
)

type FlowsService interface {
	FindFlows(*models.Flows) (*models.Flows, error)
	CreateFlows(*models.Flows) (*models.Flows, error)
	UpdateFlows(*models.Flows) (*models.Flows, error)
}

type DefaultFlowsService struct {
	repo repositories.FlowsRepository
}

func (s DefaultFlowsService) FindFlows(req *models.Flows) (*models.Flows, error) {
	f, err := s.repo.FindOne(req)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (s DefaultFlowsService) CreateFlows(req *models.Flows) (*models.Flows, error) {
	f := &models.Flows{
		FlowsStarts: req.FlowsStarts,
		Channel:     req.Channel,
	}

	newFlows, err := s.repo.Insert(f)
	if err != nil {
		return nil, err
	}
	return newFlows, nil
}

func (s DefaultFlowsService) UpdateFlows(req *models.Flows) (*models.Flows, error) {
	f := &models.Flows{
		FlowsStarts: req.FlowsStarts,
		Channel:     req.Channel,
	}
	updatedFlows, err := s.repo.Update(f)
	if err != nil {
		return nil, err
	}
	return updatedFlows, nil
}

func NewFlowsService(repo repositories.FlowsRepository) DefaultFlowsService {
	return DefaultFlowsService{repo}
}
