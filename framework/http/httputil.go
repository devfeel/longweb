// httpHelper
package httputil

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

//定义设置了超时时间的httpclient
var currClient *http.Client = &http.Client{
	Transport: &http.Transport{
		Dial: func(netw, addr string) (net.Conn, error) {
			c, err := net.DialTimeout(netw, addr, time.Second*300)
			if err != nil {
				fmt.Println("dail timeout", err)
				return nil, err
			}
			return c, nil
		},
		MaxIdleConnsPerHost:   10,
		ResponseHeaderTimeout: time.Second * 200,
	},
}

func HttpGet(url string) (body string, contentType string, intervalTime int64, errReturn error) {
	startTime := time.Now()
	intervalTime = 0
	contentType = ""
	body = ""
	errReturn = nil

	resp, err := currClient.Get(url)
	if err != nil {
		intervalTime = int64(time.Now().Sub(startTime) / time.Millisecond)
		errReturn = err
		return
	}
	defer resp.Body.Close()

	bytebody, err := ioutil.ReadAll(resp.Body)
	intervalTime = int64(time.Now().Sub(startTime) / time.Millisecond)
	if err != nil {
		intervalTime = int64(time.Now().Sub(startTime) / time.Millisecond)
		errReturn = err
		return
	}
	body = string(bytebody)
	contentType = resp.Header.Get("Content-Type")
	intervalTime = int64(time.Now().Sub(startTime) / time.Millisecond)
	return
}

func HttpPost(url string, postbody string, bodyType string) (body string, contentType string, intervalTime int64, errReturn error) {
	startTime := time.Now()
	intervalTime = 0
	contentType = ""
	body = ""
	errReturn = nil
	postbytes := bytes.NewBuffer([]byte(postbody))
	//resp, err := currClient.Post(url, "application/x-www-form-urlencoded", postbytes)
	//resp, err := currClient.Post(url, "application/json", postbytes)
	if bodyType == "" {
		bodyType = "application/x-www-form-urlencoded"
	}
	resp, err := currClient.Post(url, bodyType, postbytes)
	if err != nil {
		intervalTime = int64(time.Now().Sub(startTime) / time.Millisecond)
		errReturn = err
		return
	}
	defer resp.Body.Close()

	bytebody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		intervalTime = int64(time.Now().Sub(startTime) / time.Millisecond)
		errReturn = err
		return
	}
	body = string(bytebody)
	contentType = resp.Header.Get("Content-Type")
	intervalTime = int64(time.Now().Sub(startTime) / time.Millisecond)
	return

}

//从指定query集合获取指定key的值
func GetQuery(querys url.Values, key string) string {
	if len(querys[key]) > 0 {
		return querys[key][0]
	}
	return ""
}

//获取当前目录
func GetCurrentDirectory() (string, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "", err
	}
	return strings.Replace(dir, "\\", "/", -1), nil
}
