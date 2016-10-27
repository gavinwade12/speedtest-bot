package main

import (
		"strconv"
		"regexp"
		"encoding/base64"
		"crypto/rand"
		"io/ioutil"
		"time"
		"net/url"
		"crypto/sha1"
		"crypto/hmac"
		"strings"
)

func QEscape(toEscape string) string {
	escaped := url.QueryEscape(toEscape)
	return strings.Replace(escaped, "+", "%20", -1)
}

func GetHeaderValues() *Header {
	var header Header
	header.SignatureMethod = "HMAC-SHA1"
	header.Version = "1.0"

	consumerKey := AssignHeaderValues("consumer_key.txt")
	consumerSecret := AssignHeaderValues("consumer_secret.txt")
	accessToken := AssignHeaderValues("access_token.txt")
	accessTokenSecret := AssignHeaderValues("access_token_secret.txt")

	header.ConsumerKey = string(consumerKey)
	header.Token = string(accessToken)
	header.SigningKey = QEscape(string(consumerSecret)) + "&" + QEscape(string(accessTokenSecret))
	header.PartialBaseString = "POST&" + QEscape("https://api.twitter.com/1.1/statuses/update.json") + "&"
	return &header
}

func AssignHeaderValues(fileName string) []byte {
	headerVal, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	return headerVal
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

	paramString := "include_entities=true&oauth_consumer_key=" + QEscape(header.ConsumerKey) + "&oauth_nonce=" + QEscape(header.Nonce) + "&oauth_signature_method=" + header.SignatureMethod +
		"&oauth_timestamp=" + header.Timestamp + "&oauth_token=" + QEscape(header.Token) + "&oauth_version=" + header.Version + "&status=" + status

	base := []byte(header.PartialBaseString + QEscape(paramString))
	h := hmac.New(sha1.New, []byte(header.SigningKey))
	h.Write(base)
	header.Signature = QEscape(base64.StdEncoding.EncodeToString(h.Sum(nil)))
	return header
}

func GetCompleteHeaderString(header Header) string {
	return "OAuth oauth_consumer_key=\"" + QEscape(header.ConsumerKey) + "\", oauth_nonce=\"" + QEscape(header.Nonce) + "\", oauth_signature=\"" + header.Signature + "\", oauth_signature_method=\"" +
		header.SignatureMethod + "\", oauth_timestamp=\"" + header.Timestamp + "\", oauth_token=\"" + QEscape(header.Token) + "\", oauth_version=\"" + header.Version + "\""
}
