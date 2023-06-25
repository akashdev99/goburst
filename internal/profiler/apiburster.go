package profiler

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/schollz/progressbar/v3"
)

type ApiBurster struct {
	Method  string
	ApiUrl  string
	Headers []string
	// ApiBody    string
	Iterations int
	Done       chan bool
	Err        chan error
}

// Need to add body
// For post requests
func NewApiBurster(method string, url string, headers []string, iterations int) *ApiBurster {
	return &ApiBurster{
		Method:     method,
		ApiUrl:     url,
		Headers:    headers,
		Iterations: iterations,
		Done:       make(chan bool, 1),
		Err:        make(chan error, 1),
	}
}

func (profiler *ApiBurster) BurstRequests(method string, url string, headers []string, iteration int) {
	httpClient, request, err := createHttpClient(method, url, headers)
	if err != nil {
		fmt.Printf("Failed create http client : %v \n", err)
		profiler.Err <- err
		return
	}

	bar := progressbar.Default(int64(iteration))

	startTime := time.Now().Unix()
	for i := 0; i < iteration; i++ {
		err := makeRequest(httpClient, request)
		if err != nil {
			fmt.Printf("API Failed , stopping profiling at count %v!!! : %v \n", i, err)
			profiler.Err <- err
			return
		}
		bar.Add(1)
	}

	endTime := time.Now().Unix()
	fmt.Printf("Total Time took to complete %v request = %v second \n", iteration, endTime-startTime)
	profiler.Done <- true
}

func createHttpClient(method string, url string, headers []string) (*http.Client, *http.Request, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}

	var key, value string
	for _, header := range headers {
		key, value = getHeaderKeyValue(header)
		req.Header.Add(key, value)
	}

	return client, req, nil
}

func getHeaderKeyValue(header string) (string, string) {
	headerPair := strings.Split(header, ":")
	return strings.TrimSpace(headerPair[0]), strings.TrimSpace(headerPair[1])
}

func makeRequest(client *http.Client, req *http.Request) error {
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if res.StatusCode != 200 {
		fmt.Printf("Response Body \n: %v\n", string(body))
		return fmt.Errorf("response status code %d", res.StatusCode)
	}
	return nil
}
