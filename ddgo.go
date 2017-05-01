package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const (
	Version = "0.0.1"
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

type DatadogInformation struct {
	AuthenticationEP     string
	AuthenticationParams []string
}

var Arguments Argument = Argument{}
var DDKeys DatadogKeys = DatadogKeys{}
var DDInformation DatadogInformation = DatadogInformation{
	AuthenticationEP:     "https://app.datadoghq.com/api/v1/validate",
	AuthenticationParams: []string{"api_key"},
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
		fmt.Println(b)
		return
	}

	//================
	// Load seacrets
	// ToDo How to manage secrets?
	targetFile := filepath.Join("..", "secrets", "dd.yaml")

	b, err := ioutil.ReadFile(targetFile)
	if err != nil {
		fmt.Println(err)
	}

	err = yaml.Unmarshal(b, &DDKeys)
	// fmt.Println(DDKeys.Datadog.Api_key)
	//================

	doGet()
}

func doGet() {
	client := &http.Client{}
	values := url.Values{}
	values.Add("api_key", DDKeys.Datadog.Api_key)

	req, err := http.NewRequest("GET", DDInformation.AuthenticationEP, nil)
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
}
