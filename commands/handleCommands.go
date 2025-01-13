package commands

import (
	"context"
	"fmt"

	// "github.com/google/uuid"

	"github.com/pkg/errors"
	"github.com/streemtech/go-archipelago/api"
)

func (c *Client) handleCommand(ctx context.Context, cmd api.Command) error {
	cmdKey, err := cmd.Discriminator()

	if err != nil {
		return errors.New("Failed to extract cmd discriminator")
	}

	switch api.CommandKeys(cmdKey) {
	case api.CommandKeyBounced:
		d, err := cmd.AsBounced()
		if err != nil {
			return errors.Wrap(err, "failed to get data matching command key")
		}
		err = c.handleBouncedCommand(ctx, d)
		if err != nil {
			return errors.Wrap(err, "failed to handle Bounced Command")
		}

	case api.CommandKeyConnected:
		d, err := cmd.AsConnected()
		if err != nil {
			return errors.Wrap(err, "failed to get data matching command key")
		}
		err = c.handleConnectedCommand(ctx, d)
		if err != nil {
			return errors.Wrap(err, "failed to handle Connected Command")
		}

	case api.CommandKeyConnectionRefused:
		d, err := cmd.AsConnectionRefused()
		if err != nil {
			return errors.Wrap(err, "failed to get data matching command key")
		}
		err = c.handleConnectionRefusedCommand(ctx, d)
		if err != nil {
			return errors.Wrap(err, "failed to handle ConnectionRefused Command")
		}

	case api.CommandKeyDataPackage:
		d, err := cmd.AsDataPackage()
		if err != nil {
			return errors.Wrap(err, "failed to get data matching command key")
		}
		err = c.handleDataPackageCommand(ctx, d)
		if err != nil {
			return errors.Wrap(err, "failed to handle DataPackage Command")
		}

	case api.CommandKeyInvalidPacket:
		d, err := cmd.AsInvalidPacket()
		if err != nil {
			return errors.Wrap(err, "failed to get data matching command key")
		}
		err = c.handleInvalidPacketCommand(ctx, d)
		if err != nil {
			return errors.Wrap(err, "failed to handle InvalidPacket Command")
		}

	case api.CommandKeyLocationInfo:
		d, err := cmd.AsLocationInfo()
		if err != nil {
			return errors.Wrap(err, "failed to get data matching command key")
		}
		err = c.handleLocationInfoCommand(ctx, d)
		if err != nil {
			return errors.Wrap(err, "failed to handle LocationInfo Command")
		}

	case api.CommandKeyPrintJSON:
		d, err := cmd.AsPrintJSON()
		if err != nil {
			return errors.Wrap(err, "failed to get data matching command key")
		}
		err = c.handlePrintJSONCommand(ctx, d)
		if err != nil {
			return errors.Wrap(err, "failed to handle PrintJSON Command")
		}

	case api.CommandKeyReceivedItems:
		d, err := cmd.AsReceivedItems()
		if err != nil {
			return errors.Wrap(err, "failed to get data matching command key")
		}
		err = c.handleReceivedItemsCommand(ctx, d)
		if err != nil {
			return errors.Wrap(err, "failed to handle ReceivedItems Command")
		}

	case api.CommandKeyRetrieved:
		d, err := cmd.AsRetrieved()
		if err != nil {
			return errors.Wrap(err, "failed to get data matching command key")
		}
		err = c.handleRetrievedCommand(ctx, d)
		if err != nil {
			return errors.Wrap(err, "failed to handle Retrieved Command")
		}

	case api.CommandKeyRoomInfo:
		d, err := cmd.AsRoomInfo()
		if err != nil {
			return errors.Wrap(err, "failed to get data matching command key")
		}
		err = c.handleRoomInfoCommand(ctx, d)
		if err != nil {
			return errors.Wrap(err, "failed to handle RoomInfo Command")
		}

	case api.CommandKeyRoomUpdate:
		d, err := cmd.AsRoomUpdate()
		if err != nil {
			return errors.Wrap(err, "failed to get data matching command key")
		}
		err = c.handleRoomUpdateCommand(ctx, d)
		if err != nil {
			return errors.Wrap(err, "failed to handle RoomUpdate Command")
		}

	case api.CommandKeySetReply:
		d, err := cmd.AsSetReply()
		if err != nil {
			return errors.Wrap(err, "failed to get data matching command key")
		}
		err = c.handleSetReplyCommand(ctx, d)
		if err != nil {
			return errors.Wrap(err, "failed to handle SetReply Command")
		}

	case api.CommandKeyBounce,
		api.CommandKeyConnect,
		api.CommandKeyConnectUpdate,
		api.CommandKeyGet,
		api.CommandKeyGetDataPackage,
		api.CommandKeyLocationChecks,
		api.CommandKeyLocationScouts,
		api.CommandKeySay,
		api.CommandKeySet,
		api.CommandKeySetNotify,
		api.CommandKeyStatusUpdate,
		api.CommandKeySync,
		api.CommandKeyUpdateHint:
		return errors.Errorf("command %s should not be sent to client by the server", cmdKey)
	default:
		return errors.Errorf("unknown discriminator cmd %s", cmdKey)
	}
	return nil
}

