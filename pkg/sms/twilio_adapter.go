package sms

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// TwilioAdapter wraps around Twilio's phone number verification APIs
// to provide methods for sending verification codes to phone numbers
// and for checking if users provided the correct codes received via SMS
type TwilioAdapter struct {
	client *http.Client
	host   string
	locale string
	apiKey string
}

// NewTwilioAdapter returns a pointer to a value of TwilioTwilioAdapter
func NewTwilioAdapter(client *http.Client, host, locale, apiKey string) *TwilioAdapter {
	return &TwilioAdapter{
		client: client,
		host:   host,
		locale: locale,
		apiKey: apiKey,
	}
}

// SendCode sends a X-digits verification code via SMS to authenticate
// and validate the phoneNumber given by the user.
func (tv *TwilioAdapter) SendCode(countryCode, phoneNumber string) error {
	form := url.Values{}
	form.Add("via", "sms")
	form.Add("country_code", countryCode)
	form.Add("phone_number", phoneNumber)
	form.Add("locale", tv.locale)

	endpoint := fmt.Sprintf("%s/protected/json/phones/verification/start", tv.host)
	req, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-Authy-API-Key", tv.apiKey)

	res, err := tv.client.Do(req)
	if err != nil {
		return err
	}

	_, _ = io.Copy(ioutil.Discard, res.Body)
	res.Body.Close()

	if res.StatusCode >= 400 {
		return errors.New(res.Status)
	}

	return nil
}
