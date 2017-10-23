package receivers

import (
	"encoding/json"
	"sync"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	log "github.com/sirupsen/logrus"
)

type MongoInventoryFact struct {
	Brand           string `bson:"brand",json:"brand"`
	DocType         string `bson:"docType",json:"docType"`
	Pool            string `bson:"poolId",json:"pool"`
	ProductID       string `bson:"productId",json:"productId"`
	Availability    string `bson:"availability",json:"availability"`
	BackOrderLevel  int    `bson:"backOrderLevel",json:"backOrderLevel"`
	Backorderable   string `bson:"backorderable",json:"backorderable"`
	ShipmentDate    int    `bson:"shipmentDate",json:"shipmentDate"`
	SiteID          string `bson:"siteId",json:"siteId"`
	SkuID           string `bson:"skuId",json:"skuId"`
	StockLevel      int    `bson:"stockLevel",json:"stockLevel"`
	StoreStockLevel int    `bson:"storeStockLevel",json:"storeStockLevel"`
}

type InventoryFactMessage struct {
	Brand     string `json:"brand"`
	DocType   string `json:"docType"`
	Pool      string `json:"pool"`
	ProductID string `json:"productId"`
	Skus      []struct {
		Availability    string `json:"availability"`
		BackOrderLevel  int    `json:"backOrderLevel"`
		Backorderable   string `json:"backorderable"`
		ShipmentDate    int    `json:"shipmentDate"`
		SiteID          string `json:"siteId"`
		SkuID           string `json:"skuId"`
		StockLevel      int    `json:"stockLevel"`
		StoreStockLevel int    `json:"storeStockLevel"`
	} `json:"skus"`
}

type MongoReceiver struct {
	Wg              *sync.WaitGroup
	Cs              chan string
	Done            chan bool
	MongoServers    string
	MongoDatabase   string
	MongoCollection string
}

func (receiver MongoReceiver) PutMessages() {
	defer receiver.Wg.Done()
	log.Println("Starting Mongo Receiver...")

	var err error
	var mongo_session *mgo.Session
	mongo_session, err = mgo.Dial(receiver.MongoServers)
	if err != nil {
		log.Println("Error: ", err)
		return
	}
	defer mongo_session.Close()

	mongo_session.SetMode(mgo.Monotonic, true)
	c := mongo_session.DB(receiver.MongoDatabase).C(receiver.MongoCollection)

	for {
		select {
		case message := <-receiver.Cs:
			//log.Debug(message)
			var InventoryFact InventoryFactMessage
			err = json.Unmarshal([]byte(message), &InventoryFact)
			if err != nil {
				log.Error("Error: ", err.Error())
			}

			for j := 0; j < len(InventoryFact.Skus); j++ {
				p := MongoInventoryFact{
					DocType:         InventoryFact.DocType,
					Brand:           InventoryFact.Brand,
					Pool:            InventoryFact.Pool,
					ProductID:       InventoryFact.ProductID,
					SkuID:           InventoryFact.Skus[j].SkuID,
					SiteID:          InventoryFact.Skus[j].SiteID,
					StockLevel:      InventoryFact.Skus[j].StockLevel,
					StoreStockLevel: InventoryFact.Skus[j].StoreStockLevel,
					Availability:    InventoryFact.Skus[j].Availability,
					BackOrderLevel:  InventoryFact.Skus[j].BackOrderLevel,
					Backorderable:   InventoryFact.Skus[j].Backorderable,
					ShipmentDate:    InventoryFact.Skus[j].ShipmentDate,
				}

				if invFactChanged(p, c) {
					selector := bson.M{"brand": p.Brand, "poolId": p.Pool, "skuId": p.SkuID}
					upsertdata := bson.M{"$set": p}
					_, err2 := c.Upsert(selector, upsertdata)
					if err2 != nil {
						log.Error("Error: ", err2)
					}
					log.Debugf("Saved Fact %s/%s/%s)", p.Brand, p.Pool, p.SkuID)
				} else {
					log.Debugf("Fact was not changed, not saving it...(%s/%s/%s)", p.Brand, p.Pool, p.SkuID)
				}
			}

		case <-receiver.Done:
			log.Println("Receiver Done received.")
			return
		}
	}
}

func invFactChanged(invFact MongoInventoryFact, coll *mgo.Collection) bool {
	oldFact := MongoInventoryFact{}
	selector := bson.M{"brand": invFact.Brand, "poolId": invFact.Pool, "skuId": invFact.SkuID}
	coll.Find(selector).One(&oldFact)
	if oldFact.StockLevel != invFact.StockLevel || oldFact.BackOrderLevel != invFact.BackOrderLevel ||
		oldFact.ShipmentDate != invFact.ShipmentDate {
		return true
	}
	return false
}
