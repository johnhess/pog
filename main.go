package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

const port = 8080
const gatewayURL = "https://ohttp-gateway.jthess.com/gateway"
const respContentType = "message/ohttp-res"
const reqContentType = "message/ohttp-req"

func validateRequest(r *http.Request) error {
	if r.Method != http.MethodPost {
		return fmt.Errorf("unsupported method: %s", r.Method)
	}
	if r.Header.Get("Content-Type") != reqContentType {
		return fmt.Errorf("unsupported content type %s (expected: %s)", r.Header.Get("Content-Type"), reqContentType)
	}
	return nil
}

func validateGatewayResponse(r *http.Response) error {
	if r.Header.Get("Content-Type") != respContentType {
		return fmt.Errorf("unexpected content type from gateway: %s", r.Header.Get("Content-Type"))
	}
	return nil
}

func handleProxy(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Proxying request to gateway: %s\n", gatewayURL)
	// Validate the request
	if err := validateRequest(r); err != nil {
		http.Error(w, fmt.Sprintf("Bad request: %s", err.Error()), http.StatusBadRequest)
		return
	}
	fmt.Print("Request validated\n")

	gatewayRequest, err := http.NewRequest(http.MethodPost, gatewayURL, r.Body)
	gatewayRequest.Header.Set("Content-Type", reqContentType)
	if err != nil {
		// TODO: Carefully construct errors to conform to specification
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	fmt.Print("Request created\n")
	gatewayResponse, err := http.DefaultClient.Do(gatewayRequest)
	if err != nil {
		// TODO: Carefully construct errors to conform to specification
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	fmt.Print("Request sent\n")
	defer gatewayResponse.Body.Close()

	// Validate the gateway response
	if err := validateGatewayResponse(gatewayResponse); err != nil {
		http.Error(w, fmt.Sprintf("Bad gateway response: %s", err.Error()), http.StatusBadGateway)
		gatewayResponseAsString, _ := io.ReadAll(gatewayResponse.Body)
		fmt.Printf("Gateway responded in a bad way.  Status: %s.  Message: %s",
			gatewayResponse.Status, gatewayResponseAsString)
		return
	}
	fmt.Print("Response validated\n")

	fmt.Printf("Gateway response: %s\n", gatewayResponse.Status)
	// Copy the response body
	w.WriteHeader(http.StatusOK)
	io.Copy(w, gatewayResponse.Body)
}

func main() {
	fmt.Printf("Starting server on :%d\n", port)
	http.HandleFunc("/", handleProxy)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
