package main

import (
	"encoding/json"
	"os"
)

func WriteObject(object any, name string) error {

	// Сериализуем данные
	var response, err = json.MarshalIndent(object, "", "\t")
	if err != nil {
		return err
	}

	// Создаем файл
	file, err := os.Create(name)
	if err != nil {
		return err
	}

	// Пишем в файл
	_, err = file.WriteString(string(response))
	if err != nil {
		return err
	}

	// Закрываем
	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}
