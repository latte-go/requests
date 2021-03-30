package requests

func NewRequest() *Requests{
	r := &Requests{
		timeout: 30,
		headers: map[string]string{},
		cookies: map[string]string{},
	}
	return r
}