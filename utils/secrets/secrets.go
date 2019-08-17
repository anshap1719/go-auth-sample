package secrets

// These options need to change on deploy.
// These are all development keys; it is ok for them to be stored on Github.

const (
	Scheme   = "http"
	Hostname = "localhost:3000"
	URL      = Scheme + "://" + Hostname

	RecaptchaSiteKey  = ""
	RecaptchaSecret   = ""
	KeyLocation       = ""
	IsRelease         = false
	GoogleMapsKey     = ""
	GoogleMapsSigning = ""

	GoogleID        = ""
	GoogleSecret    = ""
	AmazonID        = ""
	AmazonSecret    = ""
	FacebookID      = ""
	FacebookSecret  = ""
	LinkedInID      = ""
	LinkedInSecret  = ""
	MicrosoftID     = ""
	MicrosoftSecret = ""
	TwitterKey      = ""
	TwitterSecret   = ""

	SendgridAPIKey = ""

	// Authentication keys
	JWTPublicKey = `-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAutkfBmxsO6eV8zDxaNwf
qcJTYcxHVrdOkX3RB9bedcNuzx8h2bat/FxulwzvgsazTHxZAyIuBFzYHNXTvUel
YuNrT4uTlE55EZAISkCXB4UkBkibMWdsr+SMI2cOzDXyxyc4YuzReazkoNnyWokb
SPeHepNYN+wiR16mXtsR2D7FcH1kgxAOnTCSz+vqDG2p9TKkYqCV6h5QGteI/CLf
2M3HQHWOk7a+4S7c267rx8WCb+Wzn278MQJl6/eBhCv5hPAALriQqQXIgQ+CvWMr
NDk00UjGnCDRrTz4An5gzj1Q9PM1RkyU7HMD6gCX5s+pe0F6UDuqDZ+HIQALT9qg
0Q2Kly/+wRKfKouO1ofaXFDVaN+KcNO0t7V23v2s2JyOnaXIkCnoIzh8q8Qbrxwk
mWuKLNlTHzjP0M8rJnprDhcdZV4rB8ettmg7FWjjcWf+W9QOkPI+7/g9kHGqm9EV
ARSWb3s8ooKC2i1/Gz/e13v0/8vxghW0+IwrW6XAEHJH55WFuDNe3f8EeyWY/z2A
Yv/Z5iNr5APhKx0Y7EJHlONw5cPtsDYvdp3PYFQgnXSuAwnmXqzFONoinNY8Dz6v
fp8rdYyJkENdSI9Gd5GTBJWCQW/C5kiPQpwyTmDtEu0c82+8c+dMMVpVct2rmRPe
6hNBiDHQflfHjVHDljrHd3sCAwEAAQ==
-----END PUBLIC KEY-----
`
	JWTPrivateKey = `-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----
`
)
