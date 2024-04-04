package main

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
)

func buildGetCurrentBlockRequest() ([]byte, error) {
	req := JsonRPCRequest{
		Jsonrpc: "2.0",
		Method:  "eth_blockNumber",
		Params:  make([]interface{}, 0),
		ID:      int16(1),
	}
	reqByte, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	return reqByte, err
}

func GetCurrentBlock() int {
	// if latestBlockNumber != 0 || time.Now().Unix()-latestRequestTimestamp < 10 {
	// 	return latestBlockNumber
	// }

	req, err := buildGetCurrentBlockRequest()
	if err != nil {
		return 0
	}

	resp, err := SendPostRequest(req)
	if err != nil {
		return 0
	}

	result := JsonRPCResponse{}
	err = json.Unmarshal(resp, &result)
	if err != nil && &result == nil {
		return 0
	}

	hexResult := strings.TrimPrefix(fmt.Sprint(result.Result), "0x")
	bi := new(big.Int)
	_, success := bi.SetString(hexResult, 16)
	if success == false {
		return 0
	}
	number := bi.Uint64()
	return int(number)
}
