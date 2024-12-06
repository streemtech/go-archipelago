package goarchipelago

import (
	"context"
	"fmt"
	"time"

	"github.com/coder/websocket"
)

// Test is a testing function to see if I can get archipelago BS to function.
func Test() {

	rootCtx := context.Background()
	conn, _, err := websocket.Dial(rootCtx, "wss://archipelago.gg:59452", &websocket.DialOptions{})

	if err != nil {
		panic(err.Error())
	}

	go func() {
		for {
			_, data, readErr := conn.Read(rootCtx)
			if readErr != nil {
				panic(readErr.Error())
			}
			fmt.Println(string(data))

			select {
			case <-rootCtx.Done():
				return
			default:
			}
		}
	}()

	time.Sleep(time.Second * 5)

}
