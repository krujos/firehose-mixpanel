package main

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/cloudfoundry-community/cfenv"
	"github.com/cloudfoundry/sonde-go/events"
)

var mixPanelChanel = make(chan map[string]interface{}, 50)

//SendEventsToMixPanel does batch posts of firehose events to mix channel
func SendEventsToMixPanel(mixPanel *cfenv.Service, msgChan chan *events.Envelope) {
	for i := 0; i < 3; i++ {
		go mixPanelWorker(strconv.Itoa(i))
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

//EventToJSON turns a firehose event into a json representation
func EventToJSON(event *events.Envelope) []byte {
	return []byte("garbage")
}

func mixPanelWorker(id string) {
	log.Println("Created a sender with id " + id)
	events := make([]map[string]interface{}, 50)
	count := 0
	for {
		events = append(events, <-mixPanelChanel)
		count++
		if 50 == count {
			log.Println(id + " Received 50 events!")
			count = 0
			events = make([]map[string]interface{}, 50)
			j, _ := json.Marshal(events)
			log.Printf("%v", j)
		}
	}
}
