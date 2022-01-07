package repositories

import (
	"context"
	"errors"

	"github.com/weni/whatsapp-router/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const CONFIG_COLLECTION = "config"

type ConfigRepository interface {
	Create(config *models.Config) error
	GetFirst() (*models.Config, error)
	Update(config *models.Config) (*models.Config, error)
	FindOne(config *models.Config) (*models.Config, error)
}

type configRepository struct {
	DB *mongo.Database
}

func (c configRepository) Create(config *models.Config) error {
	_, err := c.DB.Collection(CONFIG_COLLECTION).InsertOne(context.TODO(), config)
	if err != nil {
		return errors.New("unexpected database error - " + err.Error())
	}
	return nil
}

func (c configRepository) GetFirst() (*models.Config, error) {
	cursor, err := c.DB.Collection(CONFIG_COLLECTION).Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	cursor.Next(context.TODO())
	var config *models.Config
	if err = cursor.Decode(&config); err != nil {
		return nil, nil
	}
	return config, nil
}

func (c *configRepository) FindOne(config *models.Config) (*models.Config, error) {
	var conf models.Config
	q := bson.M{
		"_id": config.ID,
	}
	if err := c.DB.Collection(CONFIG_COLLECTION).FindOne(context.TODO(), q).Decode(&conf); err != nil {
		return nil, errors.New("config not found")
	}
	return &conf, nil
}

func (c *configRepository) Update(config *models.Config) (*models.Config, error) {
	q := bson.M{
		"_id": config.ID,
	}
	d, err := bson.Marshal(config)
	if err != nil {
		return nil, errors.New("internal server error: " + err.Error())
	}
	dc := bson.D{}
	err = bson.Unmarshal(d, &dc)
	if err != nil {
		return nil, errors.New("internal server error: " + err.Error())
	}
	_, err = c.DB.Collection(CONFIG_COLLECTION).UpdateOne(
		context.TODO(),
		q,
		bson.D{
			{Key: "$set", Value: dc},
		},
	)
	if err != nil {
		return nil, errors.New("unexpected database error - " + err.Error())
	}

	return config, nil
}

func NewConfigRepository(dbClient *mongo.Database) ConfigRepository {
	return &configRepository{dbClient}
}
