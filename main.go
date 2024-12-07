package main

import (
	"context"
	"sync"
	"time"

	// "github.com/google/uuid"

	"github.com/k0kubun/pp/v3"
	"github.com/pkg/errors"
	"github.com/streemtech/go-archipelago/api"
	"github.com/streemtech/go-archipelago/network"
)

type Client struct {
	dataPackages map[string]api.GameData

	conn network.Client
	wg   sync.WaitGroup
}

func ptr[T any](data T) *T {
	return &data
}

// Testing only. Will be removed shortly.
func main() {

	client := &Client{
		dataPackages: map[string]api.GameData{},
	}

	client.wg.Add(1)

	networkClient, err := network.NewClient(network.ClientProps{
		Url:              "wss://archipelago.gg:59024",
		CommandsCallback: client.commandCallback,
		CloseCallback:    client.closeCallback,
	})

	if err != nil {
		panic(err.Error())
	}
	client.conn = networkClient

	client.conn.Start()

	time.Sleep(time.Second * 5)

	err = errors.Wrap(client.Connect(context.Background(), api.Connect{
		Game:          "Clique",
		ItemsHandling: api.ItemHandlingFlagOtherWorlds,
		Name:          "JeffClique1",
		Password:      "",
		SlotData:      true,
		// Uuid:          uuid.NewString(),
		Version: api.NetworkVersion{
			Major: ptr(5),
			Minor: ptr(0),
			Build: ptr(1),
		},
		Tags: []api.Tags{
			// api.TagValueHintGame,
		},
	}), "failed to send connect")

	if err != nil {
		pp.Println(err.Error())
	}

	client.wg.Wait()

}

func (c *Client) commandCallback(ctx context.Context, cmds api.Commands) {
	for _, v := range cmds {
		err := c.handleCommand(ctx, v)
		if err != nil {
			pp.Println("command error", err.Error())
		}
	}
}

func (c *Client) closeCallback(err error) {
	pp.Println("Service Closed", err.Error())
	c.wg.Done()
}
