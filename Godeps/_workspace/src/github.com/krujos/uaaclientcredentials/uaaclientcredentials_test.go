package uaaclientcredentials

import (
	"net/http"
	"net/url"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("UAAClientCredentials", func() {
	var url *url.URL
	var uaaCC *UAAClientCredentials

	BeforeEach(func() {
		url, _ = url.Parse("https://uaa.at.your.place.com")
	})

	Describe("Creationism", func() {
		It("makes an initiliazed object", func() {
			uaaCC, _ := New(url, true, "client_id", "client_secret")
			Expect(uaaCC).NotTo(BeNil())
		})

		It("should complain about an empty client id", func() {
			uaaCC, err := New(url, false, "", "client_secret")
			Expect(uaaCC).To(BeNil())
			Expect(err).ToNot(BeNil())
		})

		It("should complain about an empty client secret", func() {
			uaaCC, err := New(url, false, "client_id", "")
			Expect(uaaCC).To(BeNil())
			Expect(err).ToNot(BeNil())
		})

		Describe("token preconditions", func() {
			BeforeEach(func() {
				uaaCC, _ = New(url, true, "client_id", "client_secret")
			})

			It("Should not have a valid token yet", func() {
				Expect(uaaCC.expiresAt.Unix()).To(BeNumerically("<", time.Now().Unix()))
			})

		})
	})

	Describe("SSL Validation", func() {

		It("Should skip ssl validation", func() {
			uaaCC, _ = New(url, true, "client_id", "client_secret")
			config := uaaCC.getTLSConfig()
			Expect(config.InsecureSkipVerify).To(BeTrue())
		})

		It("Should skip not ssl validation", func() {
			uaaCC, _ = New(url, false, "client_id", "client_secret")
			config := uaaCC.getTLSConfig()
			Expect(config.InsecureSkipVerify).To(BeFalse())
		})
	})

	Describe("Token Acquisition", func() {
		var server *ghttp.Server
		var statusCode int
		var responseBody UAATokenResponse

		BeforeEach(func() {
			server = ghttp.NewTLSServer()
			url, _ = url.Parse(server.URL())
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/oauth/token", "grant_type=client_credentials"),
					ghttp.VerifyBasicAuth("client_id", "client_secret"),
					ghttp.RespondWithJSONEncodedPtr(&statusCode, &responseBody),
				),
			)
			uaaCC, _ = New(url, true, "client_id", "client_secret")
		})

		AfterEach(func() {
			server.Close()
		})

		Context("when the request is all good", func() {
			BeforeEach(func() {
				statusCode = http.StatusOK
				responseBody = UAATokenResponse{
					AccessToken: "test_token",
					TokenType:   "bearer",
					ExpiresIn:   43199,
					Scope:       "cloud_controller.admin",
					Jti:         "145450ab-c78f-46dd-90f4-51a40c2bc2c0",
				}
			})

			It("should fetch a credential", func() {
				results, err := uaaCC.getJSON()
				Ω(server.ReceivedRequests()).Should(HaveLen(1))
				Expect(results).ToNot(BeNil())
				Expect(err).To(BeNil())
			})

			It("should marshal UAA responses into creds", func() {
				uaaCC.getToken()
				Ω(server.ReceivedRequests()).Should(HaveLen(1))
				Expect(uaaCC.accessToken).To(Equal("test_token"))
				Expect(uaaCC.expiresAt.Unix()).To(BeNumerically(">", time.Now().Unix()))
			})

			Describe("when the token is expired", func() {
				BeforeEach(func() {
					duration, _ := time.ParseDuration("-5m")
					uaaCC.expiresAt = time.Now().Add(duration)
					Expect(uaaCC.expiresAt.Unix()).To(BeNumerically("<", time.Now().Unix()))
				})

				It("Should refresh my token", func() {
					token, _ := uaaCC.GetBearerToken()
					Expect(token).NotTo(BeNil())
					Expect(uaaCC.expiresAt.Unix()).To(BeNumerically(">", time.Now().Unix()))
				})
			})
		})

		Context("when the request is unauthorized", func() {
			BeforeEach(func() {
				statusCode = http.StatusUnauthorized
				responseBody = UAATokenResponse{}
			})

			It("should return an error if the status code is not 200", func() {
				results, err := uaaCC.getJSON()
				Ω(server.ReceivedRequests()).Should(HaveLen(1))
				Expect(err).ToNot(BeNil())
				Expect(results).To(BeNil())
			})

			It("Should give me an error when refresh fails", func() {
				token, err := uaaCC.GetBearerToken()
				Ω(server.ReceivedRequests()).Should(HaveLen(1))
				Expect(err).ToNot(BeNil())
				Expect(token).To(Equal(""))
			})
		})

	})

	Describe("Bearer Tokens", func() {
		BeforeEach(func() {
			uaaCC, _ = New(url, true, "client_id", "client_secret")
			duration, _ := time.ParseDuration("1h")
			uaaCC.expiresAt = time.Now().Add(duration)
			uaaCC.accessToken = "test_token"
		})

		It("should return a properly formatted bearer token", func() {
			token, _ := uaaCC.GetBearerToken()
			Expect(token).To(Equal("bearer test_token"))
		})
	})
})
