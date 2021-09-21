package repositories

import (
	"fmt"
	"testing"

	"github.com/weni/whatsapp-router/models"
	"github.com/weni/whatsapp-router/storage"
)

var tcInsertChannel = []struct {
	TestName string
	Channel  models.Channel
	Err      error
}{
	{
		TestName: "Insert Channel",
		Channel: models.Channel{
			UUID:  "asdfgh",
			Name:  "foo",
			Token: "foo-bar-zaz",
		},
		Err: nil,
	},
}

func TestInsertChannel(t *testing.T) {
	mongodb := storage.NewDB()
	channelRepository := ChannelRepositoryDb{DB: mongodb}

	for _, tc := range tcInsertChannel {
		t.Run(tc.TestName, func(t *testing.T) {
			err := channelRepository.Insert(&tc.Channel)
			fmt.Println(tc.Channel.ID.Hex())
			if fmt.Sprint(err) != fmt.Sprint(tc.Err) {
				t.Errorf("got %v / want %v", err, tc.Err)
			}
		})
	}
}

var tcFindOneChannel = []struct {
	TestName string
	Channel  models.Channel
	Err      error
}{
	{
		TestName: "Find one existing channel",
		Channel: models.Channel{
			UUID:  "491bab89-5998-4cf4-896b-afab3c5e1090",
			Name:  "foo",
			Token: "foo-bar-zaz",
		},
	},
}

func TestFindOneChannel(t *testing.T) {
	mongodb := storage.NewDB()
	channelRepository := ChannelRepositoryDb{DB: mongodb}

	for _, tc := range tcFindOneChannel {
		t.Run(tc.TestName, func(t *testing.T) {
			c, err := channelRepository.FindOne(&tc.Channel)
			if fmt.Sprint(err) != fmt.Sprint(tc.Err) {
				t.Errorf("got %v / want %v", err, tc.Err)
			}
			if c == nil {
				t.Errorf("got %v / want %v", c, tc.Channel)
			}
		})
	}
}
