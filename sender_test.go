package main_test

import (
	"encoding/base64"
	"encoding/json"
	"sync"

	"github.com/cloudfoundry/sonde-go/events"
	. "github.com/krujos/firehose-mixpanel"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

type MockSender struct {
	Bytes []byte
}

func (m MockSender) Send(b []byte) error {
	m.Bytes = b
	return nil
}

func testCollect(wg *sync.WaitGroup, mixPanelChan chan *[]byte) {
	defer wg.Done()
	batch := Collect(mixPanelChan)
	Ω(batch).NotTo(BeNil())
	var actual []interface{}
	err := json.Unmarshal(batch, &actual)
	Ω(err).Should(BeNil())
	Ω(actual).To(HaveLen(50))
}

var _ = Describe("Sender", func() {
	Describe("Envelope Processing", func() {

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
			SetMixPanelToken("mixpanel-token")
		})

		It("Should translate the event into json", func() {
			j := EventToJSON(event)
			Ω(j).ShouldNot(BeNil())
		})

		It("Should set the proper envelope fields", func() {
			var actual map[string]interface{}
			var props map[string]interface{}

			err := json.Unmarshal(*(EventToJSON(event)), &actual)
			Ω(err).Should(BeNil())
			props = actual["properties"].(map[string]interface{})
			Ω(props["origin"]).Should(Equal(origin))
			Ω(props["ip"]).Should(Equal(ip))
			Ω(props["job"]).Should(Equal(job))
			Ω(props["time"]).Should(Equal(float64(2)))
			Ω(props["token"]).Should(Equal("mixpanel-token"))
		})
	})

	Describe("The worker", func() {
		It("Should append 50 events", func() {
			mixPanelChan := GetMixPanelChan()
			for i := 0; i < 50; i++ {
				input := []byte("{\"foo\":\"bar\"}")
				mixPanelChan <- &input
			}
			batch := Collect(mixPanelChan)

			Ω(batch).NotTo(BeNil())
			var actual []interface{}
			err := json.Unmarshal(batch, &actual)
			Ω(err).Should(BeNil())
			Ω(actual).To(HaveLen(50))
		})

		It("Should handle 100 events in chunks of 50", func() {
			//This isn't the greatest test, if itr % 50 != 0 then the channel
			//won't block and the test will pass... it does test that we're handleing
			//in batches
			mixPanelChan := GetMixPanelChan()
			itr := 1000
			var wg sync.WaitGroup
			wg.Add(1 + (itr / 50))
			go func() {
				defer wg.Done()
				for i := 0; i < itr; i++ {
					input := []byte("{\"foo\":\"bar\"}")
					mixPanelChan <- &input
				}
			}()

			for i := itr / 50; i > 0; i-- {
				go testCollect(&wg, mixPanelChan)
			}
			wg.Wait()
		})
	})

	Describe("The sender", func() {
		var server *ghttp.Server
		statusCode := 200
		data := []byte("[{\"foo\":\"bar\"}]")

		BeforeEach(func() {
			encodedString := base64.StdEncoding.EncodeToString(data)
			server = ghttp.NewServer()
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/track"),
					ghttp.VerifyFormKV("data", encodedString),
					ghttp.RespondWith(statusCode, "1"),
				))
		})

		AfterEach(func() {
			server.Close()
		})

		It("should send base64 encoded request requests to mix pannel", func() {
			m := MixPanelSender{URL: server.URL() + "/track"}
			err := m.Send(data)
			Ω(err).Should(BeNil())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
	})

	Describe("The sender gets an error", func() {
		var server *ghttp.Server
		statusCode := 200
		data := []byte("[{\"foo\":\"bar\"}]")

		BeforeEach(func() {
			encodedString := base64.StdEncoding.EncodeToString(data)
			server = ghttp.NewServer()
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/track"),
					ghttp.VerifyFormKV("data", encodedString),
					ghttp.RespondWith(statusCode, "0"),
				))
		})

		AfterEach(func() {
			server.Close()
		})

		It("should send base64 encoded request requests to mix pannel", func() {
			m := MixPanelSender{URL: server.URL() + "/track"}
			err := m.Send(data)
			Ω(err).ShouldNot(BeNil())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
	})

})
