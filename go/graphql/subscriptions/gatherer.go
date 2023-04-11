package subscriptions

import (
	"encoding/json"
	"os"
	"os/signal"

	log "github.com/sirupsen/logrus"

	"deepwaters/go-examples/util"
)

type Gatherer struct {
	lg *log.Entry

	envName    string
	domainName string

	l3WebsocketClient     *websocketClient
	l2WebsocketClient     *websocketClient
	tradesWebsocketClient *websocketClient
}

func NewGatherer(lg *log.Logger, envName, domainName string) *Gatherer {

	lge := lg.WithFields(log.Fields{"sourceType": "gatherer", "envName": envName})
	g := Gatherer{lg: lge,
		envName:    envName,
		domainName: domainName}

	return &g
}

func (g *Gatherer) SetL3WebsocketClient(baseAssetID, quoteAssetID, customerAddress string) error { // more arguments are possible - see swap.OrderFilter

	variables := make(map[string]interface{})
	if baseAssetID != "" {
		variables["baseAssetID"] = baseAssetID
	}
	if quoteAssetID != "" {
		variables["quoteAssetID"] = quoteAssetID
	}
	if customerAddress != "" {
		variables["customerAddress"] = customerAddress
	}

	l3WebsocketClient, err := NewWebsocketClient(g.lg, g.envName, g.domainName, "L3", variables)
	if err != nil {
		return err
	}
	g.l3WebsocketClient = l3WebsocketClient
	return nil
}

func (g *Gatherer) SetL2WebsocketClient(baseAssetID, quoteAssetID string) error {

	variables := make(map[string]interface{})
	variables["baseAssetID"] = baseAssetID
	variables["quoteAssetID"] = quoteAssetID

	l2WebsocketClient, err := NewWebsocketClient(g.lg, g.envName, g.domainName, "L2", variables)
	if err != nil {
		return err
	}
	g.l2WebsocketClient = l2WebsocketClient
	return nil
}

func (g *Gatherer) SetTradesWebsocketClient(baseAssetID, quoteAssetID, customerAddress string) error { // more arguments are possible - see swap.TradeFilter

	variables := make(map[string]interface{})
	if baseAssetID != "" {
		variables["baseAssetID"] = baseAssetID
	}
	if quoteAssetID != "" {
		variables["quoteAssetID"] = quoteAssetID
	}
	if customerAddress != "" {
		variables["customerAddress"] = customerAddress
	}

	tradesWebsocketClient, err := NewWebsocketClient(g.lg, g.envName, g.domainName, "trades", variables)
	if err != nil {
		return err
	}
	g.tradesWebsocketClient = tradesWebsocketClient
	return nil
}

func (g *Gatherer) handleL3Updates(haltC <-chan struct{}, doneC chan<- struct{}) {
	for running := true; running; {
		select {
		case <-haltC:
			running = false
		case order := <-g.l3WebsocketClient.l3OutputChannel:
			orderBytes, err := json.Marshal(order)
			if err != nil {
				g.lg.Errorf("error: %s", err)
			} else {
				g.lg.Debugf("order: %s", orderBytes)
			}
		}
	}
	close(doneC)
}

func (g *Gatherer) handleL2Updates(haltC <-chan struct{}, doneC chan<- struct{}) {
	for running := true; running; {
		select {
		case <-haltC:
			running = false
		case orderBookLevelUpdate := <-g.l2WebsocketClient.l2OutputChannel:
			updateBytes, err := json.Marshal(orderBookLevelUpdate)
			if err != nil {
				g.lg.Errorf("error: %s", err)
			} else {
				g.lg.Debugf("level update: %s", updateBytes)
			}
		}
	}
	close(doneC)
}

func (g *Gatherer) handleTradeUpdates(haltC <-chan struct{}, doneC chan<- struct{}) {
	for running := true; running; {
		select {
		case <-haltC:
			running = false
		case trade := <-g.tradesWebsocketClient.tradesOutputChannel:
			tradeBytes, err := json.Marshal(trade)
			if err != nil {
				g.lg.Errorf("error: %s", err)
			} else {
				g.lg.Debugf("trade: %s", tradeBytes)
			}
		}
	}
	close(doneC)
}

func (g *Gatherer) RestartWebsocketClient(feedName string) {
	if feedName == "L3" && g.l3WebsocketClient != nil {
		g.l3WebsocketClient.restartRequiredChannel <- struct{}{}
	}

	if feedName == "L2" && g.l2WebsocketClient != nil {
		g.l2WebsocketClient.restartRequiredChannel <- struct{}{}
	}

	if feedName == "trades" && g.tradesWebsocketClient != nil {
		g.tradesWebsocketClient.restartRequiredChannel <- struct{}{}
	}
}

func (g *Gatherer) Run() {

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	stopper := util.CloseChannelsWrapper{}

	if g.l3WebsocketClient != nil {
		go g.handleL3Updates(stopper.NewHaltC(), stopper.NewDoneC())
		go g.l3WebsocketClient.Run(stopper.NewHaltC(), stopper.NewDoneC())
	}

	if g.l2WebsocketClient != nil {
		go g.handleL2Updates(stopper.NewHaltC(), stopper.NewDoneC())
		go g.l2WebsocketClient.Run(stopper.NewHaltC(), stopper.NewDoneC())
	}

	if g.tradesWebsocketClient != nil {
		go g.handleTradeUpdates(stopper.NewHaltC(), stopper.NewDoneC())
		go g.tradesWebsocketClient.Run(stopper.NewHaltC(), stopper.NewDoneC())
	}

	<-interrupt
	g.lg.Debug("interrupt. shutting down ...")
	g.lg.Trace("sending and receiving done signals ...")
	stopper.SendHaltAndReceiveDoneSignals()
	g.lg.Info("shut down gracefully")
}
