package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"time"

	dproxy "github.com/koron/go-dproxy"
)

const (
	Version = "0.0.1"
)

// const for debug
const (
	LogSeparator = "===================="
)

type Argument struct {
	configFilePath         string
	moniotoringSettingsDir string
}

type DatadogKeys struct {
	Datadog struct {
		Api_key string
		App_key string
	}
}

type Eps struct {
	End_point string
	// Params    []string
}
type DatadogInformation struct {
	Authentication         Eps
	GetAllMonitorDetails   Eps
	CreateAMonitor         Eps
	GetAllMonitoriDowntime Eps
}

var Arguments Argument = Argument{}
var DDKeys DatadogKeys = DatadogKeys{}
var DDInformation = NewDatadogInformation()

func NewDatadogInformation() *DatadogInformation {
	return &DatadogInformation{
		Authentication: Eps{
			End_point: "https://app.datadoghq.com/api/v1/validate",
			// Params:    []string{"api_key"},
		},
		GetAllMonitorDetails: Eps{
			End_point: "https://app.datadoghq.com/api/v1/monitor",
		},
		// http://docs.datadoghq.com/ja/api/?lang=console#monitor-create
		CreateAMonitor: Eps{
			End_point: "https://app.datadoghq.com/api/v1/monitor",
		},
		// http://docs.datadoghq.com/ja/api/#downtimes
		GetAllMonitoriDowntime: Eps{
			End_point: "https://app.datadoghq.com/api/v1/downtime",
		},
	}
}

func main() {
	invalidArgments := false
	var showVersion bool
	var showUsage bool

	// -v -version
	flag.BoolVar(&showVersion, "v", false, "show version")
	flag.BoolVar(&showVersion, "version", false, "show version")
	// -h -help
	flag.BoolVar(&showUsage, "h", false, "show usage")
	flag.BoolVar(&showUsage, "help", false, "show usage")
	// -f (required)
	flag.StringVar(&Arguments.configFilePath, "f", "", "set credential file path which must have api_key and app_key(application_key) to access Datadog API. ref. http://docs.datadoghq.com/api/")
	// -m (required)
	flag.StringVar(&Arguments.moniotoringSettingsDir, "m", "", "set the directory path which has monitoring definitions. e.g) ~/monitorring_setting.d")

	flag.Parse()
	if showUsage {
		flag.Usage()
		return
	}
	if showVersion {
		fmt.Println(Version)
		return
	}
	// validation
	if Arguments.configFilePath == "" {
		fmt.Println("-f is required.")
		invalidArgments = true
	}
	if Arguments.moniotoringSettingsDir == "" {
		fmt.Println("-m is required.")
		invalidArgments = true
	}
	if invalidArgments {
		fmt.Println("-h : show usage")
		return
	}

	// initialize
	if Arguments.configFilePath != "" {
		//================
		// Load seacrets
		b, err := ioutil.ReadFile(Arguments.configFilePath)
		if err != nil {
			fmt.Println(err)
		}

		err = yaml.Unmarshal(b, &DDKeys)
		// fmt.Println(DDKeys.Datadog.Api_key)
		//================
	}

	isThereAnyLongDowntime()
	// createMonitors()
}

func createMonitors() {
	var monitors []interface{}
	var conf interface{}
	var skip bool

	client := &http.Client{}
	values := url.Values{}
	values.Add("api_key", DDKeys.Datadog.Api_key)
	values.Add("application_key", DDKeys.Datadog.App_key)

	// u, err := user.Current()
	// if err != nil {
	// 	log.Fatal(err)
	// 	return
	// }
	//confPath := filepath.Join(u.HomeDir, "conf.d", "monitor_template.d"
	confPath := Arguments.moniotoringSettingsDir

	d, err := ioutil.ReadDir(confPath)
	if err != nil {
		fmt.Println("Couldn't read JSON file directory:", err)
		return
	}

	for _, f := range d {
		fmt.Println(f.Name())
		jsonBytes, err := ioutil.ReadFile(filepath.Join(confPath, f.Name()))
		if err != nil {
			fmt.Println("Couldn't read JSON file for monitoring setting:", err)
			return
		}
		monitors = append(monitors, jsonBytes)
	}

	createTargets := _duplicationCheck(monitors)
	for _, f := range monitors {
		json.Unmarshal(f.([]byte), &conf)
		monitorName, err := dproxy.New(conf).M("name").String()
		if err != nil {
			log.Fatal("Monitor name is undefined please check json file(s) ", err)
		}
		skip = true
		for _, target := range createTargets {
			if target == monitorName {
				skip = false
				break
			}
		}
		if skip == true {
			fmt.Println("skip:", monitorName)
			continue
		}

		req, err := http.NewRequest("POST", DDInformation.GetAllMonitorDetails.End_point, bytes.NewBuffer(f.([]byte)))
		if err != nil {
			fmt.Println(err)
			return
		}
		req.Header.Set("Content-Type", "application/json")

		req.URL.RawQuery = values.Encode()
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("error:", err)
			log.Println("response body:", string(b))
			return
		}
		// log.Println(LogSeparator, "createAMonitor", LogSeparator)
		fmt.Println("created:", monitorName)
		// log.Println(string(b))
		// log.Println(LogSeparator, "createAMonitor", LogSeparator)
	}
	return
}

