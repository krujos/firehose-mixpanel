Firehose to MixPanel
====================


Dump the Firehose to MixPanel, lets see what happens.


#Setup the token for uaa
$ uaac client add watchman --scope uaa.none --authorized_grant_types "client_credentials" --authorities doppler.firehose --redirect_uri http://example.com

#Setup the service for uaa
`cf cups uaa -p '{ "uri": "https://uaa.10.244.0.34.xip.io/oauth/token?grant_type=client_credentials", "client_id": "f2mp", "client_secret": "f2mp" }'`

#Setup the service for MixPanel
`cf cups mixpanel -p '{ "uri": "http://api.mixpanel.com/track/", "api_key": "your key", "api_secret": "your secret" }'`

#Setup the service for doppler
`cf cups doppler -p '{"uri": "wss://doppler.10.244.0.34.xip.io:443" }'`

#If you're using bosh lite, or otherwise skipping ssl validation
`cf set-env f2mp SSL_VERIFY false`
