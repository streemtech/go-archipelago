package managed

import (
	"context"
	"time"

	"github.com/streemtech/go-archipelago/api"
)

func (c *client) Close() (err error) {
	if c.log != nil {
		c.log.Debug("closing connection")
	}
	return c.net.Close()
}

// cause is the cause of the dathlink. For noderunner this will be an argument or default to "foo"
func (c *client) Deathlink(ctx context.Context, cause string) (err error) {
	if c.log != nil {
		c.log.Info("sending deathlink")
	}
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

func (c *client) Say(ctx context.Context, message string) (err error) {
	if c.log != nil {
		c.log.Info("sending message")
	}
	return c.cmd.Say(ctx, api.Say{
		Text: message,
	})

}

func (c *client) Collect(ctx context.Context, location int) (err error) {
	if c.log != nil {
		c.log.Info("collecting location")
	}
	return c.cmd.LocationChecks(ctx, api.LocationChecks{
		Locations: []int{location},
	})
}
