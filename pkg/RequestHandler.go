package pkg

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/tiaguinho/gosoap"
)

type HandleBody func(body string) ([]string, error)

func GetRequestUrl(indexUrl string, handleResult HandleBody) ([]string, error) {
	if indexUrl == "" {
		return nil, errors.New("the index url is empty")
	}
	resp, err := http.Get(indexUrl)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("status code is error")
	}
	defer resp.Body.Close()
	if body, err := ioutil.ReadAll(resp.Body); err != nil {
		return nil, err
	} else {
		if handleResult != nil {
			return handleResult(string(body))
		} else {
			fmt.Println(string(body))
			return nil, nil
		}
	}

}

type TransResponse struct {
	TranscodeResult string `xml:"transcodeResult"`
}

func Trans(soapUrl string, request string) (string, error) {
	httpClient := &http.Client{
		Timeout: 1500 * time.Millisecond,
	}
	soap, err := gosoap.SoapClient(soapUrl, httpClient)
	if err != nil {
		log.Fatalf("SoapClient error: %s", err)
		return "", err
	}
	params := gosoap.Params{
		"str": request,
	}
	res, err := soap.Call("transcode", params)
	if err != nil {
		log.Fatalf("Call error: %s", err)
		return "", err
	}
	trasRes := TransResponse{}
	res.Unmarshal(&trasRes)
	return trasRes.TranscodeResult, nil
}
