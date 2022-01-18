package metric

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetricInterface(t *testing.T) {

	chanelUUID := "e213g3ce-3fdb-4bdd-a9cf-3d4e9c0edf96"

	channelCreation := NewChannelCreation(chanelUUID)
	assert.NotNil(t, channelCreation)

	contactMessage := NewContactMessage(chanelUUID)
	assert.NotNil(t, contactMessage)

	contactActivation := NewContactActivation(chanelUUID)
	assert.NotNil(t, contactActivation)

	contactActivated := NewContactActivated(chanelUUID)
	assert.NotNil(t, contactActivated)

}
