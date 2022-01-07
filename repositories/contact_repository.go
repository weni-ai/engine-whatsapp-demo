package repositories

import (
	"context"
	"errors"

	"github.com/weni/whatsapp-router/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const CONTACT_COLLECTION = "contact"

type ContactRepository interface {
	Insert(contact *models.Contact) (*models.Contact, error)
	FindOne(contact *models.Contact) (*models.Contact, error)
	Update(contact *models.Contact) (*models.Contact, error)
}

type ContactRepositoryDb struct {
	DB *mongo.Database
}

func (c ContactRepositoryDb) Insert(contact *models.Contact) (*models.Contact, error) {
	result, err := c.DB.Collection(CONTACT_COLLECTION).InsertOne(context.TODO(), contact)
	if err != nil {
		return nil, errors.New("unexpected database error - " + err.Error())
	}

	if id, ok := result.InsertedID.(primitive.ObjectID); ok {
		contact.ID = id
	}

	return contact, nil
}

func (c ContactRepositoryDb) FindOne(contact *models.Contact) (*models.Contact, error) {
	var cont models.Contact
	qry := bson.M{
		"urn": contact.URN,
	}
	if err := c.DB.Collection(CONTACT_COLLECTION).FindOne(context.TODO(), qry).Decode(&cont); err != nil {
		return nil, errors.New("contact not found")
	}
	return &cont, nil
}

func (c ContactRepositoryDb) Update(contact *models.Contact) (*models.Contact, error) {
	q := bson.M{
		"urn": contact.URN,
	}
	d, err := bson.Marshal(contact)
	if err != nil {
		return nil, errors.New("internal server error: " + err.Error())
	}
	dc := bson.D{}
	err = bson.Unmarshal(d, &dc)
	if err != nil {
		return nil, errors.New("internal server error: " + err.Error())
	}
	result, err := c.DB.Collection(CONTACT_COLLECTION).UpdateOne(
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
		contact.ID = id
	}
	return contact, nil
}

func NewContactRepositoryDb(dbClient *mongo.Database) ContactRepositoryDb {
	return ContactRepositoryDb{dbClient}
}
