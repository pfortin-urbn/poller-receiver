package receivers

import (
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

/*
UPDATE uc_inventory
    SET val1 = in_val1,
        val2 = in_val2
    WHERE val3 = in_val3;

IF ( sql%rowcount = 0 )
    THEN
    INSERT INTO tablename
        VALUES (in_val1, in_val2, in_val3);
END IF;
*/

type OracleReceiver struct {
	Wg sync.WaitGroup
	Cs chan string
}

func (receiver OracleReceiver) PutMessages() {
	defer receiver.Wg.Done()
	log.Println("Starting Oracle Receiver...")
	for {
		data := <-receiver.Cs
		log.Println("Adding Data to Oracle:", data)
		time.Sleep(100 * time.Millisecond)
	}
}
