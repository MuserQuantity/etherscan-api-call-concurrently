package core

import (
	"net/http"
	"net/url"
	"sync"
	"time"
)

type EtherscanApiCaller struct {
	ApiKeys      []string
	ApiCallLock  []*sync.Mutex
	ApiCallCount []int
	ApiCallTime  []time.Time
	Index        int
	indexLock    *sync.Mutex
}

func InitApiCaller(apiKeys []string) *EtherscanApiCaller {
	length := len(apiKeys)
	var apiCallLock []*sync.Mutex
	var apiCallCount []int
	var apiCallTime []time.Time
	for i := 0; i < length; i++ {
		apiCallLock = append(apiCallLock, &sync.Mutex{})
		apiCallCount = append(apiCallCount, 0)
		apiCallTime = append(apiCallTime, time.Now())
	}
	return &EtherscanApiCaller{
		ApiKeys:      apiKeys,
		ApiCallLock:  apiCallLock,
		ApiCallCount: apiCallCount,
		ApiCallTime:  apiCallTime,
		Index:        0,
		indexLock:    &sync.Mutex{},
	}
}

// GetApiKey api key recurring acquisition
func (caller *EtherscanApiCaller) GetApiKey() (index int, apiKey string) {
	caller.indexLock.Lock()
	index = caller.Index
	apiKey = caller.ApiKeys[index]
	caller.Index = (caller.Index + 1) % len(caller.ApiKeys)
	caller.indexLock.Unlock()
	return
}

func (caller *EtherscanApiCaller) Request(url string, params url.Values, retry bool) (resp *http.Response, err error) {
	index, apiKey := caller.GetApiKey()
	caller.ApiCallLock[index].Lock()
	if caller.ApiCallCount[index] == 5 {
		var now = time.Now()
		subDuration := now.Sub(caller.ApiCallTime[index])
		if subDuration.Milliseconds() < 1200 {
			time.Sleep(time.Second - subDuration)
		}
		caller.ApiCallCount[index] = 0
	}
	if caller.ApiCallCount[index] == 0 {
		caller.ApiCallTime[index] = time.Now()
	}
	params.Add("apikey", apiKey)
	resp, err = http.Get(url + "?" + params.Encode())
	caller.ApiCallCount[index]++
	caller.ApiCallLock[index].Unlock()
	if err != nil {
		if retry {
			return caller.Request(url, params, retry)
		}
	}
	return
}
