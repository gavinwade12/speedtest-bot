package main

type Speeds struct {
	Download string `json:'download'`
	Upload  string `json:'upload'`
}

type Header struct {
	ConsumerKey string
	Nonce string
	Signature string
	SignatureMethod string
	Timestamp string
	Token string
	Version string
	SigningKey string
	PartialBaseString string
}