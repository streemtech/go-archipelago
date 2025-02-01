package managed

import (
	"time"

	"github.com/streemtech/go-archipelago/api"
)

//Types.go contains the types used for certain callbacks so that the right data is returned.

type ReceivedItems struct {
	Raw            api.ReceivedItems
	ItemName       string
	SourceSlotName string
	DestSlotName   string
	LocationName   string
}

type DeathLink struct {
	Raw    api.Bounced
	Time   time.Time
	Cause  string
	Source string
}
