package requests

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"
)

type REST_METHOD string

type Requests struct {
	client     *http.Client
	url        string
	data       interface{}
	BodyBuffer *bytes.Buffer
	err        error

	method            string
	timeout           time.Duration
	disableKeepAlives bool

	headers   map[string]string
	cookies   map[string]string
	transport *http.Transport
	proxy     func(r *http.Request) (*url.URL, error)
}

func (this *Requests) Request(method, url string, data ...interface{}) (r *Response, err error) {
	method = strings.ToUpper(method)
	response := &Response{}

	start := time.Now().UnixNano() / 1e6
	defer this.elapsedTime(start, response)

	if method == "" || url == "" {
		return nil, errors.New("parameter method and url is required")
	}

	this.url = url
	if len(data) > 0 {
		this.data = data[0]
	} else {
		this.data = ""
	}

	var (
		req  *http.Request
		body io.Reader
	)
	this.client = this.buildClient()
	this.method = method

	if method == "GET" || method == "DELETE" {
		url, err = buildUrl(url, this.data)
		if err != nil {
			return nil, err
		}
		this.url = url
	}

	body, err = this.buildBody(data...)
	if err != nil {
		return nil, err
	}

	req, err = http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	this.initHeaders(req)
	this.initCookies(req)

	resp, err := this.client.Do(req)
	if err != nil {
		return nil, err
	}

	response.url = url
	response.resp = resp
	return response, nil
}

func (this *Requests) DisableKeepAlives(v bool) *Requests {
	this.disableKeepAlives = v
	return this
}

// Set headers
func (this *Requests) SetHeaders(headers map[string]string) *Requests {
	if headers != nil || len(headers) > 0 {
		for k, v := range headers {
			this.headers[k] = v
		}
	}
	return this
}

// Init headers
func (this *Requests) initHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for k, v := range this.headers {
		req.Header.Set(k, v)
	}
}

// Set cookies
func (this *Requests) SetCookies(cookies map[string]string) *Requests {
	if cookies != nil || len(cookies) > 0 {
		for k, v := range cookies {
			this.cookies[k] = v
		}
	}
	return this
}

// Init cookies
func (this *Requests) initCookies(req *http.Request) {
	for k, v := range this.cookies {
		req.AddCookie(&http.Cookie{
			Name:  k,
			Value: v,
		})
	}
}

func (this *Requests) SetPorxy(v func(r *http.Request) (*url.URL, error)) *Requests {
	this.proxy = v
	return this
}

func (this *Requests) SetTimeout(d time.Duration) *Requests {
	this.timeout = d
	return this
}

func (this *Requests) Transport(v *http.Transport) *Requests {
	this.transport = v
	return this
}

func (this *Requests) getTransport() http.RoundTripper {
	if this.transport == nil {
		return http.DefaultTransport
	}

	this.transport.DisableKeepAlives = this.disableKeepAlives

	if this.proxy != nil {
		this.transport.Proxy = this.proxy
	}

	return http.RoundTripper(this.transport)
}

func (this *Requests) Post(url string, data interface{}) (response *Response, err error) {
	return this.Request("post", url, data)
}

func (this *Requests) Get(url string, data interface{}) (response *Response, err error) {
	return this.Request("get", url, data)
}

func (this *Requests) elapsedTime(n int64, resp *Response) {
	end := time.Now().UnixNano() / 1e6
	resp.time = end - n
}

func (this *Requests) buildClient() *http.Client {
	if this.client == nil {
		this.client = &http.Client{
			Timeout:   time.Second * this.timeout,
			Transport: this.getTransport(),
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

func parseQuery(url string) ([]string, error) {
	urlList := strings.Split(url, "?")
	if len(urlList) < 2 {
		return make([]string, 0), nil
	}
	query := make([]string, 0)
	for _, val := range strings.Split(urlList[1], "&") {
		v := strings.Split(val, "=")
		if len(v) < 2 {
			return make([]string, 0), errors.New("query parameter error")
		}
		query = append(query, fmt.Sprintf("%s=%s", v[0], v[1]))
	}
	return query, nil
}

func buildUrl(url string, data ...interface{}) (r string, err error) {
	query, err := parseQuery(url)
	if err != nil {
		return url, err
	}

	if len(data) > 0 && data[0] != nil {
		t := reflect.TypeOf(data[0]).String()
		switch t {
		case "map[string]interface {}":
			for k, v := range data[0].(map[string]interface{}) {
				vv := ""
				if reflect.TypeOf(v).String() == "string" {
					vv = v.(string)
				} else {
					b, err := json.Marshal(v)
					if err != nil {
						return url, err
					}
					vv = string(b)
				}
				query = append(query, fmt.Sprintf("%s=%s", k, vv))
			}
		case "string":
			param := data[0].(string)
			if param != "" {
				query = append(query, param)
			}
		default:
			return url, errors.New("Unsupported data type.")
		}

	}

	list := strings.Split(url, "?")

	if len(query) > 0 {
		return fmt.Sprintf("%s?%s", list[0], strings.Join(query, "&")), nil
	}

	return list[0], nil
}
