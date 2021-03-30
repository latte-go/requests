package requests

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type Response struct {
	time int64
	url string
	resp *http.Response
	body []byte
}

func (this *Response)Response() *http.Response{
	if this != nil{
		return this.resp
	}
	return nil
}

func (this *Response)StatusCode() int{
	if this.resp == nil{
		return 0
	}
	return this.resp.StatusCode
}

func (this *Response)Url() string{
	if this != nil{
		return this.url
	}
	return ""
}

func (this *Response)Headers() http.Header{
	if this != nil{
		return this.resp.Header
	}
	return nil
}

func (this *Response)Cookies()[]*http.Cookie{
	if this != nil{
		return this.resp.Cookies()
	}
	return []*http.Cookie{}
}

func (this *Response)Body()([]byte,error){
	if this == nil{
		return []byte{},errors.New("HttpRequest.Response is nil.")
	}
	defer this.resp.Body.Close()

	if len(this.body) > 0{
		return this.body,nil
	}

	if this.resp == nil || this.resp.Body == nil{
		return nil, errors.New("response or body is nil")
	}

	b,err := ioutil.ReadAll(this.resp.Body)
	if err != nil{
		return nil, err
	}
	this.body = b

	return b, err
}

func (this *Response)BodyText() (string,error){
	b,err := this.Body()
	if err != nil{
		return "",nil
	}
	return string(b),nil
}

func (this *Response)BodyToMap()(map[string]interface{},error){
	b,err := this.Body()
	if err != nil{
		return map[string]interface{}{}, err
	}

	r := map[string]interface{}{}
	if err := json.Unmarshal(b,&r);err != nil{
		return map[string]interface{}{},err
	}
	return r,nil
}

func (this *Response)BodyToStruct(v interface{})(error){
	b,err := this.Body()
	if err != nil{
		return err
	}

	if err := json.Unmarshal(b,&v);err != nil{
		return err
	}
	return nil
}

func (this *Response)Close()error{
	if this != nil{
		return this.resp.Body.Close()
	}
	return nil
}