package repositories

import (
	"context"

	"github.com/weni/whatsapp-router/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const CHANNEL_COLLECTION = "channel"

type ChannelRepository interface {
	Insert(*models.Channel) error
	FindOne(*models.Channel) (*models.Channel, error)
	FindById(string) (*models.Channel, error)
	FindByToken(string) (*models.Channel, error)
}

type ChannelRepositoryDb struct {
	DB *mongo.Database
}

func (c ChannelRepositoryDb) Insert(channel *models.Channel) error {
	result, err := c.DB.Collection(CHANNEL_COLLECTION).InsertOne(context.TODO(), channel)
	if err != nil {
		return err
	}
	if id, ok := result.InsertedID.(primitive.ObjectID); ok {
		channel.ID = id
	}
	return nil
}

func (c ChannelRepositoryDb) FindOne(channel *models.Channel) (*models.Channel, error) {
	var ch models.Channel
	qry := bson.M{
		"uuid": channel.UUID,
	}
	if err := c.DB.Collection(CHANNEL_COLLECTION).FindOne(context.TODO(), qry).Decode(&ch); err != nil {
		return nil, err
	}
	return &ch, nil
}

func (c ChannelRepositoryDb) FindById(id string) (*models.Channel, error) {
	var ch models.Channel
	objId, _ := primitive.ObjectIDFromHex(id)
	qry := bson.M{
		"_id": objId,
	}
	if err := c.DB.Collection(CHANNEL_COLLECTION).FindOne(context.TODO(), qry).Decode(&ch); err != nil {
		return nil, err
	}
	return &ch, nil
}

func (c ChannelRepositoryDb) FindByToken(token string) (*models.Channel, error) {
	var ch models.Channel
	qry := bson.M{
		"token": token,
	}
	if err := c.DB.Collection(CHANNEL_COLLECTION).FindOne(context.TODO(), qry).Decode(&ch); err != nil {
		return nil, err
	}
	return &ch, nil
}

func NewChannelRepositoryDb(dbClient *mongo.Database) ChannelRepositoryDb {
	return ChannelRepositoryDb{dbClient}
}
