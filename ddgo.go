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
	"os/user"
	"path/filepath"

	"gopkg.in/yaml.v2"

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
	configFilePath string
}

type DatadogKeys struct {
	Datadog struct {
		Api_key string
		App_key string
	}
}

type Eps struct {
	End_point string
	Params    []string
}
type DatadogInformation struct {
	Authentication       Eps
	GetAllMonitorDetails Eps
	CreateAMonitor       Eps
}

var Arguments Argument = Argument{}
var DDKeys DatadogKeys = DatadogKeys{}
var DDInformation = NewDatadogInformation()

func NewDatadogInformation() *DatadogInformation {
	return &DatadogInformation{
		Authentication: Eps{
			End_point: "https://app.datadoghq.com/api/v1/validate",
			Params:    []string{"api_key"},
		},
		GetAllMonitorDetails: Eps{
			End_point: "https://app.datadoghq.com/api/v1/monitor",
			Params:    []string{"api_key", "application_key", "from"},
		},
		// http://docs.datadoghq.com/ja/api/?lang=console#monitor-create
		CreateAMonitor: Eps{
			End_point: "https://app.datadoghq.com/api/v1/monitor",
			Params:    []string{"type", "query", "name", "message"},
		},
	}
}

type DatadogMonitor struct {
	Type    string `json:"type"`
	Query   string `json:"query"`
	Name    string `json:"name"`
	Message string `json:"message"`
}

func NewDatadogMonitor() *DatadogMonitor {
	return &DatadogMonitor{
		Type:    "",
		Query:   "",
		Name:    "",
		Message: "",
	}
}

func main() {
	var showVersion bool

	// -v -version
	flag.BoolVar(&showVersion, "v", false, "show version")
	flag.BoolVar(&showVersion, "version", false, "show version")
	// -f
	flag.StringVar(&Arguments.configFilePath, "f", "", "set configration file path")

	flag.Parse()
	if showVersion {
		fmt.Println(Version)
		return
	}
	if Arguments.configFilePath != "" {
		//================
		// Load seacrets
		// ToDo How to manage secrets?
		// targetFile := filepath.Join("..", "secrets", "dd.yaml")

		// b, err := ioutil.ReadFile(targetFile)
		b, err := ioutil.ReadFile(Arguments.configFilePath)
		if err != nil {
			fmt.Println(err)
		}

		err = yaml.Unmarshal(b, &DDKeys)
		// fmt.Println(DDKeys.Datadog.Api_key)
		//================
	}

	createMonitors()
}

func createMonitors() {
	// var jsonBytes []byte
	var monitors []interface{}
	var conf interface{}
	var skip bool

	client := &http.Client{}
	values := url.Values{}
	values.Add("api_key", DDKeys.Datadog.Api_key)
	values.Add("application_key", DDKeys.Datadog.App_key)

	// var conf interface{}
	u, err := user.Current()
	if err != nil {
		log.Fatal(err)
		return
	}
	confPath := filepath.Join(u.HomeDir, "conf.d", "monitor_template.d")

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
	log.Println(createTargets)
	for _, f := range monitors {
		// for _, t := range createTargets {
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

		// log.Println(string(ioutil.ReadAll(resp.Header)))

	}
	return
}

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
			log.Println(monitorName, "is already exist.")
		} else {
			createTargets = append(createTargets, monitorName)
		}
	}
	log.Println(createTargets)
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
	// log.Println("monitorName:", monitorName)
	// log.Println("Param:", name)
	if err != nil {
		log.Println("Monitoring setting : ", name, " is not found on Datadog")
		// monitor : name is not found on Datadog
		return false
	}
	if name == monitorName {
		// log.Println(name, "is already exist.")
		return true
	}
	return false
}
