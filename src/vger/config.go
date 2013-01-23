package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type configuration map[string]string

var config configuration

func readJson(path string, object interface{}) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("read json:%s\n%s\n", path, string(data))

	json.Unmarshal(data, object)
}

func readConfig() configuration {
	c := configuration{}
	readJson("config.json", &c)
	if c == nil {
		log.Fatal("wrong config file.")
	}
	return c
}
