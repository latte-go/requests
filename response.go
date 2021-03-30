package requests

import "net/http"

type Response struct {
	time int64
	url string
	resp *http.Response
	body []byte
}
