package sbis

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/rs/zerolog/log"
)

type ListKKTbyOrganization struct {
	client *Client
}

type ListKKTResponse struct {
	RegID              string `json:"regId"`
	Model              string `json:"model"`
	FactoryID          string `json:"factoryId"`
	FsNumber           string `json:"fsNumber"`
	KktSalesPointSPPID int    `json:"kktSalesPointSPPId"`
	FirstShiftDate     string `json:"firstShiftDate"`
	Address            string `json:"address"`
	KktSalesPoint      string `json:"kktSalesPoint"`
	Status             int    `json:"status"`
	Kpp                string `json:"kpp"`
	OrganizationName   string `json:"organizationName"`
	FsFinishDate       string `json:"fsFinishDate"`
	LicenseStartDate   string `json:"licenseStartDate"`
	LicenseFinishDate  string `json:"licenseFinishDate"`
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
		return nil, ErrEmptyREQ
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
