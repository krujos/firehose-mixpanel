package main

import (
	"log"

	"github.com/cloudfoundry-community/cfenv"
	"github.com/cloudfoundry/sonde-go/events"
)

var mixPanelChanel = make(chan map[string]interface{}, 50)

//SendEventsToMixPanel does batch posts of firehose events to mix channel
func SendEventsToMixPanel(mixPanel *cfenv.Service, msgChan chan *events.Envelope) {
	for i := 0; i < 3; i++ {
		go sender()
	}

	for msg := range msgChan {
		eventType := msg.GetEventType()
		event := map[string]interface{}{
			"EventType":  eventType,
			"Time":       msg.GetTimestamp() / 1000000000,
			"Origin":     msg.GetOrigin(),
			"Deployment": msg.GetDeployment(),
			"Job":        msg.GetJob(),
			"Index":      msg.GetIndex(),
			"Ip":         msg.GetIp(),
		}
		mixPanelChanel <- event
	}
}

func sender() {
	log.Println("Created a sender!")

	events := make([]map[string]interface{}, 50)
	for i := 0; i < 50; i++ {
		log.Println("GotOne")
		events = append(events, <-mixPanelChanel)
	}
	log.Println("Received 50 events!")
}
