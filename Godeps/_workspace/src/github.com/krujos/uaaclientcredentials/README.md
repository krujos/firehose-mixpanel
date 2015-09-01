#UAA Client Credentials
This package provides a go implementation for the client credentials flow for CloudFoundry UAA. It will exchange a client id and a client secret for a bearer token. 

##Usage

```go
import (
	"github.com/krujos/uaaclientcredentials"
)

func main() {
	uaaURL, err := url.Parse("https://uaa.10.244.0.34.xip.io")
	if nil != err {
		panic("Failed to parse uaa url!")
	}

	creds, err := uaaclientcredentials.New(uaaURL, true, "my_client", "my_secret")
	if nil != err {
		panic("Failed to obtain creds!")
	}
	
	token, err := creds.GetBearerToken()
	if nil != err {
		panic(err)
	}
	
	consumer.DoSomething(token)

}
	
```

See [watchman](https://github.com/krujos/watchman) for a slightly less contrived example.