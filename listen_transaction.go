package main

import (
	"encoding/json"
	"time"
)

func buildNewFilterRequest(address string) ([]byte, error) {
	params := make(map[string]string, 0)
	params["address"] = address

	req := JsonRPCRequest{
		Jsonrpc: "2.0",
		Method:  "eth_newFilter",
		Params:  []interface{}{params},
		ID:      int16(1),
	}
	reqByte, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	return reqByte, err
}

func buildFilterChangesRequest(filterID string) ([]byte, error) {
	req := JsonRPCRequest{
		Jsonrpc: "2.0",
		Method:  "eth_getFilterChanges",
		Params:  []interface{}{filterID},
		ID:      int16(1),
	}
	reqByte, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	return reqByte, err
}

func getFilterChanges(filterID string) ([]Transaction, error) {
	req, err := buildFilterChangesRequest(filterID)
	if err != nil {
		return nil, err
	}

	resp, err := SendPostRequest(req)
	if err != nil {
		return nil, err
	}

	result := JsonRPCFilterChangeResponse{}
	err = json.Unmarshal(resp, &result)

	if err != nil || &result == nil || result.Error != nil {
		return nil, err
	}
	return result.Result, nil
}

func ListenToTransactions() {
	for {
		for addr, filterID := range SubscribeAddrFilterID {
			tl, err := getFilterChanges(filterID)
			if err != nil {
				continue
			}
			SubscribedAddrTx[addr] = append(SubscribedAddrTx[addr], tl...)
		}

		time.Sleep(10 * time.Second)
	}
}

func GetTransactions(address string) []Transaction {
	if txs, ok := SubscribedAddrTx[address]; ok {
		return txs
	}
	return []Transaction{}
}
