### 类似于Python Requests 的golang http请求库

#### requests 方法


**1、Post**

    req := requests.NewRequest()   
    
     params := map[string]interface{}{
        "page":123,
	}
    
    
    req.Post("http://192.168.0.125:8090/agent/getAgentList",params)

**1、Get**

    req := requests.NewRequest()   
    
     params := map[string]interface{}{
        "page":123,
    }
    
    
    req.Get("http://192.168.0.125:8090/agent/getAgentList",params)



**3、SetHeaders**

    req := requests.NewRequest()
    
    headers := map[string]string{
        "Content-Type":"application/json",
    }

    req.SetHeaders(headers)

**4、SetCookies**

    ` req := requests.NewRequest().Post()
    
    cookies := map[string]string{
    "uid":"123",
    }
    
    req.SetCookies(cookies).Post()`

#### response 方法

    resp, err := req.Post("http://127.0.0.1:8000/") //res is a http.Response object

**StatusCode() int **

    resp.StatusCode()

**Body() ([]byte,error)**

    body,err := resp.Body()
    fmt.Println(string(body))

**Close() error**

    resp.Close()

**BodyText() (string error)**

    body,err := resp.BodyText()
    fmt.Println(body)

**BodyToMap()(map[string]interface{},error)**

    body,err := resp.BodyToMap()
    fmt.Println(body)

**BodyToStruct(v interface)(error)**
type bodyStruct struct {}
body := bodyStruct{}
err := resp.BodyToMap(&body)
fmt.Println(body)