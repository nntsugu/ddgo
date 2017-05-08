package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"gopkg.in/yaml.v2"
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
	EndPoint       string
	ResourceParams []string
	QueryParams    []string
}
type DatadogInformation struct {
	Authentication       Eps
	GetAllMonitorDetails Eps
	GetUser              Eps
	GetMonitor           Eps
}

var Arguments Argument = Argument{}
var DDKeys DatadogKeys = DatadogKeys{}
var DDInformation = NewDatadogInformation()

func NewDatadogInformation() *DatadogInformation {
	return &DatadogInformation{
		Authentication: Eps{
			EndPoint:       "https://app.datadoghq.com/api/v1/validate",
			ResourceParams: []string{""},
			QueryParams:    []string{"api_key"},
		},
		GetAllMonitorDetails: Eps{
			EndPoint:       "https://app.datadoghq.com/api/v1/monitor",
			ResourceParams: []string{""},
			QueryParams:    []string{"api_key", "application_key", "from"},
		},
		GetUser: Eps{
			EndPoint:       "https://app.datadoghq.com/api/v1/user",
			ResourceParams: []string{"user"},
			QueryParams:    []string{"api_key", "application_key"},
		},
		GetMonitor: Eps{
			EndPoint:       "https://app.datadoghq.com/api/v1/monitor",
			ResourceParams: []string{"monitor_id"},
			QueryParams:    []string{"api_key", "application_key"},
		},
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
		b, err := ioutil.ReadFile(Arguments.configFilePath)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(b))
		// return
	}

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

	//doGet()
	//getAllMonitorDetails()
	getAllUser()
	getMonitor(1523957)
}

func doGet() {
	client := &http.Client{}
	values := url.Values{}
	values.Add("api_key", DDKeys.Datadog.Api_key)

	req, err := http.NewRequest("GET", DDInformation.Authentication.EndPoint, nil)
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

	fmt.Println(req.URL.RawQuery)
	fmt.Println(resp)

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Println("doGet: ", string(b))
}

func getAllMonitorDetails() {
	client := &http.Client{}
	values := url.Values{}
	values.Add("api_key", DDKeys.Datadog.Api_key)
	values.Add("application_key", DDKeys.Datadog.App_key)

	req, err := http.NewRequest("GET", DDInformation.GetAllMonitorDetails.EndPoint, nil)
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

	fmt.Println(req.URL.RawQuery)
	fmt.Println(resp)

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Println(LogSeparator, "getAllMonitorDetails", LogSeparator)
	log.Println(string(b))
	log.Println(LogSeparator, "getAllMonitorDetails", LogSeparator)
}

func getAllUser() {
	client := &http.Client{}
	values := url.Values{}
	values.Add("api_key", DDKeys.Datadog.Api_key)
	values.Add("application_key", DDKeys.Datadog.App_key)

	req, err := http.NewRequest("GET", DDInformation.GetUser.EndPoint, nil)
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

	fmt.Println(req.URL.RawQuery)
	fmt.Println(resp)

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Println(LogSeparator, "getAllUser", LogSeparator)
	log.Println(string(b))
	log.Println(LogSeparator, "getAllUser", LogSeparator)

	content := []byte(b)
	ioutil.WriteFile("/tmp/gogoogooo", content, os.ModePerm)
}

func getMonitor(monitorId int) {
	client := &http.Client{}
	values := url.Values{}
	values.Add("api_key", DDKeys.Datadog.Api_key)
	values.Add("application_key", DDKeys.Datadog.App_key)

	log.Println(monitorId)
	req, err := http.NewRequest("GET", DDInformation.GetMonitor.EndPoint, nil)
	// req, err := http.NewRequest("GET", DDInformation.GetMonitor.EndPoint+"/"+string(monitorId), nil)
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

	fmt.Println(req.URL.RawQuery)
	fmt.Println(resp)

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Println(LogSeparator, "getMonitor", LogSeparator)
	log.Println(string(b))
	log.Println(LogSeparator, "getMonitor", LogSeparator)

	content := []byte(b)
	ioutil.WriteFile("/tmp/gogoogooo", content, os.ModePerm)
}
