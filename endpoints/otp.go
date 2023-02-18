package endpoints

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/supabase-community/gotrue-go/types"
)

const otpPath = "/otp"

// POST /otp
// One-Time-Password. Will deliver a magiclink or SMS OTP to the user depending
// on whether the request contains an email or phone key.
//
// If CreateUser is true, the user will be automatically signed up if the user
// doesn't exist.
func (c *Client) OTP(req types.OTPRequest, redirect string) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	r, err := c.newRequest(otpPath, http.MethodPost, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	if redirect != "" {
		q := r.URL.Query()
		q.Add("redirect_to", redirect)
		r.URL.RawQuery = q.Encode()
	}

	resp, err := c.client.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fullBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("response status code %d", resp.StatusCode)
		}
		return fmt.Errorf("response status code %d: %s", resp.StatusCode, fullBody)
	}

	return nil
}
