package main

import (
	"GoTrackDb/config"
	"GoTrackDb/influxdb"
	"GoTrackDb/model"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"
	"path/filepath"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	client "github.com/influxdata/influxdb1-client/v2"
)

// var ic *client.Client
var ms = model.Markers{}
var xmlstruct, filename string
var mutex = &sync.Mutex{}

func setEnv() {
	fmt.Println("OS type: ", runtime.GOOS)

	switch runtime.GOOS {
	case "linux":
		filename = config.Cfg.DataXMLPathLinux
	case "darwin", "windows", "nacl":
		filename = config.Cfg.DataXMLPathOther
	default:
		filename = "data.xml"
	}
}

func readConf() {
	exePath, err := os.Executable()
	if err != nil {
			log.Fatal(err)
	}

	exeDir := filepath.Dir(exePath)
	configPath := filepath.Join(exeDir, "config.file")

	config.Cfg, err = config.ReadConfig(configPath)
	if err != nil {
		log.Fatalf("Error reading config file at %s: %w", configPath, err)
	}
}

func markerCleanup(mPtr *map[string]model.Marker) {
	for {
		fmt.Printf("Marker housekeeping started\n")
		n := time.Now()
		for k, v := range *mPtr {
			lastSeen := time.Unix(v.UnixTime, 0)
			if d := math.Abs(n.Sub(lastSeen).Hours() / 24); d >= 7 {
				mutex.Lock()
				delete(*mPtr, k)
				mutex.Unlock()
				fmt.Printf("key: %s, car: %s, last seen: %v days ago\n", k, v.Car, math.Round(d))
			}
		}
		fmt.Printf("Marker housekeeping done %s\n", time.Now())
		time.Sleep(60 * time.Second)
	}
}

func writeXmlData(fPtr *string, data string) error {
	file, err := os.Create(*fPtr)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.WriteString(file, data)
	if err != nil {
		return err
	}
	return file.Sync()
}

func handler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var n = model.Marker{}
	var innerXML = ""
	for key, value := range r.Form {
		fmt.Printf("%s: %s, ", key, value)
		field := reflect.ValueOf(&n).Elem().FieldByName(cases.Title(language.Und).String(key))
		if field.IsValid() {
			field.SetString(value[0])
		} else {
			fmt.Printf("Error %s", r.URL)
		}
	}
	fmt.Println("\n")
	n.UnixTime = time.Now().Unix()
	n.Time = time.Now().Format("2006-01-02 15:04:05")
	if r.FormValue("sn") != "" && r.FormValue("lat") != "0.0" && r.FormValue("lon") != "0.0" {
		mutex.Lock()

		//keys for sorting hashmap
		ms.M[r.FormValue("sn")] = n
		keys := make([]string, 0)
		for k := range ms.M {
			keys = append(keys, k)
		}

		//sort keys
		sort.Strings(keys)
		xmlstruct = string(model.Header + model.Markersbegin)

		//write sorted
		for _, k := range keys {
			v := ms.M[k]
			markerstruct := reflect.ValueOf(v)
			fieldtype := markerstruct.Type()
			innerXML += "  <marker "
			for i := 0; i < markerstruct.NumField(); i++ {
				innerXML += fmt.Sprintf("%s=\"%v\" ", strings.ToLower(fieldtype.Field(i).Name), markerstruct.Field(i).Interface())
			}
			innerXML += "/>\n"
		}
		xmlstruct += innerXML
		mutex.Unlock()
		xmlstruct += string(model.Markersend)
		writeXmlData(&filename, xmlstruct)

		//writedb
		ic := influxdb.Connect(model.Dbhost, model.Dbport, model.Dbuser, model.Dbpass)
		if ic != nil {
			influxdb.CreateNewBatchPoint(ic, model.Dbname, model.Precision, model.MeasurementPoint, n)
		}
	}
	w.WriteHeader(200)
	fmt.Fprintf(w, n.Time)
}

func initialize() {
	readConf()
	setEnv()
	model.SetDbParams()
	ic := influxdb.Connect(model.Dbhost, model.Dbport, model.Dbuser, model.Dbpass)
	fmt.Println(model.Dbuser)
	fmt.Println(model.Dbhost)
	if ic != nil {
		defer ic.Close()
		q := client.Query{
			Command:  fmt.Sprintf("create database %s", model.Dbname),
			Database: model.Dbname,
		}
		ic.Query(q)
	}
	ms.M = make(map[string]model.Marker)
	go markerCleanup(&ms.M)
}

func main() {
	initialize()
	http.HandleFunc(config.Cfg.Pattern, handler)
	fmt.Println(config.Cfg.Pattern)
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {})
	log.Fatal(http.ListenAndServe(":"+config.Cfg.Port, nil))
}
