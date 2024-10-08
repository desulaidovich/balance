# balance

### Запросы
| Запрос          | Метод | Параметры                                     |
|-----------------|-------|-----------------------------------------------|
| /wallet/create  | POST  | money=INT_VALUE&level=INT_VALUE               |
| /wallet/hold    | POST  | wallet_id=INT_VALUE&money=INT_VALUE           |
| /wallet/dishold | POST  | wallet_id=INT_VALUE&money=INT_VALUE           |
| /wallet/edit    | POST  | wallet_id=INT_VALUE&money=INT_VALUE&type_id=INT_VALUE |
| /wallet/get     | GET   | wallet_id=INT_VALUE                           |

### type_id
| ID | Назначение |
|----|------------|
| 1  | Списание   |
| 2  | Пополнение |


### NATS example
```go
package main

import (
	"log"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

const (
	REQUEST_ERROR_CODE    = 500
	REQUEST_NO_ERROR_CODE = 0
)

type JSONMessage struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Data    *Data  `json:"data,omitempty"`
}

type Data struct {
	WalletData *WalletData `json:"wallet_data,omitempty"`
	WalletID   int         `json:"wallet_id,omitempty"`
}

type WalletData struct {
	CreateAt       string          `json:"created_at,omitempty"`
	UpdatedAt      string          `json:"updated_at,omitempty"`
	Hold           int             `json:"hold,omitempty"`
	Balance        int             `json:"balance,omitempty"`
	Identification *Identification `json:"identification,omitempty"`
}

type Identification struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

func main() {
	nc, err := nats.Connect(nats.DefaultURL, nats.Name("Balance"), nats.ReconnectWait(40*time.Second),
		nats.ErrorHandler(func(_ *nats.Conn, _ *nats.Subscription, err error) {
			log.Printf("Error: %v", err)
		}))
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	getStatusTxt := func(nc *nats.Conn) string {
		switch nc.Status() {
		case nats.CONNECTED:
			return "Connected"
		case nats.CLOSED:
			return "Closed"
		default:
			return "Other"
		}
	}
	log.Printf("The connection is %v\n", getStatusTxt(nc))

	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		log.Fatal(err)
	}
	defer ec.Close()

	wg := sync.WaitGroup{}
	wg.Add(1)

	if _, err := ec.Subscribe("created", func(data *JSONMessage) {
		log.Printf("[created]\tcode:%d, message:%s ", data.Code, data.Message)
		log.Printf("\t\twallet id:%d, balance:%d hold:%d",
			data.Data.WalletID, data.Data.WalletData.Balance, data.Data.WalletData.Hold)
		log.Printf("\t\tidentification %d (%s)",
			data.Data.WalletData.Identification.ID, data.Data.WalletData.Identification.Name)
		wg.Done()
	}); err != nil {
		log.Fatal(err)
	}

	if _, err := ec.Subscribe("holded", func(data *JSONMessage) {
		log.Printf("[holded]\tcode:%d, message:%s ", data.Code, data.Message)
		log.Printf("\t\twallet id:%d, balance:%d hold:%d",
			data.Data.WalletID, data.Data.WalletData.Balance, data.Data.WalletData.Hold)
		wg.Done()
	}); err != nil {
		log.Fatal(err)
	}

	if _, err := ec.Subscribe("disholded", func(data *JSONMessage) {
		log.Printf("[disholded]\tcode:%d, message:%s ", data.Code, data.Message)
		log.Printf("\t\twallet id:%d, balance:%d hold:%d",
			data.Data.WalletID, data.Data.WalletData.Balance, data.Data.WalletData.Hold)
		wg.Done()
	}); err != nil {
		log.Fatal(err)
	}

	if _, err := ec.Subscribe("edited", func(data *JSONMessage) {
		log.Printf("[edited]\tcode:%d, message:%s ", data.Code, data.Message)
		log.Printf("\t\twallet id:%d, balance:%d hold:%d",
			data.Data.WalletID, data.Data.WalletData.Balance, data.Data.WalletData.Hold)
		log.Printf("\t\tidentification %d (%s)",
			data.Data.WalletData.Identification.ID, data.Data.WalletData.Identification.Name)
		wg.Done()
	}); err != nil {
		log.Fatal(err)
	}

	if _, err := ec.Subscribe("got", func(data *JSONMessage) {
		log.Printf("[edited]\tcode:%d, message:%s ", data.Code, data.Message)
		log.Printf("\t\twallet id:%d, balance:%d hold:%d",
			data.Data.WalletID, data.Data.WalletData.Balance, data.Data.WalletData.Hold)
		log.Printf("\t\tidentification %d (%s)",
			data.Data.WalletData.Identification.ID, data.Data.WalletData.Identification.Name)
		wg.Done()
	}); err != nil {
		log.Fatal(err)
	}

	wg.Wait()
}
```


