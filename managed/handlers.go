package managed

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	// "github.com/k0kubun/pp/v3"
	"github.com/pkg/errors"
	"github.com/streemtech/go-archipelago/api"
)

func (c *client) handleRoomInfo(ctx context.Context, cmd api.RoomInfo) (err error) {
	if c.log != nil {
		c.log.Debug("received room info")
	}
	defer func() {
		if r := recover(); r != nil {
			// pp.Println("panic in room info", r)
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
		if c.log != nil {
			c.log.Debug("requested data package")
		}
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
		//dont need to wait for data package, can just end it here.
		c.wg.Done()
	}

	return nil
}

func (c *client) handleDataPackage(ctx context.Context, cmd api.DataPackage) (err error) {
	// pp.Println(cmd)
	if c.log != nil {
		c.log.Debug("received data package")
	}
	for game, data := range cmd.Data.Games {
		c.dataPackages[game] = data
	}
	c.wg.Done()
	return nil
}

func (c *client) handleConnected(ctx context.Context, cmd api.Connected) (err error) {
	c.connected = true
	c.connectedData = cmd
	// pp.Println(c.connectedData)
	if c.log != nil {
		c.log.Debug("connected to server")
	}
	c.wg.Done()
	return nil
}
func (c *client) handleConnectionRefused(ctx context.Context, cmd api.ConnectionRefused) (err error) {
	c.connectionRefused = cmd
	// fmt.Println("Get Handle Connection")
	if c.log != nil {
		c.log.Error("connection to server refused")
	}

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
		//TODO0 set the received items data.

		var slot api.NetworkSlot
		for _, v := range c.connectedData.SlotInfo {
			if v.Name == c.slot {
				slot = v
				break
			}
		}

		for _, v := range cmd.Items {
			sourceSlotInfo := c.connectedData.SlotInfo[v.Player]

			destGame := c.dataPackages[slot.Game]

			sourceGameData := c.dataPackages[sourceSlotInfo.Game]

			//set item name
			itemName := ""
			for name, id := range destGame.ItemNameToId {
				if id == v.Item {
					itemName = name
					break
				}
			}

			//set location name
			locationName := ""
			for name, id := range sourceGameData.LocationNameToId {
				if id == v.Location {
					locationName = name
					break
				}
			}

			rec := ReceivedItems{
				Raw:            cmd,
				SourceSlotName: sourceSlotInfo.Name,
				DestSlotName:   slot.Name,
				ItemName:       itemName,
				LocationName:   locationName,
			}
			err = c.ri(ctx, rec)
			if err != nil {
				return errors.Wrap(err, "failed to send received to handler")
			}
		}

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
	if c.log != nil {
		c.log.Debug("received bounce command")
	}
	deathlink := false
	if cmd.Tags != nil {

		for _, v := range *cmd.Tags {
			if v == api.TagValueDeathLink {
				deathlink = true
			}

		}
	}
	if deathlink && c.dl != nil {
		dl := DeathLink{
			Raw: cmd,
		}

		cause, _ := cmd.Data["cause"]
		source, _ := cmd.Data["source"]
		timeField, _ := cmd.Data["time"]

		dl.Cause, _ = cause.(string)
		dl.Source, _ = source.(string)
		timeFloat, _ := timeField.(float64)

		remainder := int64(timeFloat*1_000_000_000.0) % 1_000_000_000

		time.Unix(int64(timeFloat), remainder)

		return c.dl(ctx, dl)
	}
	return nil
}

func (c *client) handleClose(err error) {
	if c.log != nil {
		c.log.Error("Got close")
	}
	if c.dc == nil {
		return
	}
	c.dc(err)

}
