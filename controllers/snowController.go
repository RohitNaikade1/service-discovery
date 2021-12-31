package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
		Logger.Error("Error creating http request: " + err.Error())
	}

	req.Header.Add("Authorization", "Basic YWRtaW46OHRsVE5wTlNzTTFy")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		Logger.Error("Error sending http request: " + err.Error())
	}

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		Logger.Error("Error reading http response body: " + err.Error())
	}

	Logger.Info(string(responseBody))
	Logger.Info(res.Status)

	var s models.Results
	json.Unmarshal(responseBody, &s)

	return s.Result.SysID
}

func DeActive(data string) {
	payload := strings.NewReader(data)
	url := "https://dev55842.service-now.com/api/631287/pocapi/deactive"
	method := "POST"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		Logger.Error(err.Error())
		return
	}
	req.Header.Add("Authorization", "Basic YWRtaW46OHRsVE5wTlNzTTFy")

	res, err := client.Do(req)
	if err != nil {
		Logger.Error(err.Error())
		return
	}
	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		Logger.Error(err.Error())
		return
	}
	Logger.Info(string(responseBody))
	Logger.Info(res.Status)
}