func (c *Client) handleBouncedCommand(ctx context.Context, cmd api.Bounced) (err error) {
	if c.BouncedCommandHandler != nil {
		return c.BouncedCommandHandler(ctx, cmd)
	}
	fmt.Println("Got Unhandled Bounced Command")
	return nil
}
func (c *Client) handleConnectedCommand(ctx context.Context, cmd api.Connected) (err error) {
	if c.ConnectedCommandHandler != nil {
		return c.ConnectedCommandHandler(ctx, cmd)
	}
	fmt.Println("Got Unhandled Connected Command")
	return nil
}
func (c *Client) handleConnectionRefusedCommand(ctx context.Context, cmd api.ConnectionRefused) (err error) {
	if c.ConnectionRefusedCommandHandler != nil {
		return c.ConnectionRefusedCommandHandler(ctx, cmd)
	}
	fmt.Println("Got Unhandled ConnectionRefused Command")
	return nil
}
func (c *Client) handleDataPackageCommand(ctx context.Context, cmd api.DataPackage) (err error) {
	if c.DataPackageCommandHandler != nil {
		return c.DataPackageCommandHandler(ctx, cmd)
	}
	fmt.Println("Got Unhandled DataPackage Command")

	return nil
}
func (c *Client) handleInvalidPacketCommand(ctx context.Context, cmd api.InvalidPacket) (err error) {
	if c.InvalidPacketCommandHandler != nil {
		return c.InvalidPacketCommandHandler(ctx, cmd)
	}
	fmt.Println("Got Unhandled InvalidPacket Command")
	return nil
}
func (c *Client) handleLocationInfoCommand(ctx context.Context, cmd api.LocationInfo) (err error) {
	if c.LocationInfoCommandHandler != nil {
		return c.LocationInfoCommandHandler(ctx, cmd)
	}
	fmt.Println("Got Unhandled LocationInfo Command")
	return nil
}
func (c *Client) handlePrintJSONCommand(ctx context.Context, cmd api.PrintJSON) (err error) {
	if c.PrintJSONCommandHandler != nil {
		return c.PrintJSONCommandHandler(ctx, cmd)
	}
	fmt.Println("Got Unhandled PrintJSON Command")
	return nil
}
func (c *Client) handleReceivedItemsCommand(ctx context.Context, cmd api.ReceivedItems) (err error) {
	if c.ReceivedItemsCommandHandler != nil {
		return c.ReceivedItemsCommandHandler(ctx, cmd)
	}
	fmt.Println("Got Unhandled ReceivedItems Command")
	return nil
}
func (c *Client) handleRetrievedCommand(ctx context.Context, cmd api.Retrieved) (err error) {
	if c.RetrievedCommandHandler != nil {
		return c.RetrievedCommandHandler(ctx, cmd)
	}
	fmt.Println("Got Unhandled Retrieved Command")
	return nil
}
func (c *Client) handleRoomInfoCommand(ctx context.Context, cmd api.RoomInfo) (err error) {
	if c.RoomInfoCommandHandler != nil {
		return c.RoomInfoCommandHandler(ctx, cmd)
	}
	fmt.Println("Got Unhandled RoomInfo Command")
	return nil
}
func (c *Client) handleRoomUpdateCommand(ctx context.Context, cmd api.RoomUpdate) (err error) {
	if c.RoomUpdateCommandHandler != nil {
		return c.RoomUpdateCommandHandler(ctx, cmd)
	}
	fmt.Println("Got Unhandled RoomUpdate Command")
	return nil
}
func (c *Client) handleSetReplyCommand(ctx context.Context, cmd api.SetReply) (err error) {
	if c.SetReplyCommandHandler != nil {
		return c.SetReplyCommandHandler(ctx, cmd)
	}
	fmt.Println("Got Unhandled SetReply Command")
	return nil
}
