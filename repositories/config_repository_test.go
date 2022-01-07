package repositories

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/weni/whatsapp-router/models"
	"github.com/weni/whatsapp-router/storage"
)

func TestConfigRepository(t *testing.T) {
	db := storage.NewTestDB()
	defer storage.CloseDB(db)
	storage.CleanupDB(db)
	repo := NewConfigRepository(db)

	conf1 := models.Config{
		Token: "qwert1234",
	}
	// test Create
	err := repo.Create(&conf1)
	assert.Nil(t, err)

	conf2, err := repo.GetFirst()
	assert.Nil(t, err)
	assert.Equal(t, conf1.Token, conf2.Token)

	// test FindOne
	conf3, err := repo.FindOne(conf2)
	assert.Nil(t, err)
	assert.Equal(t, conf2, conf3)

	newToken := "asdfg12345"
	conf3.Token = newToken
	// test Update
	conf4, err := repo.Update(conf3)
	assert.Nil(t, err)
	assert.Equal(t, conf3, conf4)
}
