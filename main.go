package main

import (
	"context"
	"fmt"
	"time"

	"github.com/k0kubun/pp/v3"
	"github.com/streemtech/go-archipelago/api"
	"github.com/streemtech/go-archipelago/network"
)

// Testing only. Will be removed shortly.
func main() {
	network.NewClient(network.ClientProps{
		Url: "wss://archipelago.gg:59024",
		CommandsCallback: func(ctx context.Context, cmds api.Commands) {
			ri, err := cmds[0].AsRoomInfo()
			pp.Println(ri)
			if err != nil {
				fmt.Println(err.Error())
			}
		},
		CloseCallback: func(err error) {
			pp.Println("Service Closed", err.Error())
		},
	})
	time.Sleep(time.Second * 5)
}
