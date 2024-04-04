package main

import (
	"bytes"
	"io"
	"net/http"
)

type JsonRPCRequest struct {
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int16         `json:"id"`
}
type ResponseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type JsonRPCResponse struct {
	Jsonrpc string         `json:"jsonrpc"`
	Result  interface{}    `json:"result"`
	ID      int            `json:"id"`
	Error   *ResponseError `json:"error"`
}

type Transaction struct {
	Hash             string `json:"hash"`
	From             string `json:"from"`
	To               string `json:"to"`
	Value            string `json:"value"`
	BlockNumber      string `json:"blockNumber"`
	TransactionIndex string `json:"transactionIndex"`
}

type JsonRPCFilterChangeResponse struct {
	Jsonrpc string         `json:"jsonrpc"`
	Result  []Transaction  `json:"result"`
	ID      int            `json:"id"`
	Error   *ResponseError `json:"error"`
}

func SendPostRequest(requestParams []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", ENTRYPOINT, bytes.NewBuffer(requestParams))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}
