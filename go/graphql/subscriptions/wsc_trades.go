package subscriptions

import (
	"deepwaters/go-examples/graphql/swap"
	"encoding/json"
	"fmt"
)

func (wsc *websocketClient) setupTrades() error {
	variables := ""
	if len(wsc.subscriptionVariables) != 0 {
		variables = "(where: {"
		for k, v := range wsc.subscriptionVariables {
			variables += k + ": "
			bytes, err := json.Marshal(v)
			if err != nil {
				return nil
			}
			variables += string(bytes) + ", "

		}
		variables = variables[:len(variables)-2] + "})"
	}
	wsc.subscriptionQuery = fmt.Sprintf(`subscription {
		trades%s %s
	}`, variables, swap.TradeResponseFields)

	wsc.tradesOutputChannel = make(chan swap.Trade, 1024)
	return nil
}

func (wsc *websocketClient) parseAndOutputTrades(wMessage wrappedReceivedMessage) {
	var d swap.TradesFeedFrame

	if err := json.Unmarshal(wMessage.message, &d); err != nil {
		err = fmt.Errorf("error parsing output: %w ", err)
		wsc.readErrorChannel <- wrappedError{iteration: wMessage.iteration, err: err}
		return
	}

	if d.Errors != nil {
		err := fmt.Errorf("error in response: %#v ", d.Errors)
		wsc.readErrorChannel <- wrappedError{iteration: wMessage.iteration, err: err}
		return
	}

	for _, trade := range d.Data.Trades {
		wsc.tradesOutputChannel <- trade
	}
}
