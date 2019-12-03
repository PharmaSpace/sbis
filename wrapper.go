package sbis

import (
	"github.com/rs/zerolog/log"
)

func GetReceipts(inn string, dateFrom, dateTo string, options ...Option) ([]*ListOfFiscalDocResponse, error) {
	var response []*ListOfFiscalDocResponse

	options = append(options, SetInn(inn))
	sbis, err := NewClient(options...)

	if err != nil {
		log.Fatal().Msgf(err.Error())
	}
	if sbis != nil {

		// Список ККТ по организации
		list, err := sbis.ListKKTbyOrganization.Get(nil)
		if err != nil {
			return nil, err
		}
		if len(list) > 0 {
			for _, regValue := range list {

				// Список фискальных накопителей по ККТ
				list, err := sbis.ListOfFiscalDriver.Get(regValue.RegID, nil)
				if err != nil {
					return nil, err
				}
				if len(list) > 0 {
					for _, value := range list {

						// Список фискальных документов по фискальному накопителю
						list, err := sbis.ListOfFiscalDoc.Get(value.StorageID, regValue.RegID, dateFrom, dateTo)
						if err != nil {
							return nil, err
						}

						response = list
					}
				}
			}
		}
	}

	return response, nil
}
