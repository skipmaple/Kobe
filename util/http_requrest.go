package util

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func HttpClient() *http.Client {
	client := &http.Client{Timeout: 10 * time.Second}
	return client
}

func SendRequest(client *http.Client, method string, endpoint string, values map[string]string) []byte {
	//endpoint := "https://httpbin.org/post"
	//values := map[string]string{"foo": "baz"}
	jsonData, err := json.Marshal(values)

	req, err := http.NewRequest(method, endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Error Occurred. %+v", err)
	}

	response, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request to API endpoint. %+v", err)
	}

	// Close the connection to reuse it
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Couldn't parse response body. %+v", err)
	}

	return body
}

// Usage Example:

//func main() {
//	// c should be re-used for further calls
//	c := HttpClient()
//	response := SendRequest(c, http.MethodPost, "https://httpbin.org/post", map[string]string{"foo": "baz"})
//	log.Println("Response Body:", string(response))
//}
