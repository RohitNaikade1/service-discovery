package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"service-discovery/models"
	"strings"
)

func AddSysID(data map[string]interface{}, id string) map[string]interface{} {

	r := models.SysIDs{SysID: id, Status: "active"}
	sysresponse, err := json.Marshal(r)
	if err != nil {
		fmt.Println(err)
	}
	json.Unmarshal(sysresponse, &data)
	return data
}

func PostApi(body string) (sysid string) {

	payload := strings.NewReader(body)
	req, err := http.NewRequest(
		http.MethodPost,
		"https://dev55842.service-now.com/api/631287/pocapi/send",
		payload,
	)
	if err != nil {
		log.Fatalf("Error creating http request: %v", err)
	}

	req.Header.Add("Authorization", "Basic YWRtaW46OHRsVE5wTlNzTTFy")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Error sending http request: %v", err)
	}

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Error reading http response body: %v", err)
	}

	log.Println(string(responseBody))
	log.Println(res.Status)

	var s models.Results
	json.Unmarshal(responseBody, &s)

	//fmt.Println(s.Result.SysID)

	return s.Result.SysID
}

func DeActive(data string) {
	//fmt.Println(data)
	//body := fmt.Sprint(data)
	payload := strings.NewReader(data)
	url := "https://dev55842.service-now.com/api/631287/pocapi/deactive"
	method := "POST"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Authorization", "Basic YWRtaW46OHRsVE5wTlNzTTFy")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(responseBody))
	fmt.Println(res.Status)
}
