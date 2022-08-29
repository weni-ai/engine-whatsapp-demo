package repositories

import (
	"context"
	"errors"

	"github.com/weni/whatsapp-router/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const FLOWS_COLLECTION = "flows"

type FlowsRepository interface {
	Insert(flows *models.Flows) (*models.Flows, error)
	FindOne(flows *models.Flows) (*models.Flows, error)
	Update(flows *models.Flows) (*models.Flows, error)
}

type FlowsRepositoryDb struct {
	DB *mongo.Database
}

func (f FlowsRepositoryDb) Insert(flows *models.Flows) (*models.Flows, error) {
	result, err := f.DB.Collection(FLOWS_COLLECTION).InsertOne(context.TODO(), flows)
	if err != nil {
		return nil, errors.New("unexpected database error - " + err.Error())
	}

	if id, ok := result.InsertedID.(primitive.ObjectID); ok {
		flows.ID = id
	}

	return flows, nil
}

func (f FlowsRepositoryDb) FindOne(flows *models.Flows) (*models.Flows, error) {
	var fl models.Flows
	qry := bson.M{
		"channel_uuid": flows.Channel,
	}
	if err := f.DB.Collection(FLOWS_COLLECTION).FindOne(context.TODO(), qry).Decode(&fl); err != nil {
		return nil, errors.New("flows not found")
	}
	return &fl, nil
}

func (f FlowsRepositoryDb) Update(flows *models.Flows) (*models.Flows, error) {
	q := bson.M{
		"flows_starts": flows.FlowsStarts,
	}
	d, err := bson.Marshal(flows)
	if err != nil {
		return nil, errors.New("internal server error: " + err.Error())
	}
	dc := bson.D{}
	err = bson.Unmarshal(d, &dc)
	if err != nil {
		return nil, errors.New("internal server error: " + err.Error())
	}
	result, err := f.DB.Collection(FLOWS_COLLECTION).UpdateOne(
		context.TODO(),
		q,
		bson.D{
			{Key: "$set", Value: dc},
		},
	)
	if err != nil {
		return nil, errors.New("unexpected database error - " + err.Error())
	}

	if id, ok := result.UpsertedID.(primitive.ObjectID); ok {
		flows.ID = id
	}
	return flows, nil
}

func NewFlowsRepositoryDb(dbClient *mongo.Database) FlowsRepositoryDb {
	return FlowsRepositoryDb{dbClient}
}
