package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"service-discovery/database"
	"strings"
)

func FetchDataFromDatabase() {
	arr := database.ListCollectionNames(database.Database())
	length := len(arr)
	for i := 0; i < length; i++ {
		fmt.Println(arr[i])
		//fmt.Println(string(GetDataFromDatabase(arr[i])))
	}
}

func GetDataFromDatabase(collection string, id string, data []byte) (response bool) {
	//var s []string
	arr := database.GetAllDocuments(database.Database(), collection)
	for _, doc := range arr {
		//fmt.Println(doc)
		if doc["id"] == id {
			result, err := json.Marshal(doc)
			if err != nil {
				fmt.Println(err)
			}

			cloud_data := strings.Replace(string(data), "{", "", -1)
			db_data := strings.Replace(string(result), "{", "", -1)
			cloudData := strings.Split(string(cloud_data), ",")
			database := strings.Split(string(db_data), ",")

			//cloudData = RemoveIndex(cloudData, 0)
			databaseData := RemoveIndex(database, 0)
			//databaseData1 := RemoveIndex(databaseData, 0)
			fmt.Println("From cloud: ", cloudData)
			fmt.Println("From database: ", databaseData)
			res := bytes.Compare(data, result)
			fmt.Println("res: ", res)
			if res == 0 {
				response = true
			} else {
				response = false
			}
		} else {
			response = false
		}

		//fmt.Println(string(result))
		//s = append(s, string(result))
	}
	return response
}

func RemoveIndex(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}
