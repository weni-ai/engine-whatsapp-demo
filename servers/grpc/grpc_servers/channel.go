package grpc_servers

import (
	"context"

	"github.com/teris-io/shortid"
	"github.com/weni/whatsapp-router/models"
	"github.com/weni/whatsapp-router/repositories"
	"github.com/weni/whatsapp-router/servers/grpc/pb"
	"go.mongodb.org/mongo-driver/mongo"
)

type ChannelServer struct {
	DB *mongo.Database
}

func NewChannelCServer(db *mongo.Database) *ChannelServer {
	return &ChannelServer{DB: db}
}

func (c *ChannelServer) CreateChannel(ctx context.Context, req *pb.ChannelRequest) (*pb.ChannelResponse, error) {
	var channel models.Channel
	channel.UUID = req.GetUuid()
	channel.Name = req.GetName()
	token, _ := shortid.Generate()

	channel.Token = token

	channelRepository := repositories.ChannelRepositoryDb{DB: c.DB}

	channelRepository.Insert(&channel)

	return &pb.ChannelResponse{
		Token: channel.Token,
	}, nil
}
