package sbis

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type ListOfFiscalDriver struct {
	client *Client
}

type ListOfFiscalDriverResponse struct {
	StorageID     string `json:"storageId"`
	Model         string `json:"model"`
	Status        int    `json:"status"`
	EffectiveFrom string `json:"effectiveFrom"`
	EffectiveTo   string `json:"effectiveTo"`
}

func NewListOfFiscalDriver(client *Client) *ListOfFiscalDriver {
	return &ListOfFiscalDriver{client: client}
}

func (c *ListOfFiscalDriver) Get(regId string, orgType *OrgType) ([]*ListOfFiscalDriverResponse, error) {
	var ar []*ListOfFiscalDriverResponse

	if len(c.client.inn) == 0 || len(regId) == 0 {
		return nil, ErrEmptyREQ
	}

	path := fmt.Sprintf("ofd/v1/orgs/%s/kkts/%s/storages", c.client.inn, regId)
	if orgType != nil {
		path = fmt.Sprintf("ofd/v1/orgs/%s/kkts/%s/storages?status=%d", c.client.inn, regId, orgType)
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
