package main

import (
	//"net/http"
	"net/url"
	"path/filepath"
	"io/ioutil"
	"fmt"
)

func main() {
	targetFile := filepath.Join("..", "secrets", "dd.yaml")

	println(targetFile)

	b, err := ioutil.ReadFile(targetFile)
	if err != nil {
		fmt.Println(err)
	}

	str := string(b)
	fmt.Println(str)
	//file, _ := os.Open(targetFile)
	//defer file.Close()
}

func doGet() {
	values := url.Values{}
	values.Add("api_key", "")
}