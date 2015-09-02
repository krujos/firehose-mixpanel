package main_test

import (
	"encoding/json"

	"github.com/cloudfoundry/sonde-go/events"
	. "github.com/krujos/firehose-mixpanel"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type MockSender struct {
	Bytes []byte
}

func (m MockSender) Send(b []byte) error {
	m.Bytes = b
	return nil
}

var _ = Describe("Sender", func() {
	var event *events.Envelope
	origin := "origin"
	var timestamp int64 = 2000000000
	deployment := "deployment"
	ip := "10.10.10.1"
	job := "job"
	index := "0"
	eventType := events.Envelope_HttpStart
	BeforeEach(func() {
		event = &events.Envelope{
			Origin:     &origin,
			Timestamp:  &timestamp,
			Deployment: &deployment,
			Ip:         &ip,
			Job:        &job,
			Index:      &index,
			EventType:  &eventType,
		}
	})

	It("Should translate the event into json", func() {
		j := EventToJSON(event)
		Ω(j).ShouldNot(BeNil())
	})

	It("Should set the proper envelope fields", func() {
		var actual map[string]interface{}
		err := json.Unmarshal(*(EventToJSON(event)), &actual)
		Ω(err).Should(BeNil())
		Ω(actual["origin"]).Should(Equal(origin))
		Ω(actual["ip"]).Should(Equal(ip))
		Ω(actual["job"]).Should(Equal(job))
		Ω(actual["time"]).Should(Equal(float64(2)))
	})

	Describe("The worker", func() {
		It("Should append 50 events", func() {
			mixPanelChan := GetMixPanelChan()
			var batch []byte
			for i := 0; i < 50; i++ {
				bytes := []byte("foo bar")
				mixPanelChan <- &bytes
			}
			batch = Collect(mixPanelChan)
			Ω(batch).NotTo(BeNil())
		})
	})
})
