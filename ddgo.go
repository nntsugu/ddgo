package main

import (
	"net/http"
	"net/url"
	"path/filepath"
	"io/ioutil"
	"fmt"

	 "gopkg.in/yaml.v2"
)

type DatadogKeys struct {
	Datadog struct {
		Api_key string
		App_key string
	}
}

type DatadogInformation struct {
	AuthenticationEP string
	AuthenticationParams []string
}

var DDKeys DatadogKeys = DatadogKeys{}
var DDInformation DatadogInformation = DatadogInformation{
	AuthenticationEP: "https://app.datadoghq.com/api/v1/validate",
	AuthenticationParams:  []string{"api_key"},
}



func main() {
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