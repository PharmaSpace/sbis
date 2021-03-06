package sbis

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"
)

type ListOfFiscalDoc struct {
	client *Client
}

type Item struct {
	Quantity    float64 `json:"quantity"`
	Name        string  `json:"name"`
	Sum         int     `json:"sum"`
	Price       int     `json:"price"`
	Nds         int     `json:"nds"`
	NdsSum      int     `json:"ndsSum"`
	ProductType int     `json:"productType"`
	PaymentType int     `json:"paymentType"`
}

type Receipt struct {
	ReceiptCode             int     `json:"receiptCode"`
	RawData                 string  `json:"rawData"`
	ReceiveDateTime         string  `json:"receiveDateTime"`
	SendDateTime            string  `json:"sendDateTime"`
	DateTime                int     `json:"dateTime"`
	UserInn                 string  `json:"userInn"`
	TotalSum                int     `json:"totalSum"`
	Operator                string  `json:"operator"`
	CashTotalSum            int     `json:"cashTotalSum"`
	KktRegID                string  `json:"kktRegId"`
	ShiftNumber             int     `json:"shiftNumber"`
	FiscalDocumentNumber    int     `json:"fiscalDocumentNumber"`
	FiscalDriveNumber       string  `json:"fiscalDriveNumber"`
	RequestNumber           int     `json:"requestNumber"`
	OperationType           int     `json:"operationType"`
	TaxationType            int     `json:"taxationType"`
	Items                   []*Item `json:"items"`
	FiscalSign              int64   `json:"fiscalSign"`
	EcashTotalSum           int     `json:"ecashTotalSum"`
	NdsNo                   int     `json:"ndsNo"`
	FiscalDocumentFormatVer int     `json:"fiscalDocumentFormatVer"`
	PrepaidSum              int     `json:"prepaidSum"`
	CreditSum               int     `json:"creditSum"`
	ProvisionSum            int     `json:"provisionSum"`
	Url                     string  `json:"url"`
}

type ListOfFiscalDocResponse struct {
	Receipt Receipt `json:"receipt"`
}

func NewListOfFiscalDoc(client *Client) *ListOfFiscalDoc {
	return &ListOfFiscalDoc{client: client}
}

func (r *Receipt) GetUrl() {
	tm := time.Unix(int64(r.DateTime), 0)
	r.Url = "https://ofd.sbis.ru/rec/" + r.KktRegID + "/" + tm.Format("020106") + "/" + strconv.Itoa(int(r.FiscalSign))
}

func (c *ListOfFiscalDoc) Get(storages string, regId string, dateFrom, dateTo string) ([]*ListOfFiscalDocResponse, error) {
	var ar []*ListOfFiscalDocResponse

	// ofd/v1/orgs/<inn>/kkts/<regId>/storages/<storageId>/docs?dateFrom=<dateFrom>&dateTo=<dateTo>&shiftNumber=<shiftNumber>&startId=<startId>&limit=<limit>

	if len(c.client.inn) == 0 || len(regId) == 0 {
		return nil, ErrEmptyREQ
	}

	path := fmt.Sprintf("ofd/v1/orgs/%s/kkts/%s/storages/%s/docs?dateFrom=%s&dateTo=%s&limit=1000", c.client.inn, regId, storages, dateFrom, dateTo)
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
		for _, item := range ar {
			item.Receipt.GetUrl()
		}
	}

	if c.client.verbose {
		c.client.logger.Info().Msgf("%s - success", path)
	}
	return ar, nil

}
