package repositories

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/weni/whatsapp-router/models"
	"github.com/weni/whatsapp-router/storage"
)

var tcInsertFlows = []struct {
	TestName string
	Flows    models.Flows
	Err      error
}{
	{
		TestName: "test insert flows",
		Flows: models.Flows{
			FlowsStarts: []models.Flow{
				{
					Name:    "test_flow1",
					UUID:    "507b6703-cc80-41fc-8a1b-cca573518dbb",
					Keyword: "hello1",
				},
				{
					Name:    "test_flow2",
					UUID:    "a76b3106-5e3d-462d-a0fc-4817c0d73ce7",
					Keyword: "hello2",
				},
				{
					Name:    "test_flow3",
					UUID:    "d7c97de5-bd06-4d7f-904f-63a7f8dd6b9d",
					Keyword: "hello3",
				},
			},
			Channel: dummyChannel.UUID,
		},
	},
}

func TestInsertFlows(t *testing.T) {
	mongodb := storage.NewTestDB()
	defer storage.CloseDB(mongodb)
	flowsRepository := FlowsRepositoryDb{DB: mongodb}

	for _, tc := range tcInsertFlows {
		t.Run(tc.TestName, func(t *testing.T) {
			f, err := flowsRepository.Insert(&tc.Flows)
			if fmt.Sprint(err) != fmt.Sprint(tc.Err) {
				t.Errorf("got %v / want %v", err, tc.Err)
			}

			if f == nil {
				t.Errorf("got %v / want %v", f, reflect.TypeOf(tc.Flows))
			}
		})
	}
}

var tcFindOneFlows = []struct {
	TestName string
	Flows    models.Flows
	Err      error
}{
	{
		TestName: "Find one existing flows",
		Flows: models.Flows{
			FlowsStarts: []models.Flow{
				{
					Name:    "test_flow1",
					UUID:    "507b6703-cc80-41fc-8a1b-cca573518dbb",
					Keyword: "hello1",
				},
				{
					Name:    "test_flow2",
					UUID:    "a76b3106-5e3d-462d-a0fc-4817c0d73ce7",
					Keyword: "hello2",
				},
				{
					Name:    "test_flow3",
					UUID:    "d7c97de5-bd06-4d7f-904f-63a7f8dd6b9d",
					Keyword: "hello3",
				},
			},
			Channel: dummyChannel.UUID,
		},
	},
}

func TestFindOneFlows(t *testing.T) {
	mongodb := storage.NewTestDB()
	defer storage.CloseDB(mongodb)
	flowsRepository := FlowsRepositoryDb{DB: mongodb}

	for _, tc := range tcFindOneFlows {
		t.Run(tc.TestName, func(t *testing.T) {
			f, err := flowsRepository.FindOne(&tc.Flows)
			if fmt.Sprint(err) != fmt.Sprint(tc.Err) {
				t.Errorf("got %v / want %v", err, tc.Err)
			}
			if f == nil {
				t.Errorf("got %v / want %v", f, tc.Flows)
			}
		})
	}
}

var tcUpdateFlows = []struct {
	TestName       string
	Flows          models.Flows
	Err            error
	ChannelToUdate models.Channel
}{
	{
		TestName: "Find one existing flows",
		Flows: models.Flows{
			FlowsStarts: []models.Flow{
				{
					Name:    "test_flow1",
					UUID:    "507b6703-cc80-41fc-8a1b-cca573518dbb",
					Keyword: "hello1",
				},
				{
					Name:    "test_flow2",
					UUID:    "a76b3106-5e3d-462d-a0fc-4817c0d73ce7",
					Keyword: "hello2",
				},
				{
					Name:    "test_flow3",
					UUID:    "d7c97de5-bd06-4d7f-904f-63a7f8dd6b9d",
					Keyword: "hello3",
				},
			},
			Channel: dummyChannel.UUID,
		},
		ChannelToUdate: dummyChannel2,
	},
}

func TestUpdateFlows(t *testing.T) {
	mongodb := storage.NewTestDB()
	defer storage.CloseDB(mongodb)
	flowsRepository := FlowsRepositoryDb{DB: mongodb}

	for _, tc := range tcUpdateFlows {
		t.Run(tc.TestName, func(t *testing.T) {
			flowsToUpdate := &models.Flows{
				FlowsStarts: tc.Flows.FlowsStarts,
				Channel:     tc.ChannelToUdate.UUID,
			}
			f, err := flowsRepository.Update(flowsToUpdate)
			if fmt.Sprint(err) != fmt.Sprint(tc.Err) {
				t.Errorf("got %v / want %v", err, tc.Err)
			}
			if f == nil {
				t.Errorf("got %v / want %v", f, flowsToUpdate)
			}
		})
	}
}
