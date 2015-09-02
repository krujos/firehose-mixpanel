package main

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/cloudfoundry-community/cfenv"
	"github.com/cloudfoundry/sonde-go/events"
)

var mixPanelChanel = make(chan *[]byte, 50)

//GetMixPanelChan returns the channel to send events to MixPanel, used as a
//test hook
func GetMixPanelChan() chan *[]byte {
	return mixPanelChanel
}

//Sender interface is what you must implmenet to send something to mixpanel
type Sender interface {
	Send(bytes []byte) error
}

//MixPanelSender sends to MixPanel
type MixPanelSender struct{}

//Send to MixPanel
func (m MixPanelSender) Send(bytes []byte) error {
	log.Fatal("NYI")
	return errors.New("NYOI")
}

//SendEventsToMixPanel does batch posts of firehose events to mix channel
func SendEventsToMixPanel(mixPanel *cfenv.Service, msgChan chan *events.Envelope) {
	for i := 0; i < 3; i++ {
		//		go MixPanelWorker(strconv.Itoa(i), MixPanelSender{})
	}

	for msg := range msgChan {
		j := EventToJSON(msg)
		mixPanelChanel <- j
	}
}

//EventToJSON turns a firehose event into a json representation
func EventToJSON(event *events.Envelope) *[]byte {
	data := map[string]interface{}{
		"event":      event.String(),
		"time":       event.GetTimestamp() / 1000000000,
		"origin":     event.GetOrigin(),
		"deployment": event.GetDeployment(),
		"job":        event.GetJob(),
		"index":      event.GetIndex(),
		"ip":         event.GetIp(),
	}
	j, err := json.Marshal(data)
	if nil != err {
		log.Print("Failed to marshal event")
		log.Print(data)
	}
	return &j
}

//Collect gathers 50 events from the channel and returns
//them as a batch
func Collect(channel chan *[]byte) []byte {
	events := "["
	count := 0
	for {
		event := <-channel
		events += string(*event)
		count++
		if 50 == count {
			log.Println("Received 50 events!")
			events += "]"
			return []byte(events)
		}
		events += ","
	}
}

//MixPanelWorker collects events to send to mixpanel in batches of 50
// func MixPanelWorker(id string, sender Sender) {
// 	log.Println("Created a sender with id " + id)
//   for {
//     batch := Collect(mixPanelChanel)
//     log.Printf("%v", batch)
//   }
// }
