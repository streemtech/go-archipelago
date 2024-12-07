package main

import (
	"context"
	"fmt"

	// "github.com/google/uuid"

	"github.com/pkg/errors"
	"github.com/streemtech/go-archipelago/api"
)

func (c *Client) Bounce(ctx context.Context, data api.Bounce) (err error) {
	command := api.Command{}
	command.MergeBounce(data)
	fmt.Println("Sending Bounce command")
	err = c.conn.SendCommand([]api.Command{
		command,
	})
	if err != nil {
		return errors.Wrap(err, "failed to send Bounce command")
	}
	return nil
}
func (c *Client) Connect(ctx context.Context, data api.Connect) (err error) {
	command := api.Command{}
	command.MergeConnect(data)
	fmt.Println("Sending Connect command")
	err = c.conn.SendCommand([]api.Command{
		command,
	})
	if err != nil {
		return errors.Wrap(err, "failed to send Connect command")
	}
	return nil
}
func (c *Client) ConnectUpdate(ctx context.Context, data api.ConnectUpdate) (err error) {
	command := api.Command{}
	command.MergeConnectUpdate(data)
	fmt.Println("Sending ConnectUpdate command")
	err = c.conn.SendCommand([]api.Command{
		command,
	})
	if err != nil {
		return errors.Wrap(err, "failed to send ConnectUpdate command")
	}
	return nil
}
func (c *Client) Get(ctx context.Context, data api.Get) (err error) {
	command := api.Command{}
	command.MergeGet(data)
	fmt.Println("Sending Get command")
	err = c.conn.SendCommand([]api.Command{
		command,
	})
	if err != nil {
		return errors.Wrap(err, "failed to send Get command")
	}
	return nil
}
func (c *Client) GetDataPackage(ctx context.Context, data api.GetDataPackage) (err error) {
	command := api.Command{}
	command.MergeGetDataPackage(data)
	fmt.Println("Sending GetDataPackage command")
	err = c.conn.SendCommand([]api.Command{
		command,
	})
	if err != nil {
		return errors.Wrap(err, "failed to send GetDataPackage command")
	}
	return nil
}
func (c *Client) LocationChecks(ctx context.Context, data api.LocationChecks) (err error) {
	command := api.Command{}
	command.MergeLocationChecks(data)
	fmt.Println("Sending LocationChecks command")
	err = c.conn.SendCommand([]api.Command{
		command,
	})
	if err != nil {
		return errors.Wrap(err, "failed to send LocationChecks command")
	}
	return nil
}
func (c *Client) LocationScouts(ctx context.Context, data api.LocationScouts) (err error) {
	command := api.Command{}
	command.MergeLocationScouts(data)
	fmt.Println("Sending LocationScouts command")
	err = c.conn.SendCommand([]api.Command{
		command,
	})
	if err != nil {
		return errors.Wrap(err, "failed to send LocationScouts command")
	}
	return nil
}
func (c *Client) Say(ctx context.Context, data api.Say) (err error) {
	command := api.Command{}
	command.MergeSay(data)
	fmt.Println("Sending Say command")
	err = c.conn.SendCommand([]api.Command{
		command,
	})
	if err != nil {
		return errors.Wrap(err, "failed to send Say command")
	}
	return nil
}
func (c *Client) Set(ctx context.Context, data api.Set) (err error) {
	command := api.Command{}
	command.MergeSet(data)
	fmt.Println("Sending Set command")
	err = c.conn.SendCommand([]api.Command{
		command,
	})
	if err != nil {
		return errors.Wrap(err, "failed to send Set command")
	}
	return nil
}
func (c *Client) SetNotify(ctx context.Context, data api.SetNotify) (err error) {
	command := api.Command{}
	command.MergeSetNotify(data)
	fmt.Println("Sending SetNotify command")
	err = c.conn.SendCommand([]api.Command{
		command,
	})
	if err != nil {
		return errors.Wrap(err, "failed to send SetNotify command")
	}
	return nil
}
func (c *Client) StatusUpdate(ctx context.Context, data api.StatusUpdate) (err error) {
	command := api.Command{}
	command.MergeStatusUpdate(data)
	fmt.Println("Sending StatusUpdate command")
	err = c.conn.SendCommand([]api.Command{
		command,
	})
	if err != nil {
		return errors.Wrap(err, "failed to send StatusUpdate command")
	}
	return nil
}
func (c *Client) Sync(ctx context.Context, data api.Sync) (err error) {
	command := api.Command{}
	command.MergeSync(data)
	fmt.Println("Sending Sync command")
	err = c.conn.SendCommand([]api.Command{
		command,
	})
	if err != nil {
		return errors.Wrap(err, "failed to send Sync command")
	}
	return nil
}
func (c *Client) UpdateHint(ctx context.Context, data api.UpdateHint) (err error) {
	command := api.Command{}
	command.MergeUpdateHint(data)
	fmt.Println("Sending UpdateHint command")
	err = c.conn.SendCommand([]api.Command{
		command,
	})
	if err != nil {
		return errors.Wrap(err, "failed to send UpdateHint command")
	}
	return nil
}
