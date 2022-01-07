package services

import (
	"github.com/weni/whatsapp-router/models"
	"github.com/weni/whatsapp-router/repositories"
)

type ConfigService interface {
	CreateOrUpdate(*models.Config) (*models.Config, error)
	GetConfig() (*models.Config, error)
}

type configService struct {
	repo repositories.ConfigRepository
}

func (s *configService) GetConfig() (*models.Config, error) {
	return s.repo.GetFirst()
}

func (s *configService) CreateOrUpdate(conf *models.Config) (*models.Config, error) {
	cf, _ := s.repo.GetFirst()
	if cf == nil {
		c := &models.Config{
			Token: conf.Token,
		}
		err := s.repo.Create(c)
		if err != nil {
			return nil, err
		}
		return s.repo.GetFirst()
	} else {
		cf.Token = conf.Token
		return s.repo.Update(cf)
	}
}

func NewConfigService(repo repositories.ConfigRepository) ConfigService {
	return &configService{repo}
}
