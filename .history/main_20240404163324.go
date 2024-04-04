package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const ENTRYPOINT = "https://cloudflare-eth.com"

var latestBlockNumber int
var latestRequestTimestamp int64
var SubscribeAddrFilterID map[string]string
var SubscribedAddrTx map[string][]Transaction

func main() {

	http.HandleFunc("/", handleRequest)

	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/getcurrentblock":
		if r.Method == "GET" {
			getCurrentBlockHandler(w, r)
		} else {
			respondWithError(w, http.StatusNotFound, "Method is not supported for /getcurrentblock")
		}
	case "/subscribe":
		if r.Method == "POST" {
			subscribeHandler(w, r)
		} else {
			respondWithError(w, http.StatusNotFound, "Method is not supported for /subscribe")
		}
	case "/getTransactions":
		if r.Method == "GET" {
			getTransactionsHandler(w, r)
		} else {
			respondWithError(w, http.StatusNotFound, "Method is not supported for /getTransactions")
		}
	default:
		respondWithError(w, http.StatusNotFound, "Path not found.")
	}
}

func getCurrentBlockHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		respondWithError(w, http.StatusNotFound, "Method is not supported.")
		return
	}

	blockNumber := GetCurrentBlock()
	respondWithJSON(w, http.StatusOK, map[string]int{"current_block": blockNumber})
}

func subscribeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		respondWithError(w, http.StatusNotFound, "Method is not supported.")
		return
	}

	var data struct {
		Address string `json:"address"`
	}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	success := Subscribe(data.Address)
	if !success {
		respondWithError(w, http.StatusInternalServerError, "Failed to subscribe.")
		return
	}
	respondWithJSON(w, http.StatusCreated, map[string]string{"message": "Subscribed to address: " + data.Address})
}

func getTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Address string `json:"address"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	transactions := GetTransactions(data.Address)
	respondWithJSON(w, http.StatusOK, transactions)
}

func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}
