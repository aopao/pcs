package utils

import (
	"net/http"
	"io/ioutil"
	"github.com/astaxie/beego/logs"
	"net/url"
)

//Send Request by Get Method
func GetRequest(url string) ([]byte, error) {
	//http.NewRequest()
	response, err := http.Get(url)
	if err != nil {
		logs.Error(err)
		//os.Exit(1)
		return nil, err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logs.Error(err)
		//os.Exit(1)
		return nil, err
	}
	return body, nil
}

//Send Request by Get Method with some params
func GetRequestWithParams(url string, params map[string]string) ([]byte, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logs.Error(err)
		//os.Exit(1)
		return nil, err
	}
	q := request.URL.Query()
	for k, v := range params {
		logs.Info("GET:paramKey=%v, paramValue=%v\n", k, v)
		q.Set(k, v)
	}

	request.URL.RawQuery = q.Encode()

	logs.Info(request.URL.String())

	//var resp *http.Response
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	//logs.Info(json.NewDecoder(response.Body))
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	return body, nil
}

//Send Request by Post Method with some params
func PostRequest(urladdr string, params map[string]string) ([]byte, error) {
	var requestParam = url.Values{}
	for k, v := range params {
		logs.Info("POST:paramKey=%v, paramValue=%v\n", k, v)
		requestParam.Set(k, v)
	}
	resp, err := http.PostForm(urladdr, requestParam)

	if err != nil {
		logs.Error(err)
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	return body, nil
}
