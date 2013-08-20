package util

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"sync"
)

var locker sync.Mutex = sync.Mutex{}

func WriteJson(path string, object interface{}) error {
	locker.Lock()
	defer locker.Unlock()

	data, err := json.Marshal(object)
	if err != nil {
		log.Print(err)
		return err
	}

	return ioutil.WriteFile(path, data, 0666)
}

func ReadJson(path string, object interface{}) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Print(err)
		return err
	}

	return json.Unmarshal(data, &object)
}
