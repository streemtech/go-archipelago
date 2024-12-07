package main

import (
	"context"
	"fmt"

	// "github.com/google/uuid"

	"github.com/k0kubun/pp/v3"
	"github.com/pkg/errors"
	"github.com/streemtech/go-archipelago/api"
)

func (c *Client) handleCommand(ctx context.Context, cmd api.Command) error {
	cmdKey, err := cmd.Discriminator()

	if err != nil {
		return errors.New("Failed to extract cmd discriminator")
	}

	fmt.Printf("got %s command\n", cmdKey)

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
		return errors.Errorf("command %s should not be sent by server to client", cmdKey)
	default:
		return errors.Errorf("unknown discriminator cmd %s", cmdKey)
	}
	return nil
}

func (c *Client) handleBouncedCommand(ctx context.Context, cmd api.Bounced) (err error) {
	return nil
}
func (c *Client) handleConnectedCommand(ctx context.Context, cmd api.Connected) (err error) {
	return nil
}
func (c *Client) handleConnectionRefusedCommand(ctx context.Context, cmd api.ConnectionRefused) (err error) {
	return nil
}
func (c *Client) handleDataPackageCommand(ctx context.Context, cmd api.DataPackage) (err error) {
	for game, data := range cmd.Data.Games {
		c.dataPackages[game] = data
	}

	pp.Println(fmt.Sprintf("got %d items in data package\n", len(c.dataPackages)))

	return nil
}
func (c *Client) handleInvalidPacketCommand(ctx context.Context, cmd api.InvalidPacket) (err error) {
	return nil
}
func (c *Client) handleLocationInfoCommand(ctx context.Context, cmd api.LocationInfo) (err error) {
	return nil
}
func (c *Client) handlePrintJSONCommand(ctx context.Context, cmd api.PrintJSON) (err error) {
	return nil
}
func (c *Client) handleReceivedItemsCommand(ctx context.Context, cmd api.ReceivedItems) (err error) {
	return nil
}
func (c *Client) handleRetrievedCommand(ctx context.Context, cmd api.Retrieved) (err error) {
	return nil
}
func (c *Client) handleRoomInfoCommand(ctx context.Context, cmd api.RoomInfo) (err error) {
	gamesToGetDPFor := []string{}
	for gameKey, gameHash := range cmd.DatapackageChecksums {
		if game, ok := c.dataPackages[gameKey]; !ok || game.Checksum != gameHash {
			gamesToGetDPFor = append(gamesToGetDPFor, gameKey)
			//data package does not have game or datapackage checksum does not match the game checksum.
		}
	}

	if len(gamesToGetDPFor) > 0 {
		err = c.GetDataPackage(ctx, api.GetDataPackage{
			Games: &gamesToGetDPFor,
		})
		if err != nil {
			return errors.Wrap(err, "failed to request game data packages")
		}
	}

	return nil
}
func (c *Client) handleRoomUpdateCommand(ctx context.Context, cmd api.RoomUpdate) (err error) {
	return nil
}
func (c *Client) handleSetReplyCommand(ctx context.Context, cmd api.SetReply) (err error) {
	return nil
}
