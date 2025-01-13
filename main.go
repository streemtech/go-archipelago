package main

import (

	// "github.com/google/uuid"

	"context"
	"fmt"
	"time"

	"github.com/k0kubun/pp/v3"
	"github.com/pkg/errors"
	"github.com/streemtech/go-archipelago/api"
	"github.com/streemtech/go-archipelago/managed"
)

func ptr[T any](data T) *T {
	return &data
}

// Testing only. Will be removed shortly.
func main() {

	now := time.Now()
	mc, err := managed.NewClient("wss://archipelago.gg:12345", "example", managed.WithOnPrintJSON(func(ctx context.Context, cmd api.PrintJSON, receiverSlotInfo api.NetworkSlot, finderSlotInfo api.NetworkSlot, item string, location_where_found string) error {
		pp.Println(cmd, receiverSlotInfo, finderSlotInfo, item, location_where_found)
		return nil
	}))
	if err != nil {
		panic(err.Error())
	}

	// time.Sleep(time.Second)
	mc.Deathlink(context.Background(), "Bot Builder triggered a deathlink")
	// time.Sleep(time.Hour * 24)

	err = mc.Close()
	if err != nil {
		panic(errors.Wrap(err, "failed to close").Error())
	}
	fmt.Println(time.Since(now))

}
