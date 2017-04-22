package main

import (
	//"net/http"
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

var Datadog DatadogInformation = DatadogInformation{
	AuthenticationEP: "https://app.datadoghq.com/api/v1/validate",
	AuthenticationParams:  []string{"?api_key="},

}


func main() {
	targetFile := filepath.Join("..", "secrets", "dd.yaml")

	println(targetFile)

	b, err := ioutil.ReadFile(targetFile)
	if err != nil {
		fmt.Println(err)
	}

	// str := string(b)
	// fmt.Println(str)
	// fmt.Println(Datadog.AuthenticationEP)

	m := make(map[interface{}]interface{})
	err = yaml.Unmarshal(b, &m)
	if err != nil {
		fmt.Println(err)
	}


	// fmt.Println("%s", m["datadog"].(map[interface{}]interface{})["api_key"])
	fmt.Println("%s", m["datadog"])
	var d DatadogKeys
	err = yaml.Unmarshal(b, &d)
	fmt.Println(d.Datadog.Api_key)
}

func doGet() {
	values := url.Values{}
	values.Add("api_key", "")
}