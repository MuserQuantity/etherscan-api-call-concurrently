package main

import (
	"github.com/MuserQuantity/etherscan-api-call-concurrently/core"
	"io"
	"log"
	"math/rand"
	"net/url"
	"strconv"
)

func main() {
	var apiKeys = []string{
		"",
	}
	caller := core.InitApiCaller(apiKeys)
	for i := 0; i < 50; i++ {
		var params = url.Values{}
		params.Add("module", "proxy")
		params.Add("action", "eth_getCode")
		params.Add("address", strconv.Itoa(rand.Int()))
		go func(caller *core.EtherscanApiCaller, params url.Values) {
			resp, err := caller.Request("https://api.etherscan.io/api", params)
			if err != nil {
				log.Println(err)
				return
			}
			defer resp.Body.Close()
			byteData, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Println(err)
				return
			}
			log.Println(string(byteData))
		}(caller, params)
	}
	for {

	}
}
