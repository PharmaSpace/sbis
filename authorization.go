package sbis

import (
	"encoding/json"
	"io/ioutil"
)

type Authorization struct {
	client *Client
}

type AuthResponse struct {
	Sid   string `json:"sid"`
	Token string `json:"token"`
}

func NewAuthorization(client *Client) *Authorization {
	return &Authorization{client: client}
}

func (a *Authorization) GetSID() (*AuthResponse, error) {
	var ar AuthResponse
	path := "/oauth/service/"

	req, err := a.client.NewRequest(false, "POST", path, &a.client.config)
	if err != nil {
		return nil, err
	}

	resp, err := a.client.client.Do(req)
	if resp != nil {
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(data, &ar)
		if err != nil {
			return nil, err
		}
	}

	if a.client.verbose {
		a.client.logger.Info().Msgf("login success sid %s", ar.Sid)
	}
	return &ar, nil
}
