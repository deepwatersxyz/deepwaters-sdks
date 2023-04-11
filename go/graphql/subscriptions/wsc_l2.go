package subscriptions

import (
	"deepwaters/go-examples/graphql/swap"
	"encoding/json"
	"fmt"
)

func (wsc *websocketClient) setupL2() error {
	variables := "("
	for k, v := range wsc.subscriptionVariables {
		variables += k + ": "
		bytes, err := json.Marshal(v)
		if err != nil {
			return nil
		}
		variables += string(bytes) + ", "

	}
	variables = variables[:len(variables)-2] + ")"
	wsc.subscriptionQuery = fmt.Sprintf(`subscription {
		orderBook%s %s
	}`, variables, swap.OrderBookLevelUpdateFields)

	wsc.l2OutputChannel = make(chan swap.OrderBookLevelUpdate, 1024)
	return nil
}

func (wsc *websocketClient) parseAndOutputL2(wMessage wrappedReceivedMessage) {
	var d swap.L2FeedFrame

	if err := json.Unmarshal(wMessage.message, &d); err != nil {
		err = fmt.Errorf("error parsing output: %w", err)
		wsc.readErrorChannel <- wrappedError{iteration: wMessage.iteration, err: err}
		return
	}

	if d.Errors != nil {
		err := fmt.Errorf("error in response: %#v", d.Errors)
		wsc.readErrorChannel <- wrappedError{iteration: wMessage.iteration, err: err}
		return
	}

	for _, l2Update := range d.Data.OrderBook {
		wsc.l2OutputChannel <- l2Update
	}
}
