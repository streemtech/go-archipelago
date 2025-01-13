package managed

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/streemtech/go-archipelago/api"
	"github.com/streemtech/go-archipelago/commands"
	"github.com/streemtech/go-archipelago/network"
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

// This will likely be changed to with global item recieved.
func WithOnPrintJSON(callback func(ctx context.Context, cmd api.PrintJSON, receiverSlotInfo api.NetworkSlot, finderSlotInfo api.NetworkSlot, item string, location_where_found string) error) Option {
	return func(client *client) {
		client.pj = callback
	}
}

// sets the check to be told that you have recieved an item.
func WithOnReceivedItems(callback func(ctx context.Context, cmd api.ReceivedItems) error) Option {
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
func WithOnDeathlink(callback func(ctx context.Context, cmd api.Bounced) error) Option {
	return func(client *client) {
		client.dl = callback
	}
}

type Option func(client *client)

type client struct {
	cmd *commands.Client
	net network.Client

	dataPackages map[string]api.GameData

	wg *sync.WaitGroup

	connected         bool
	connectionRefused api.ConnectionRefused
	connectedData     api.Connected

	game     string
	password string
	slot     string

	//print json command callback. Will eventually be replaced with dedicated, cleaner "Get Item" callback
	pj func(ctx context.Context, cmd api.PrintJSON, receiverSlotInfo api.NetworkSlot, finderSlotInfo api.NetworkSlot, item string, location_where_found string) error
	ri func(ctx context.Context, cmd api.ReceivedItems) error
	ru func(ctx context.Context, cmd api.RoomUpdate) error
	dl func(ctx context.Context, cmd api.Bounced) error
}

func NewClient(address string, slotName string, opts ...Option) (c *client, err error) {

	c = &client{
		wg:           &sync.WaitGroup{},
		dataPackages: map[string]api.GameData{},
		slot:         slotName,
	}

	//commands must be set here, and not in the generator itself to prevent a nul pointer as c is null before its created.
	c.cmd = &commands.Client{
		RoomInfoCommandHandler:          c.handleRoomInfo,
		DataPackageCommandHandler:       c.handleDataPackage,
		ConnectedCommandHandler:         c.handleConnected,
		ConnectionRefusedCommandHandler: c.handleConnectionRefused,
		PrintJSONCommandHandler:         c.handlePrintJson,
		ReceivedItemsCommandHandler:     c.handleReceivedItems,
		RoomUpdateCommandHandler:        c.handleRoomUpdate,
	}

	for _, v := range opts {
		v(c)
	}

	c.net, err = network.NewClient(network.ClientProps{
		Url:              address,
		CommandsCallback: c.cmd.CommandCallback,
		CloseCallback:    c.cmd.CloseCallback,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect create network client")
	}
	c.cmd.Conn = c.net

	// fmt.Println("Starting Network client")

	//start the command client, which initiates the method to listen for data packages etc. and wait for the WAIT command.
	c.wg.Add(1)
	err = c.cmd.Conn.Start()
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize client")
	}

	// fmt.Println("Waiting Network client")
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

func (c *client) Close() (err error) {
	return c.net.Close()
}

// cause is the cause of the dathlink. For noderunner this will be an argument or default to "foo"
func (c *client) Deathlink(ctx context.Context, cause string) (err error) {
	return c.cmd.Bounce(ctx, api.Bounce{
		Tags: &[]api.Tags{
			api.TagValueDeathLink,
		},
		Data: map[string]interface{}{
			"time":   float64(time.Now().Unix()) + 0.0000001,
			"cause":  cause,
			"source": c.slot,
		},
	})
}
