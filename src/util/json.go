package util

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

func WriteJson(path string, object interface{}) error {
	data, err := json.Marshal(object)
	if err != nil {
		log.Print(err)
		return err
	}

	ioutil.WriteFile(path, data, 0666)
	return nil
}

func ReadJson(path string, object interface{}) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Print(err)
		return err
	}

	json.Unmarshal(data, &object)
	return nil
}
