package helpers

import "encoding/json"

func PrintError(err error) {
	if err != nil {
		Logger.Error(err.Error())
	}
}

func Encode(data interface{}) []byte {
	response, err := json.Marshal(data)
	PrintError(err)
	return response
}
