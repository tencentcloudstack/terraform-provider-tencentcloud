package common

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type Response interface {
	ParseErrorFromHTTPResponse(body []byte) error
}

type BaseResponse struct {
}

type ErrorResponse struct {
	Response struct {
		Error struct {
			Code    string `json:"Code"`
			Message string `json:"Message"`
		} `json:"Error" omitempty`
		RequestId string `json:"RequestId"`
	} `json:"Response"`
}

type DeprecatedAPIErrorResponse struct {
	Code     int    `json:"code"`
	Message  string `json:"message"`
	CodeDesc string `json:"codeDesc"`
}

func (r *BaseResponse) ParseErrorFromHTTPResponse(body []byte) (err error) {
	resp := &ErrorResponse{}
	err = json.Unmarshal(body, resp)
	if err != nil {
		return
	}
	if resp.Response.Error.Code != "" {
		return NewAPIError(resp.Response.Error.Code, resp.Response.Error.Message, -1)
	}

	deprecated := &DeprecatedAPIErrorResponse{}
	err = json.Unmarshal(body, deprecated)
	if err != nil {
		return
	}
	if deprecated.Code != 0 {
		return NewAPIError(deprecated.CodeDesc, deprecated.Message, deprecated.Code)
	}
	return nil
}

func ParseFromHttpResponse(hr *http.Response, response Response) (err error) {
	defer hr.Body.Close()
	body, err := ioutil.ReadAll(hr.Body)
	if err != nil {
		return
	}
	log.Printf("[DEBUG] Response Body=%s", body)
	err = response.ParseErrorFromHTTPResponse(body)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &response)
	return
}
