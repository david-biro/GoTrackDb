package influxdb

import (
	"GoTrackDb/model"
	"fmt"
	"log"
	"reflect"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	client "github.com/influxdata/influxdb1-client/v2"
)

// Connect ...
func Connect(host string, port string, user string, password string) client.Client {
	addr := fmt.Sprintf("http://%s:%s", host, port)
	conf := client.HTTPConfig{Addr: addr, Username: user, Password: password, Timeout: model.Timeout * time.Millisecond}

	fmt.Println(addr)

	ic, err := client.NewHTTPClient(conf)
	if err != nil {
		//log.Fatal(err)
		fmt.Println("[WARN] InfluxDB: Creating HTTPClient")
		return nil
	}

	if _, _, err := ic.Ping(1); err != nil {
		//log.Fatal(err)
		fmt.Println("[WARN] InfluxDB: Ping failed")
		return nil
	}
	return ic
}

// CreateNewBatchPoint ...
func CreateNewBatchPoint(ic client.Client, dbname string, precision string, mpoint string, marker model.Marker) {
	defer ic.Close()
	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  dbname,
		Precision: precision,
	})

	tags := map[string]string{"": ""}
	/*fields := map[string]interface{}{
		"idle":   10.1,
		"system": 53.3,
		"user":   46.6,
	}*/
	fields := make(map[string]interface{})

	for _, i := range model.GetInfluxFieldNames() {
		fields[i] = reflect.ValueOf(&marker).Elem().FieldByName(cases.Title(language.Und).String(i)).Interface()
	}
	pt, err := client.NewPoint(mpoint, tags, fields, time.Now())
	if err != nil {
		fmt.Println("[WARN] InfluxDB: ", err.Error())
	}
	bp.AddPoint(pt)

	ic.Write(bp)
	err = ic.Write(bp)
	if err != nil {
		fmt.Println("[WARN] InfluxDB: Insert data error: ", err)
	} else {
		fmt.Println("Data inserted to InfluxDB")
	}
	if err := ic.Close(); err != nil {
		log.Println("[WARN] InfluxDB: Error closing socket: ", err)
	}
}
