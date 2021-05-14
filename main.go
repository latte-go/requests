package requests

import "time"

func NewRequest() *Requests {
	r := &Requests{
		timeout: time.Second * 300,
		headers: map[string]string{},
		cookies: map[string]string{},
	}
	return r
}
