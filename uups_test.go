package main_test

import (
	"github.com/cloudfoundry-community/cfenv"
	. "github.com/krujos/firehose-mixpanel"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Uups", func() {

	validEnv := []string{
		`VCAP_APPLICATION={"instance_id":"451f045fd16427bb99c895a2649b7b2a","instance_index":0,"host":"0.0.0.0","port":61857,"started_at":"2013-08-1200:05:29+0000","started_at_timestamp":1376265929,"start":"2013-08-1200:05:29+0000","state_timestamp":1376265929,"limits":{"mem":512,"disk":1024,"fds":16384},"application_version":"c1063c1c-40b9-434e-a797-db240b587d32","application_name":"styx-james","application_uris":["styx-james.a1-app.cf-app.com"],"version":"c1063c1c-40b9-434e-a797-db240b587d32","name":"styx-james","uris":["styx-james.a1-app.cf-app.com"],"users":null}`,
		`HOME=/home/vcap/app`,
		`MEMORY_LIMIT=512m`,
		`PWD=/home/vcap`,
		`TMPDIR=/home/vcap/tmp`,
		`USER=vcap`,
		`VCAP_SERVICES={"user-provided":[{"credentials":{"api_key":"the-key","api_secret":"the-secret","uri":"http://api.mixpanel.com/track/"},"label":"user-provided","name":"mixpanel","syslog_drain_url":"","tags":[]},{"credentials":{"client_id":"f2mp",	"client_secret":"f2mp",	"uri":"https://uaa.10.244.0.34.xip.io/oauth/token?grant_type=client_credentials"},"label":"user-provided","name":"uaa","syslog_drain_url":"","tags":[]}]}`,
	}

	It("Should get a service named mixpanel", func() {
		testEnv := cfenv.Env(validEnv)
		cfenv, err := cfenv.New(testEnv)
		Ω(err).Should(BeNil())
		Ω(cfenv).ShouldNot(BeNil())

		mixPanel, err := GetUserProvidedServiceByName("mixpanel", cfenv)
		Ω(mixPanel).ShouldNot(BeNil())
		Ω(err).Should(BeNil())
	})

	It("Should get a service named uaa", func() {
		testEnv := cfenv.Env(validEnv)
		cfenv, err := cfenv.New(testEnv)
		Ω(err).Should(BeNil())
		Ω(cfenv).ShouldNot(BeNil())

		uaa, err := GetUserProvidedServiceByName("uaa", cfenv)
		Ω(uaa).ShouldNot(BeNil())
		Ω(err).Should(BeNil())
	})

	It("Should not get a service named foo", func() {
		testEnv := cfenv.Env(validEnv)
		cfenv, err := cfenv.New(testEnv)
		Ω(err).Should(BeNil())
		Ω(cfenv).ShouldNot(BeNil())

		foo, err := GetUserProvidedServiceByName("foo", cfenv)
		Ω(err).ShouldNot(BeNil())
		Ω(foo).Should(BeNil())
	})
})
