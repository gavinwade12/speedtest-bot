package main

import (
		"fmt"
		"strconv"
		"regexp"
		"encoding/base64"
		"crypto/rand"
		"io/ioutil"
		"time"
		"net/url"
		"crypto/sha1"
		"crypto/hmac"
)

func QEscape(toEscape string) string {
	return url.QueryEscape(toEscape)
}

func GetHeaderValues() *Header {
	var header Header
	header.SignatureMethod = "HMAC-SHA1"
	header.Version = "1.0"
	
	consumerKey, err := ioutil.ReadFile("consumer_key.txt")
	if err != nil {
		panic(err)
	}
	consumerSecret, err := ioutil.ReadFile("consumer_secret.txt")
	if err != nil {
		panic(err)
	}
	accessToken, err := ioutil.ReadFile("access_token.txt")
	if err != nil {
		panic(err)
	}
	accessTokenSecret, err := ioutil.ReadFile("access_token_secret.txt")
	if err != nil {
		panic(err)
	}

	header.ConsumerKey = string(consumerKey)
	header.Token = string(accessToken)
	header.SigningKey = QEscape(string(consumerSecret)) + "&" + QEscape(string(accessTokenSecret))
	header.PartialBaseString = "POST&" + QEscape("https://api.twitter.com/1.1/statuses/update.json") + "&"
	return &header
}

func UpdateHeaderValues(header Header, status string) Header {
	reg, _ := regexp.Compile("[^a-zA-Z0-9-]")
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		panic(err)
	}
	nonce := base64.StdEncoding.EncodeToString(randomBytes)
	header.Nonce = reg.ReplaceAllString(nonce, "")
	header.Timestamp = strconv.FormatInt(int64(time.Now().Unix()), 10)

	paramString := "oauth_consumer_key=" + QEscape(header.ConsumerKey)
	paramString += "&oauth_nonce=" + QEscape(header.Nonce)
	paramString += "&oauth_signature_method=" + header.SignatureMethod
	paramString += "&oauth_timestamp=" + header.Timestamp
	paramString += "&oauth_token=" + QEscape(header.Token)
	paramString += "&oauth_version=" + header.Version
	paramString += "&status=" + status

	base := []byte(header.PartialBaseString + QEscape(paramString))
	fmt.Printf("%s\n\n", string(base))
	h := hmac.New(sha1.New, []byte(header.SigningKey))
	h.Write(base)
	header.Signature = base64.StdEncoding.EncodeToString(h.Sum(nil))
	return header
}

func GetCompleteHeaderString(header Header) string {
	headerString := "OAuth "
	headerString += "oauth_consumer_key=\"" + QEscape(header.ConsumerKey) + "\", "
	headerString += "oauth_nonce=\"" + QEscape(header.Nonce) + "\", "
	headerString += "oauth_signature=\"" + header.Signature + "\", "
	headerString += "oauth_signature_method=\"" + header.SignatureMethod + "\", "
	headerString += "oauth_timestamp=\"" + header.Timestamp + "\", "
	headerString += "oauth_token=\"" + QEscape(header.Token) + "\", "
	headerString += "oauth_version=\"" + header.Version + "\""
	return headerString
}