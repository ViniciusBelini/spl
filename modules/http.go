package modules

import (
	// "fmt"
	"bytes"
	// "encoding/json"
	"io"
	"net/http"
)

type HttpResponse struct{
	Status  int
	Headers map[string]string
	Body    string
}

func makeHeaders(h map[any]any) http.Header{
	headers := http.Header{}
	for k, v := range h{
		headers.Set(k.(string), v.(string))
	}
	return headers
}

func HttpGet(url string, headers map[any]any, method string, jsonData string) (*HttpResponse, error){
	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(jsonData)))
	if err != nil {
		return nil, err
	}
	req.Header = makeHeaders(headers)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil{
		return nil, err
	}

	hmap := map[string]string{}
	for k, v := range resp.Header{
		hmap[k] = v[0]
	}

	return &HttpResponse{
		Status:  resp.StatusCode,
		Headers: hmap,
		Body:    string(bodyBytes),
	}, nil
}
