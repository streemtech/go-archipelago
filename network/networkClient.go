package network

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/coder/websocket"
	"github.com/pkg/errors"
	"github.com/streemtech/go-archipelago/api"
)

type Client interface {
	Start() (err error)
	//Send a command to the websocket/archipelago server.
	SendCommand(api.Commands) (err error)
	//close the connection to the websocket/archipelago server.
	Close() (err error)
}

type ClientImpl struct {
	conn *websocket.Conn

	rootCtx context.Context

	//CommandsCallback handles the list of commands that are returned.
	commandsCallback func(ctx context.Context, cmds api.Commands)

	//CloseCallback Optional: will be called when the websocket client is closed.
	//May contain an error if the close was triggered by a send error
	closeCallback func(err error)

	dialOpts *websocket.DialOptions

	closed bool

	//The primary websocket url to connect to.
	url string

	readLimit int64
}

type ClientProps struct {
	//The primary websocket url to connect to.
	Url string

	//CommandsCallback handles the list of commands that are returned.
	CommandsCallback func(ctx context.Context, cmds api.Commands)

	//CloseCallback Optional: will be called when the websocket client is closed.
	//May contain an error if the close was triggered by a send error
	CloseCallback func(err error)

	//Sets the context for the client request.
	RootContext context.Context

	// HTTPClient is used for the connection.
	// Its Transport must return writable bodies for WebSocket handshakes.
	// http.Transport does beginning with Go 1.12.
	HTTPClient *http.Client

	// HTTPHeader specifies the HTTP headers included in the handshake request.
	HTTPHeader http.Header

	// Host optionally overrides the Host HTTP header to send. If empty, the value
	// of URL.Host will be used.
	Host string
}

func (c ClientProps) Validate() (err error) {
	if c.Url == "" {
		return errors.New("Must set URL")
	}
	if c.CommandsCallback == nil {
		return errors.New("Must set CommandsCallback")
	}
	return nil
}

// Test is a testing function to see if I can get archipelago BS to function.
func NewClient(props ClientProps) (client Client, err error) {

	err = props.Validate()
	if err != nil {
		return nil, errors.Wrap(err, "invalid networking client props")
	}

	//set root context to default
	rootCtx := props.RootContext
	if rootCtx == nil {
		rootCtx = context.Background()
	}

	httpClient := props.HTTPClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	ci := &ClientImpl{
		readLimit:        -1,
		commandsCallback: props.CommandsCallback,
		closeCallback:    props.CloseCallback,
		rootCtx:          rootCtx,
		dialOpts: &websocket.DialOptions{
			HTTPClient: httpClient,
			HTTPHeader: props.HTTPHeader,
			Host:       props.Host,
		},
		url: props.Url,
	}

	return ci, nil

}

func (ci *ClientImpl) Start() (err error) {

	conn, _, err := websocket.Dial(ci.rootCtx, ci.url, ci.dialOpts)

	if err != nil {
		return errors.Wrap(err, "failed to dial archipelago server")
	}
	conn.SetReadLimit(ci.readLimit)

	ci.conn = conn

	go func() {
		for {
			//read in something from the websocket
			_, data, readErr := ci.conn.Read(ci.rootCtx)
			if readErr != nil {

				if !ci.closed {
					ci.closed = true
					ci.conn.Close(websocket.StatusGoingAway, "read-error")
				}
				if ci.closeCallback != nil {
					ci.closeCallback(readErr)
				}
				return
			}

			var cmds api.Commands = []api.Command{}
			unmarshalErr := json.Unmarshal(data, &cmds)
			if unmarshalErr != nil {
				if !ci.closed {
					ci.closed = true
					ci.conn.Close(websocket.StatusGoingAway, "unmarshal-error")
				}
				if ci.closeCallback != nil {
					ci.closeCallback(unmarshalErr)
				}
				return
			}

			ci.sendCommandsToUser(cmds)

			select {
			case <-ci.rootCtx.Done():
				return
			default:
			}
		}
	}()

	return nil
}

// sends a command to hte users callback to parse
func (ci *ClientImpl) sendCommandsToUser(cmds api.Commands) {
	defer func() {
		if r := recover(); r != nil {
		}
	}()

	ci.commandsCallback(ci.rootCtx, cmds)

}

// takes a command from the user and sends it to
func (ci *ClientImpl) SendCommand(cmds api.Commands) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Errorf("recovered sending command: %v", r)
		}
	}()

	//marshal to json to send commands
	data, marshalErr := json.Marshal(cmds)
	if marshalErr != nil {
		return errors.Wrap(marshalErr, "failed to marshal commands")
	}

	// fmt.Println(string(data))

	//send commands to destination server.
	writeErr := ci.conn.Write(ci.rootCtx, websocket.MessageText, data)
	if writeErr != nil {
		return errors.Wrap(writeErr, "failed to send message over connection")
	}

	return nil

}

func (ci *ClientImpl) Close() (err error) {
	if !ci.closed {
		ci.closed = true
		return errors.Wrap(ci.conn.Close(websocket.StatusGoingAway, "read-error"), "failed to close websocket connection")
	}
	return nil
}
