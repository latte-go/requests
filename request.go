package httpRequest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"time"
)

type REST_METHOD string

type Requests struct {
	url string
	data interface{}
	timeout time.Duration
	headers map[string]string
	cookies map[string]string
	client *http.Client
	method string
	BodyBuffer *bytes.Buffer
	err error
}

func (this *Requests)Request (method,url string,data ...interface{}) (r *Response,err error)  {
	response := &Response{}

	start := time.Now().UnixNano() / 1e6
	defer this.elapsedTime(start,response)

	if method == "" || url == ""{
		return nil, errors.New("parameter method and url is required")
	}

	r.url = url
	if len(data) > 0 {
		this.data = data
	}else {
		this.data = ""
	}

	var (
		req  *http.Request
		body  io.Reader
	)
	this.client = this.buildClient()
	this.method = method

	if method == "GET" || method == "DELETE"{
		url,err = buildUrl(url,this.data)
		if err != nil{
			return nil, err
		}
		r.url = url
	}

	body,err = this.buildBody(data...)
	if err != nil{
		return nil, err
	}

	req,err = http.NewRequest(method,url,body)
	if err != nil{
		return nil, err
	}

	resp,err := this.client.Do(req)
	if err != nil{
		return nil, err
	}

	response.url = url
	response.resp = resp
	return response,nil
}


func (this *Requests)Post(url string,data interface{}) (response *Response,err error) {
	return  this.Request("post",url,data)
}

func (this *Requests)Get(url string,data interface{})  (response *Response,err error) {
	return this.Request("get",url,data)
}

func (this *Requests)ResponseJson(data interface{})  {


}

func (this *Requests)elapsedTime(n int64,resp *Response){
	end := time.Now().UnixNano()/1e6
	resp.time = end - n
}

func (this *Requests)buildClient() *http.Client{
	if this.client == nil{
		this.client = &http.Client{
			Timeout: time.Second * this.timeout,
		}
	}
	return this.client
}

func (this *Requests) buildBody(d ...interface{}) (io.Reader, error) {
	if this.method == "GET" || this.method == "DELETE" || len(d) == 0 || (len(d) > 0 && d[0] == nil) {
		return nil, nil
	}

	switch d[0].(type) {
	case string:
		return strings.NewReader(d[0].(string)), nil
	case []byte:
		return bytes.NewReader(d[0].([]byte)), nil
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return bytes.NewReader(IntByte(d[0])), nil
	case *bytes.Reader:
		return d[0].(*bytes.Reader), nil
	case *strings.Reader:
		return d[0].(*strings.Reader), nil
	case *bytes.Buffer:
		return d[0].(*bytes.Buffer), nil
	default:
		if this.isJson() {
			b, err := json.Marshal(d[0])
			if err != nil {
				return nil, err
			}
			return bytes.NewReader(b), nil
		}
	}

	t := reflect.TypeOf(d[0]).String()
	if !strings.Contains(t, "map[string]interface") {
		return nil, errors.New("Unsupported data type.")
	}

	data := make([]string, 0)
	for k, v := range d[0].(map[string]interface{}) {
		if s, ok := v.(string); ok {
			data = append(data, fmt.Sprintf("%s=%v", k, s))
			continue
		}
		b, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		data = append(data, fmt.Sprintf("%s=%s", k, string(b)))
	}

	return strings.NewReader(strings.Join(data, "&")), nil
}

func (this *Requests) isJson() bool {
	if len(this.headers) > 0 {
		for _, v := range this.headers {
			if strings.Contains(strings.ToLower(v), "application/json") {
				return true
			}
		}
	}
	return false
}

func buildUrl(url string,data ...interface{}) (r string,err error){
	return
}
