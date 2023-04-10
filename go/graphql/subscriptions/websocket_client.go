package subscriptions

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"

	"deepwaters/go-examples/graphql/swap"
	"deepwaters/go-examples/util"
)

type wrappedError struct {
	iteration int
	err       error
}

type wrappedReceivedMessage struct {
	iteration int
	message   json.RawMessage
}

type websocketClient struct {
	lg *log.Entry

	logFields log.Fields

	envName    string
	domainName string
	feedName   string

	subscriptionQuery     string
	subscriptionVariables map[string]interface{}

	connection *(websocket.Conn)

	// go to the gatherer
	l3OutputChannel     chan swap.Order
	l2OutputChannel     chan swap.OrderBookLevelUpdate
	tradesOutputChannel chan swap.Trade

	readErrorChannel         chan wrappedError
	readerToProcessorChannel chan wrappedReceivedMessage
	restartRequiredChannel   chan struct{}

	iteration   int
	interrupted bool

	mu sync.Mutex
}

func (wsc *websocketClient) markInterrupted() {
	wsc.mu.Lock()
	wsc.interrupted = true
	wsc.mu.Unlock()
}

func (wsc *websocketClient) checkInterrupted() bool {
	wsc.mu.Lock()
	val := wsc.interrupted
	wsc.mu.Unlock()
	return val
}

func NewWebsocketClient(lg *log.Entry, envName, domainName, feedName string, subscriptionVariables map[string]interface{}) (*websocketClient, error) {

	wsc := websocketClient{
		envName:                  envName,
		domainName:               domainName,
		feedName:                 feedName,
		subscriptionVariables:    subscriptionVariables,
		readerToProcessorChannel: make(chan wrappedReceivedMessage, 1024),
		restartRequiredChannel:   make(chan struct{}),
		readErrorChannel:         make(chan wrappedError),
	}

	wsc.logFields = log.Fields{
		"sourceType": "websocketClient",
		"envName":    envName,
		"feedName":   feedName,
	}

	wsc.lg = lg.WithFields(wsc.logFields)

	switch {
	case feedName == "L3":
		if err := wsc.setupL3(); err != nil {
			return nil, err
		}
	case feedName == "L2":
		if err := wsc.setupL2(); err != nil {
			return nil, err
		}
	case feedName == "trades":
		if err := wsc.setupTrades(); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported feed %s", feedName)
	}
	wsc.lg.Tracef("%s", wsc.subscriptionQuery)

	wsc.lg.Info("created")
	return &wsc, nil
}

func (wsc *websocketClient) processMessage(wMessage wrappedReceivedMessage) {
	switch {
	case wsc.feedName == "L3":
		wsc.parseAndOutputL3(wMessage)
	case wsc.feedName == "L2":
		wsc.parseAndOutputL2(wMessage)
	case wsc.feedName == "trades":
		wsc.parseAndOutputTrades(wMessage)
	}
}

// gets a goroutine
func (wsc *websocketClient) processMessages(haltC <-chan struct{}, doneC chan<- struct{}) {
	for running := true; running; {
		select {
		case <-haltC:
			running = false
		case wMessage := <-wsc.readerToProcessorChannel:
			wsc.processMessage(wMessage)
		}
	}
	close(doneC)
}

// gets a goroutine
func (wsc *websocketClient) connectAndMonitor(haltC <-chan struct{}, doneC chan<- struct{}) {
	if wsc.iteration == 0 {
		wsc.lg.Info("started")
	} else {
		wsc.lg.Info("restarted")
	}

	if err := wsc.subscribe(); err != nil { // sets wsc.connection
		wsc.lg.Errorf("dial/subscribe error, retrying: %w", err)
		time.Sleep(1 * time.Second)
		wsc.connectAndMonitor(haltC, doneC)
		return
	}

	go wsc.read(wsc.iteration)

	haveReadFromHaltC := false
	for running := true; running; {
		select {
		case <-haltC: // if we get here, wsc has been marked interrupted
			haveReadFromHaltC = true
			running = false
		case wErr := <-wsc.readErrorChannel:
			if wErr.iteration != wsc.iteration {
				continue
			}
			wsc.lg.Errorf("%w ", wErr.err)
			running = false
		case <-wsc.restartRequiredChannel:
			wsc.lg.Debug("received restart instruction")
			running = false
		}
	}

	wsc.finalizeSession(haveReadFromHaltC, haltC, doneC)
}

func (wsc *websocketClient) finalizeSession(haveReadFromHaltC bool, haltC <-chan struct{}, doneC chan<- struct{}) {

	wsc.finalizeConnection() // closes the connection, which leads to the "read" goroutine ending

	if haveReadFromHaltC {
		close(doneC)
	} else if wsc.checkInterrupted() {
		<-haltC // needs to be drained
		close(doneC)
	} else {
		wsc.lg.WithFields(wsc.logFields).Info("restarting ...")
		wsc.iteration += 1
		wsc.connectAndMonitor(haltC, doneC)
	}
}

func (wsc *websocketClient) Run(haltC <-chan struct{}, doneC chan<- struct{}) {

	stopper := util.CloseChannelsWrapper{}
	go wsc.processMessages(stopper.NewHaltC(), stopper.NewDoneC())
	go wsc.connectAndMonitor(stopper.NewHaltC(), stopper.NewDoneC())
	<-haltC

	wsc.lg.Debug("shutting down ...")
	wsc.markInterrupted()
	wsc.lg.Trace("sending and receiving done signals ...")
	stopper.SendHaltAndReceiveDoneSignals()
	wsc.lg.Info("shut down gracefully")
	close(doneC)
}