// return the slice which has the list of create target moniroting settings.
func _duplicationCheck(monitors []interface{}) []string {
	var conf interface{}
	var createTargets []string
	for _, monitor := range monitors {
		json.Unmarshal(monitor.([]byte), &conf)
		monitorName, err := dproxy.New(conf).M("name").String()
		if err != nil {
			log.Fatal("Monitor name is undefined please check json file(s) ", err)
		}
		if _isMonitorExists(monitorName) {
			log.Println(monitorName, "is already exists.")
		} else {
			createTargets = append(createTargets, monitorName)
		}
	}
	return createTargets
}

func _isMonitorExists(name string) bool {
	var conf interface{}

	client := &http.Client{}
	values := url.Values{}
	values.Add("api_key", DDKeys.Datadog.Api_key)
	values.Add("application_key", DDKeys.Datadog.App_key)
	values.Add("name", name)

	req, err := http.NewRequest("GET", DDInformation.GetAllMonitorDetails.End_point, nil)
	if err != nil {
		fmt.Println(err)
		return false
	}

	req.URL.RawQuery = values.Encode()
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return false
	}
	json.Unmarshal(b, &conf)

	monitorName, err := dproxy.New(conf).A(0).M("name").String()
	if err != nil {
		// log.Println("Monitoring setting : ", name, " is not found on Datadog")
		return false
	}
	if name == monitorName {
		// log.Println(name, "is already exist.")
		return true
	}
	return false
}

func isThereAnyLongDowntime() {
	// 1 week lator
	t1 := time.Now().AddDate(0, 0, 7)
	t2 := t1
	fmt.Println(t2.Unix)
	fmt.Println(t2.Sub(t1)) // 12h0m0s

	loc, _ := time.LoadLocation("Asia/Tokyo")

	t3 := time.Date(2017, 6, 19, 8, 0, 0, 0, loc)
	fmt.Println(t2.Sub(t3)) // 12h0m0s

	// d := time.Since(t3)
	// fmt.Println(d) // 72h0m0s (Go Playgroundで実行した場合)

	var conf interface{}

	client := &http.Client{}
	values := url.Values{}
	values.Add("api_key", DDKeys.Datadog.Api_key)
	values.Add("application_key", DDKeys.Datadog.App_key)

	req, err := http.NewRequest("GET", DDInformation.GetAllMonitoriDowntime.End_point, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	req.URL.RawQuery = values.Encode()
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	json.Unmarshal(b, &conf)

	fmt.Println(conf)

	// monitorName, err := dproxy.New(conf).A(0).M("name").String()

	// downtime, err := dproxy.NewSet(conf).Len()
	// fmt.Println(downtime)
	var endDowntime float64
	var item dproxy.Proxy
	for index := 0; ; {
		fmt.Println(LogSeparator)
		fmt.Println(dproxy.New(conf).A(index))
		item = dproxy.New(conf).A(index)
		index++

		if item.M("end").Nil() == true {
			log.Println("end duration is not set")
			continue
		}
		endDowntime, err = item.M("end").Float64()
		if err != nil {
			fmt.Println(err)
			log.Println(err)
			break
		}
		// fmt.Println(string(endDowntime))
		fmt.Println(endDowntime)
	}
	// for _, downtime := range dproxy.New(conf).A {
	// 	fmt.Println(downtime)
	// }

	// if err != nil {
	// 	// log.Println("Monitoring setting : ", name, " is not found on Datadog")
	// 	return false
	// }
	// if name == monitorName {
	// 	// log.Println(name, "is already exist.")
	// 	return true
	// }
	return

}
