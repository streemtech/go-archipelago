package managed

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/k0kubun/pp/v3"
	"github.com/pkg/errors"
	"github.com/streemtech/go-archipelago/api"
)

func (c *client) handleRoomInfo(ctx context.Context, cmd api.RoomInfo) (err error) {
	defer func() {
		if r := recover(); r != nil {
			pp.Println("panic in room info", r)
			debug.PrintStack()
		}
	}()

	// pp.Println(cmd, "Help")
	gamesToGetDPFor := []string{}
	for gameKey, gameHash := range cmd.DatapackageChecksums {
		// fmt.Println(gameKey, gameHash)
		// pp.Println(c)
		// pp.Println(c.dataPackages)
		game, ok := c.dataPackages[gameKey]
		// pp.Println(ok, game)
		if !ok || game.Checksum != gameHash {
			// fmt.Println("got game to get datapackage for")
			gamesToGetDPFor = append(gamesToGetDPFor, gameKey)
			//data package does not have game or datapackage checksum does not match the game checksum.
		} else {
			fmt.Println("got datapackage for game already")
		}
	}

	// pp.Println(gamesToGetDPFor)
	if len(gamesToGetDPFor) > 0 {
		// pp.Println("Sending data package")
		err = c.cmd.GetDataPackage(ctx, api.GetDataPackage{
			Games: &gamesToGetDPFor,
		})
		if err != nil {
			//TODO this may need a special case.
			c.wg.Done()
			// pp.Println("SendGetDatapackageErr ", err.Error())
			return errors.Wrap(err, "failed to request game data packages")
		}
	} else {
		c.wg.Done()
	}

	return nil
}

func (c *client) handleDataPackage(ctx context.Context, cmd api.DataPackage) (err error) {
	pp.Println(cmd)
	for game, data := range cmd.Data.Games {
		c.dataPackages[game] = data
	}
	c.wg.Done()
	return nil
}

func (c *client) handleConnected(ctx context.Context, cmd api.Connected) (err error) {
	c.connected = true
	c.connectedData = cmd
	pp.Println(c.connectedData)

	c.wg.Done()
	return nil
}
func (c *client) handleConnectionRefused(ctx context.Context, cmd api.ConnectionRefused) (err error) {
	c.connectionRefused = cmd
	fmt.Println("Get Handle Connection")
	c.wg.Done()
	return nil
}
func (c *client) handlePrintJson(ctx context.Context, cmd api.PrintJSON) (err error) {

	if c.pj != nil {
		var receiverSlotInfo api.NetworkSlot
		var finderSlotInfo api.NetworkSlot
		var item string
		var locationWhereFound string
		if cmd.Type != nil && *cmd.Type == api.PrintJsonTypeItemSend && cmd.Receiving != nil {
			receiverSlotInfo, _ = c.connectedData.SlotInfo[*cmd.Receiving]
		}

		if cmd.Type != nil && *cmd.Type == api.PrintJsonTypeItemSend && cmd.Item != nil {
			finderSlotInfo, _ = c.connectedData.SlotInfo[cmd.Item.Player]
		}

		if game, ok := c.dataPackages[receiverSlotInfo.Game]; ok {
			for k, v := range game.ItemNameToId {
				if v == cmd.Item.Item {
					item = k
					break
				}
			}
		}
		if game, ok := c.dataPackages[finderSlotInfo.Game]; ok {
			for k, v := range game.LocationNameToId {
				if v == cmd.Item.Location {
					locationWhereFound = k
					break
				}
			}
		}

		return c.pj(ctx, cmd, receiverSlotInfo, finderSlotInfo, item, locationWhereFound)
	}
	return nil
}

func (c *client) handleReceivedItems(ctx context.Context, cmd api.ReceivedItems) (err error) {
	if c.ri != nil {
		return c.ri(ctx, cmd)
	}
	return nil
}

func (c *client) handleRoomUpdate(ctx context.Context, cmd api.RoomUpdate) (err error) {
	if c.ru != nil {
		return c.ru(ctx, cmd)
	}
	return nil
}

func (c *client) handleBounced(ctx context.Context, cmd api.Bounced) (err error) {
	deathlink := false
	if cmd.Tags != nil {

		for _, v := range *cmd.Tags {
			if v == api.TagValueDeathLink {
				deathlink = true
			}
		}
	}
	if deathlink && c.dl != nil {
		return c.dl(ctx, cmd)
	}
	return nil
}
