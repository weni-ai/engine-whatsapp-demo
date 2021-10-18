package repositories

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/weni/whatsapp-router/models"
	"github.com/weni/whatsapp-router/storage"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var dummyChannel = models.Channel{
	ID:    primitive.NewObjectID(),
	UUID:  "f11c744c-4937-4ee3-8a51-26e56eb77c4e",
	Name:  "foo",
	Token: "foo-bar-zaz",
}

var tcInsertContact = []struct {
	TestName string
	Contact  models.Contact
	Err      error
}{
	{
		TestName: "test insert contact",
		Contact: models.Contact{
			URN:     "5582988887777",
			Name:    "dummy",
			Channel: dummyChannel.ID,
		},
	},
}

func TestInsertContact(t *testing.T) {
	mongodb := storage.NewTestDB()
	defer storage.CloseDB(mongodb)
	contactRepository := ContactRepositoryDb{DB: mongodb}

	for _, tc := range tcInsertContact {
		t.Run(tc.TestName, func(t *testing.T) {
			c, err := contactRepository.Insert(&tc.Contact)
			if fmt.Sprint(err) != fmt.Sprint(tc.Err) {
				t.Errorf("got %v / want %v", err, tc.Err)
			}

			if c == nil {
				t.Errorf("got %v / want %v", c, reflect.TypeOf(tc.Contact))
			}
		})
	}
}

var tcFindOneContact = []struct {
	TestName string
	Contact  models.Contact
	Err      error
}{
	{
		TestName: "Find one existing contact",
		Contact: models.Contact{
			URN:     "5582988887777",
			Name:    "dummy",
			Channel: dummyChannel.ID,
		},
	},
}

func TestFindOneContact(t *testing.T) {
	mongodb := storage.NewTestDB()
	defer storage.CloseDB(mongodb)
	contactRepository := ContactRepositoryDb{DB: mongodb}

	for _, tc := range tcFindOneContact {
		t.Run(tc.TestName, func(t *testing.T) {
			c, err := contactRepository.FindOne(&tc.Contact)
			if fmt.Sprint(err) != fmt.Sprint(tc.Err) {
				t.Errorf("got %v / want %v", err, tc.Err)
			}
			if c == nil {
				t.Errorf("got %v / want %v", c, tc.Contact)
			}
		})
	}
}
