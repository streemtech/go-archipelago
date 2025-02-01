package managed

import (
	"context"
	"log/slog"
	"sync"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/streemtech/go-archipelago/api"
	"github.com/streemtech/go-archipelago/commands"
	"github.com/streemtech/go-archipelago/network"
	"github.com/streemtech/go-archipelago/utils"
)

func WithPassword(password string) Option {
	return func(client *client) {
		client.password = password
	}
}

func WithGame(game string) Option {
	return func(client *client) {
		client.game = game
	}
}

// triggers when a deathlink is gotten, and does something.
func WithLogger(l utils.Logger) Option {
	return func(client *client) {
		client.log = l
	}
}

// This will likely be changed to with global item recieved.
func WithOnPrintJSON(callback func(ctx context.Context, cmd api.PrintJSON, receiverSlotInfo api.NetworkSlot, finderSlotInfo api.NetworkSlot, item string, location_where_found string) error) Option {
	return func(client *client) {
		client.pj = callback
	}
}

// sets the check to be told that you have recieved an item.
func WithOnReceivedItems(callback func(ctx context.Context, cmd ReceivedItems) error) Option {
	return func(client *client) {
		client.ri = callback
	}
}

// sets the check to be told that you have collected an item somewhere. (AKA that the local state of collectables needs to be updated.)
func WithOnRoomUpdate(callback func(ctx context.Context, cmd api.RoomUpdate) error) Option {
	return func(client *client) {
		client.ru = callback
	}
}

// triggers when a deathlink is gotten, and does something.
func WithOnDeathlink(callback func(ctx context.Context, cmd DeathLink) error) Option {
	return func(client *client) {
		client.dl = callback
	}
}

// triggers when the disconnect is called
func WithOnDisconnect(callback func(err error)) Option {
	return func(client *client) {
		client.dc = callback
	}
}

type Option func(client *client)

type client struct {
	cmd *commands.Client
	net network.Client

	//contains game name to game data
	dataPackages map[string]api.GameData

	wg *sync.WaitGroup

	connected         bool
	connectionRefused api.ConnectionRefused
	connectedData     api.Connected

	game     string
	password string
	slot     string

	log utils.Logger

	//print json command callback. Will eventually be replaced with dedicated, cleaner "Get Item" callback
	pj func(ctx context.Context, cmd api.PrintJSON, receiverSlotInfo api.NetworkSlot, finderSlotInfo api.NetworkSlot, item string, location_where_found string) error
	//received items handler
	ri func(ctx context.Context, cmd ReceivedItems) error
	//room update handler
	ru func(ctx context.Context, cmd api.RoomUpdate) error
	//deathlink bounce command handler (autofilters deathlinks)
	dl func(ctx context.Context, cmd DeathLink) error
	//disconnect handler
	dc func(err error)
}

func NewClient(address string, slotName string, opts ...Option) (c *client, err error) {

	c = &client{
		wg:           &sync.WaitGroup{},
		dataPackages: map[string]api.GameData{},
		slot:         slotName,
		log:          slog.Default(),
	}

	//commands must be set here, and not in the generator itself to prevent a nul pointer as c is null before its created.
	c.cmd = &commands.Client{}
	for _, v := range opts {
		v(c)
	}
	c.cmd.RoomInfoCommandHandler = c.handleRoomInfo
	c.cmd.DataPackageCommandHandler = c.handleDataPackage
	c.cmd.ConnectedCommandHandler = c.handleConnected
	c.cmd.ConnectionRefusedCommandHandler = c.handleConnectionRefused
	c.cmd.PrintJSONCommandHandler = c.handlePrintJson
	c.cmd.ReceivedItemsCommandHandler = c.handleReceivedItems
	c.cmd.RoomUpdateCommandHandler = c.handleRoomUpdate
	c.cmd.CloseCallbackHandler = c.handleClose

	c.cmd.Log = c.log

	c.net, err = network.NewClient(network.ClientProps{
		Url:              address,
		CommandsCallback: c.cmd.CommandCallback,
		CloseCallback:    c.cmd.CloseCallback,
		Log:              c.log,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect create network client")
	}
	c.cmd.Conn = c.net

	if c.log != nil {
		c.log.Debug("starting network client")
	}

	//start the command client, which initiates the method to listen for data packages etc. and wait for the WAIT command.
	c.wg.Add(1)
	err = c.cmd.Conn.Start()
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize client")
	}

	if c.log != nil {
		c.log.Debug("waiting for ")
	}
	c.wg.Wait()

	// fmt.Println("Starting Connect")
	//send the connect packet and
	c.wg.Add(1)
	err = c.cmd.Connect(context.Background(), api.Connect{
		Game:          c.game,
		ItemsHandling: api.ItemHandlingFlagOtherWorlds | api.ItemHandlingFlagOwnWorld | api.ItemHandlingFlagStartingInventory,
		Name:          c.slot,
		Password:      c.password,
		SlotData:      true,
		Uuid:          uuid.NewString(),
		Version: api.NetworkVersion{
			Major: 5,
			Minor: 0,
			Build: 1,
			Class: "Version",
		},
		Tags: []api.Tags{
			api.TagValueDeathLink,
			api.TagValueTextOnly,
		},
	})
	if err != nil {
		c.net.Close()
		return nil, errors.Wrap(err, "failed to send connect")
	}
	c.wg.Wait()
	if false {
		c.net.Close()
		return nil, errors.Errorf("failed to connect to server: %v", c.connectionRefused.Errors)
	}

	return c, nil
}