###  [Thunder Client](https://www.thunderclient.com/)
```json
{
    "clientName": "Thunder Client",
    "collectionName": "balance",
    "collectionId": "a93bdbd1-523d-4269-90fd-4f73e59bee58",
    "dateExported": "2024-07-25T09:56:22.553Z",
    "version": "1.2",
    "folders": [],
    "requests": [
        {
            "_id": "2b2db4b0-cbf4-426d-9688-37f92e0e0f05",
            "colId": "a93bdbd1-523d-4269-90fd-4f73e59bee58",
            "containerId": "",
            "name": "Create",
            "url": "http://localhost:8080/wallet/create?money=14999&level=1",
            "method": "POST",
            "sortNum": 10000,
            "created": "2024-07-22T14:03:07.079Z",
            "modified": "2024-07-25T09:54:52.442Z",
            "headers": [],
            "params": [
                {
                    "name": "money",
                    "value": "14999",
                    "isPath": false
                },
                {
                    "name": "level",
                    "value": "1",
                    "isPath": false
                }
            ],
            "tests": [
                {
                    "type": "res-body",
                    "custom": "",
                    "action": "isjson",
                    "value": ""
                }
            ],
            "preReq": {
                "options": {
                    "clearCookies": true
                }
            }
        },
        {
            "_id": "5aad003c-9ca4-4eb4-b6bb-5a1798db42be",
            "colId": "a93bdbd1-523d-4269-90fd-4f73e59bee58",
            "containerId": "",
            "name": "Hold",
            "url": "http://localhost:8080/wallet/hold?wallet_id=7&money=1",
            "method": "POST",
            "sortNum": 20000,
            "created": "2024-07-24T12:05:31.488Z",
            "modified": "2024-07-25T09:56:05.175Z",
            "headers": [],
            "params": [
                {
                    "name": "wallet_id",
                    "value": "7",
                    "isPath": false
                },
                {
                    "name": "money",
                    "value": "1",
                    "isPath": false
                }
            ],
            "tests": [
                {
                    "type": "res-body",
                    "custom": "",
                    "action": "isjson",
                    "value": ""
                }
            ]
        },
        {
            "_id": "e22a41ab-b360-430d-a8bf-e11117bf132f",
            "colId": "a93bdbd1-523d-4269-90fd-4f73e59bee58",
            "containerId": "",
            "name": "Dishold",
            "url": "http://localhost:8080/wallet/dishold?wallet_id=6AS&money=-1",
            "method": "POST",
            "sortNum": 30000,
            "created": "2024-07-24T13:04:34.694Z",
            "modified": "2024-07-25T08:14:58.878Z",
            "headers": [],
            "params": [
                {
                    "name": "wallet_id",
                    "value": "6AS",
                    "isPath": false
                },
                {
                    "name": "money",
                    "value": "-1",
                    "isPath": false
                }
            ]
        },
        {
            "_id": "af2128ba-5a48-440b-9be8-7f5ffaf86776",
            "colId": "a93bdbd1-523d-4269-90fd-4f73e59bee58",
            "containerId": "",
            "name": "Edit",
            "url": "http://localhost:8080/wallet/edit?wallet_id=7&money=0&type_i=3",
            "method": "POST",
            "sortNum": 40000,
            "created": "2024-07-24T13:56:14.540Z",
            "modified": "2024-07-25T09:35:49.580Z",
            "headers": [],
            "params": [
                {
                    "name": "wallet_id",
                    "value": "7",
                    "isPath": false
                },
                {
                    "name": "money",
                    "value": "0",
                    "isPath": false
                },
                {
                    "name": "type_i",
                    "value": "3",
                    "isPath": false
                }
            ]
        },
        {
            "_id": "a940abd2-5f27-48d9-99c2-a27bb1a11070",
            "colId": "a93bdbd1-523d-4269-90fd-4f73e59bee58",
            "containerId": "",
            "name": "Get",
            "url": "http://localhost:8080/wallet/get?wallet_id=1",
            "method": "GET",
            "sortNum": 50000,
            "created": "2024-07-25T06:59:10.312Z",
            "modified": "2024-07-25T09:22:32.338Z",
            "headers": [],
            "params": [
                {
                    "name": "wallet_id",
                    "value": "1",
                    "isPath": false
                }
            ]
        }
    ],
    "ref": "IOxQSEqg9JTiAEtwatIeI7tKyPCs-eCxEmYXHMcDmt3KsS-44HLsimfQ0tAO6Qm_g7mocuyY-1gkKLcE5cxzXg"
}
```