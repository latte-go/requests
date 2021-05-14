package requests

import "time"

func NewRequest() *Requests {
	r := &Requests{
		timeout: time.Second * 600,
		headers: map[string]string{},
		cookies: map[string]string{},
	}
	return r
}
