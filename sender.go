package main

import (
	"encoding/json"
	"log"

	"github.com/cloudfoundry-community/cfenv"
	"github.com/cloudfoundry/sonde-go/events"
)

//SendEventsToMixPanel does batch posts of firehose events to mix channel
func SendEventsToMixPanel(mixPanel *cfenv.Service, msgChan chan *events.Envelope) {
	for msg := range msgChan {
		//eventType := msg.GetEventType()
		//time := msg.GetTimestamp() / 1000000000
		event, err := json.Marshal(msg)
		dieIfError("Could not marshal json", err)
		log.Println(event)
	}
}
