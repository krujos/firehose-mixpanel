package main_test

import (
	"encoding/json"
	"log"

	"github.com/cloudfoundry/sonde-go/events"
	. "github.com/krujos/firehose-mixpanel"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Sender", func() {
	var event *events.Envelope
	origin := "origin"
	var timestamp int64 = 2000000000
	deployment := "deployment"
	ip := "10.10.10.1"
	job := "job"
	index := "0"

	BeforeEach(func() {
		event = &events.Envelope{
			Origin:     &origin,
			Timestamp:  &timestamp,
			Deployment: &deployment,
			Ip:         &ip,
			Job:        &job,
			Index:      &index,
		}
	})

	It("Should translate the event into json", func() {
		j := EventToJSON(event)

		Ω(j).ShouldNot(BeNil())
		var marshaled interface{}
		err := json.Unmarshal(*j, &marshaled)
		log.Print(err)

		Ω(marshaled).ShouldNot(BeNil())
		Ω(err).Should(BeNil())

	})

})
