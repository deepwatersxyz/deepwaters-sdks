package subscriptions

import (
	"encoding/json"
	"fmt"

	"deepwaters/go-examples/graphql/accounting"
)

func (wsc *websocketClient) setupBalances() error {
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
		balances%s %s
	}`, variables, accounting.BalanceResponseFields)

	wsc.balancesOutputChannel = make(chan accounting.Balance, 1024)
	return nil
}

func (wsc *websocketClient) parseAndOutputBalances(wMessage wrappedReceivedMessage) {
	var d accounting.BalancesFeedFrame

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

	for _, balance := range d.Data.Balances {
		wsc.balancesOutputChannel <- balance
	}
}
