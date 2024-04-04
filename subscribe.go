package main

import (
	"encoding/json"
	"fmt"
)

// add address to observer
func Subscribe(address string) bool {
	if _, ok := SubscribeAddrFilterID[address]; ok {
		return true
	}
	req, err := buildNewFilterRequest(address)
	if err != nil {
		return false
	}

	resp, err := SendPostRequest(req)
	if err != nil {
		return false
	}

	result := JsonRPCResponse{}
	err = json.Unmarshal(resp, &result)

	if err != nil || result.Result == nil || result.Error != nil {
		return false
	}

	SubscribeAddrFilterID[address] = fmt.Sprint(result.Result)
	return true
}
