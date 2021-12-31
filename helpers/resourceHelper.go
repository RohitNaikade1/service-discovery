package helpers

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"service-discovery/models"
)

func Contains(arr []string, resourceType string) (result bool) {
	result = false
	for i := 0; i < len(arr); i++ {
		if arr[i] == resourceType {
			result = true
		}
	}
	return result
}

func ValidateResource(resourceType string) (result bool, arr []string) {
	jsonFile, err := os.Open("resources.json")
	if err != nil {
		Logger.Error(err.Error())
	}

	Logger.Info("Successfully Opened users.json")

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var resources models.Resources

	json.Unmarshal(byteValue, &resources)

	result = Contains(resources.Resource, resourceType)
	return result, resources.Resource
}
