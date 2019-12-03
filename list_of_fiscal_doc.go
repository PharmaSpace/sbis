package sbis

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type ListOfFiscalDoc struct {
	client *Client
}

type ListOfFiscalDocResponse struct {
	StorageID     string `json:"storageId"`
	Model         string `json:"model"`
	Status        int    `json:"status"`
	EffectiveFrom string `json:"effectiveFrom"`
	EffectiveTo   string `json:"effectiveTo"`
}

func NewListOfFiscalDoc(client *Client) *ListOfFiscalDoc {
	return &ListOfFiscalDoc{client: client}
}

func (c *ListOfFiscalDoc) Get(storages string, regId string, dateFrom string) ([]*ListOfFiscalDocResponse, error) {
	var ar []*ListOfFiscalDocResponse

	// ofd/v1/orgs/<inn>/kkts/<regId>/storages/<storageId>/docs?dateFrom=<dateFrom>&dateTo=<dateTo>&shiftNumber=<shiftNumber>&startId=<startId>&limit=<limit>

	if len(c.client.inn) == 0 || len(regId) == 0 {
		return nil, ErrEmptyREQ
	}

	path := fmt.Sprintf("ofd/v1/orgs/%s/kkts/%s/storages/%s/docs?dateFrom=%s", c.client.inn, regId, storages, dateFrom)
	req, err := c.client.NewRequest(true, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.client.Do(req)
	if resp != nil {
		data, err := ioutil.ReadAll(resp.Body)

		c.client.logger.Info().Msgf(string(data))

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
