package main

import (
	"github.com/cloudfoundry-community/cfenv"
	"github.com/cloudfoundry/sonde-go/events"
)

//SendEventsToMixPanel does batch posts of firehose events to mix channel
func SendEventsToMixPanel(mixPanel *cfenv.Service, msgchan chan *events.Envelope) {

}
