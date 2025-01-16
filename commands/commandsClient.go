package commands

import (
	"context"

	"github.com/streemtech/go-archipelago/api"
	"github.com/streemtech/go-archipelago/network"
	"github.com/streemtech/go-archipelago/utils"
)

type Client struct {
	Conn network.Client

	BouncedCommandHandler           func(ctx context.Context, cmd api.Bounced) (err error)
	ConnectedCommandHandler         func(ctx context.Context, cmd api.Connected) (err error)
	ConnectionRefusedCommandHandler func(ctx context.Context, cmd api.ConnectionRefused) (err error)
	DataPackageCommandHandler       func(ctx context.Context, cmd api.DataPackage) (err error)
	InvalidPacketCommandHandler     func(ctx context.Context, cmd api.InvalidPacket) (err error)
	LocationInfoCommandHandler      func(ctx context.Context, cmd api.LocationInfo) (err error)
	PrintJSONCommandHandler         func(ctx context.Context, cmd api.PrintJSON) (err error)
	ReceivedItemsCommandHandler     func(ctx context.Context, cmd api.ReceivedItems) (err error)
	RetrievedCommandHandler         func(ctx context.Context, cmd api.Retrieved) (err error)
	RoomInfoCommandHandler          func(ctx context.Context, cmd api.RoomInfo) (err error)
	RoomUpdateCommandHandler        func(ctx context.Context, cmd api.RoomUpdate) (err error)
	SetReplyCommandHandler          func(ctx context.Context, cmd api.SetReply) (err error)

	CloseCallbackHandler func(err error)
	CommandErrorHandler  func(cmd api.Command, err error)

	Log utils.Logger
}

func (c *Client) CommandCallback(ctx context.Context, cmds api.Commands) {
	for _, v := range cmds {
		err := c.handleCommand(ctx, v)
		if err != nil {
			if c.CommandErrorHandler != nil {
				c.CommandErrorHandler(v, err)
			}
		}
	}
}

func (c *Client) CloseCallback(err error) {
	if c.CloseCallbackHandler != nil {
		c.CloseCallbackHandler(err)
	}
}
