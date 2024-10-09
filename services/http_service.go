package services

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/valyala/fasthttp"
)

func doRequest(url string, requestMethod string, response interface{}) error {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(url)
	req.Header.SetMethod(requestMethod)
	req.Header.Set("Accept", "application/json")

	if err := fasthttp.Do(req, resp); err != nil {
		log.Println("[DoRequest] Error doing request", err)

		return err
	}

	if resp.StatusCode() != fasthttp.StatusOK {
		log.Println("[DoRequest] Status code is not OK", resp.StatusCode())

		return errors.New("status code is not OK")
	}

	if err := json.Unmarshal(resp.Body(), response); err != nil {
		log.Println("[DoRequest] Error deserializing request response", err)

		return err
	}

	return nil
}

func GetRequest(url string, response interface{}) error {
	return doRequest(url, "GET", response)
}
