package util

type CloseChannelsWrapper struct {
	haltChannels []chan struct{}
	doneChannels []chan struct{}
}

func (w *CloseChannelsWrapper) NewHaltC() chan struct{} {
	chn := make(chan struct{}, 1)
	w.haltChannels = append(w.haltChannels, chn)
	return chn
}

func (w *CloseChannelsWrapper) NewDoneC() chan struct{} {
	chn := make(chan struct{})
	w.doneChannels = append(w.doneChannels, chn)
	return chn
}

func (w *CloseChannelsWrapper) SendHaltAndReceiveDoneSignals() {
	for _, c := range w.haltChannels {
		c <- struct{}{}
	}
	for _, c := range w.doneChannels {
		<-c
	}
}
