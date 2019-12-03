package sbis

import (
	"github.com/rs/zerolog/log"
)

func GetReceipts(inn string, dateFrom string) {

	sbis, err := NewClient(Verbose(), SetInn(inn), SetAuthConfig(&AuthConfig{
		AppClientID: "1025293145607151",
		Login:       "PharmaSpace",
		Password:    "vfTrYE86$",
	}))

	if err != nil {
		log.Fatal().Msgf(err.Error())
	}
	if sbis != nil {

		// Список ККТ по организации
		list, err := sbis.ListKKTbyOrganization.Get(nil)
		if err != nil {
			log.Fatal().Msgf(err.Error())
		}
		if len(list) > 0 {
			for _, regValue := range list {

				// Список фискальных накопителей по ККТ
				list, err := sbis.ListOfFiscalDriver.Get(regValue.RegID, nil)
				if err != nil {
					log.Fatal().Msgf(err.Error())
				}
				if len(list) > 0 {
					for _, value := range list {

						// Список фискальных документов по фискальному накопителю
						list, err := sbis.ListOfFiscalDoc.Get(value.StorageID, regValue.RegID, dateFrom)
						if err != nil {
							log.Fatal().Msgf(err.Error())
						}

						log.Print(list)

					}
				}
			}
		}
	}
}
