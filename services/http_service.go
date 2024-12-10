package services

import (
	"errors"
	"log"

	"github.com/valyala/fasthttp"
)

func doRequest(url string, requestMethod string) ([]byte, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(url)
	req.Header.SetMethod(requestMethod)
	req.Header.Set("Accept", "application/json")

	if err := fasthttp.Do(req, resp); err != nil {
		log.Println("[DoRequest] Error doing request", err)

		return nil, err
	}

	if resp.StatusCode() != fasthttp.StatusOK {
		log.Println("[DoRequest] Status code is not OK", resp.StatusCode())

		return nil, errors.New("status code is not OK")
	}

	return resp.Body(), nil
}

func GetRequest(url string) ([]byte, error) {
	return doRequest(url, "GET")
}
