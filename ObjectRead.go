package main

import (
	"encoding/json"
	"os"
)

func ObjectRead(object any, file string) error {

	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &object)
	if err != nil {
		return err
	}

	return nil
}
