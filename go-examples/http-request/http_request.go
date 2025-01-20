package http_request

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func HTTPrequest() {
	resp, err := http.Get("https://jsonplaceholder.typicode.com/posts/1")
	if err != nil {
		log.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	// body := resp.Body.Read()
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("not a json string")
	}
	for k, v := range result {
		fmt.Println("key - ", k, " value - ", v)
	}

}
