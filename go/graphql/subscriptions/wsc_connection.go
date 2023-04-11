// graphql subscription protocol details

package subscriptions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

const (
	connectionInitMsg = "connection_init" // Client -> Server
	startMsg          = "start"           // Client -> Server
	connectionAckMsg  = "connection_ack"  // Server -> Client
	connectionKaMsg   = "ka"              // Server -> Client
	dataMsg           = "data"            // Server -> Client
	errorMsg          = "error"           // Server -> Client
)

type operationMessage struct {
	Payload json.RawMessage `json:"payload,omitempty"`
	ID      string          `json:"id,omitempty"`
	Type    string          `json:"type"`
}

// Request represents an outgoing GraphQL request
type request struct {
	Query         string                 `json:"query"`
	OperationName string                 `json:"operationName,omitempty"`
	Extensions    map[string]interface{} `json:"extensions,omitempty"`
}

func (wsc *websocketClient) subscribe() error {

	u := url.URL{Scheme: "wss", Host: wsc.domainName, Path: "/swap/graphql"}
	h := http.Header{}
	h.Set("Content-Type", "application/json")

	bodyStruct := request{
		Query: wsc.subscriptionQuery,
	}

	bodyBytes, err := json.Marshal(bodyStruct)
	if err != nil {
		return fmt.Errorf("encode: %w", err)
	}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), h)
	if err != nil {
		return fmt.Errorf("dial: %w", err)
	}

	initMessage := operationMessage{Type: connectionInitMsg}
	if err = c.WriteJSON(initMessage); err != nil {
		return fmt.Errorf("init: %w", err)
	}

	var ack operationMessage
	if err = c.ReadJSON(&ack); err != nil {
		return fmt.Errorf("ack: %w", err)
	}

	if ack.Type != connectionAckMsg {
		return fmt.Errorf("expected ack message, got %#v", ack)
	}

	var ka operationMessage
	if err = c.ReadJSON(&ka); err != nil {
		return fmt.Errorf("ka: %w", err)
	}

	if ka.Type != connectionKaMsg {
		return fmt.Errorf("expected ka message, got %#v", ka)
	}

	if err = c.WriteJSON(operationMessage{Type: startMsg, ID: "1", Payload: bodyBytes}); err != nil {
		return fmt.Errorf("start: %w", err)
	}

	wsc.connection = c

	return nil
}

// gets goroutines, one per connection. Upon failure / restart, the old one dies and it gets a new one
func (wsc *websocketClient) read(iteration int) {
	wsc.lg.Tracef("reading")

	for {
		var op operationMessage
		err := wsc.connection.ReadJSON(&op)
		if err != nil {
			err = fmt.Errorf("read error: %w ", err)
			wsc.readErrorChannel <- wrappedError{iteration: iteration, err: err}
			return
		}

		switch op.Type {
		case dataMsg:
			wsc.lg.Tracef("got data")
			wsc.readerToProcessorChannel <- wrappedReceivedMessage{iteration: iteration, message: op.Payload}
		case connectionKaMsg:
			wsc.lg.Tracef("got ka")
		case errorMsg:
			err := fmt.Errorf("protocol-level error message: %s", string(op.Payload))
			wsc.readErrorChannel <- wrappedError{iteration: iteration, err: err}
			return
		default:
			err := fmt.Errorf("expected data message, got %#v", op)
			wsc.readErrorChannel <- wrappedError{iteration: iteration, err: err}
			return
		}
	}
}

func (wsc *websocketClient) finalizeConnection() {

	wsc.connection.SetWriteDeadline(time.Now().Add(1 * time.Second))
	wsc.lg.Debug("writing close message ...")
	err := wsc.connection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err == nil {
		wsc.lg.Debug("sent close message")

	} else if !wsc.checkInterrupted() {
		wsc.lg.Warnf("write close message error: %w", err)
	}

	_ = wsc.connection.Close()
}
