package sbis

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type ListKKTbyOrganization struct {
	client *Client
}

type ListKKTResponse struct {
	RegID     string `json:"regId"`
	Model     string `json:"model"`
	FactoryID string `json:"factoryId"`
	Address   string `json:"address"`
	Status    string `json:"status"`
}

type OrgType int

const (
	OrgNotRegistered       OrgType = 0
	RegistrationInProgress OrgType = 1
	Activated              OrgType = 2
	Deregistered           OrgType = 3
	ActivationExpectation  OrgType = 4
)

func NewListKKTbyOrganization(client *Client) *ListKKTbyOrganization {
	return &ListKKTbyOrganization{client: client}
}

func (c *ListKKTbyOrganization) Get(orgType *OrgType) ([]*ListKKTResponse, error) {
	var ar []*ListKKTResponse

	if len(c.client.inn) == 0 {
		return nil, ErrEmptyINN
	}

	path := fmt.Sprintf("ofd/v1/orgs/%s/kkts", c.client.inn)
	if orgType != nil {
		path = fmt.Sprintf("ofd/v1/orgs/%s/kkts?status=%d", c.client.inn, orgType)
	}
	req, err := c.client.NewRequest(true, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.client.Do(req)
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

	if c.client.verbose {
		c.client.logger.Info().Msgf("%s - success", path)
	}
	return ar, nil
}
