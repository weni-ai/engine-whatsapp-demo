package repositories

import (
	"context"

	"github.com/weni/whatsapp-router/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const CONTACT_COLLECTION = "contact"

type ContactRepository interface {
	Insert(contact *models.Contact) (*models.Contact, error)
	FindOne(contact *models.Contact) (*models.Contact, error)
}

type ContactRepositoryDb struct {
	DB *mongo.Database
}

func (c ContactRepositoryDb) Insert(contact *models.Contact) (*models.Contact, error) {
	result, err := c.DB.Collection(CONTACT_COLLECTION).InsertOne(context.TODO(), contact)
	if err != nil {
		return nil, err
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
		return nil, err
	}
	return &cont, nil
}

func NewContactRepositoryDb(dbClient *mongo.Database) ContactRepositoryDb {
	return ContactRepositoryDb{dbClient}
}
