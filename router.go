package begonia

import "io"

type Router interface {
	Redirect(uri, to string, code int)
	Get(uri string, callback Callback)
	Post(uri string, callback Callback)
	Put(uri string, callback Callback)
	Patch(uri string, callback Callback)
	Delete(uri string, callback Callback)
	Options(uri string, callback Callback)
}

// Should be a function
type Callback interface{}

type View interface {
	Status() int
	Serialize() []byte
}

type Request struct {
}

// Input grabs input from the request body by name
func (req *Request) Input(name string) string {
	return "some input"
}

// Input grabs input from the query string by name
func (req *Request) Query(name string) string {
	panic("Implement me")
}

// Cookie grabs input from cookies by name
func (req *Request) Cookie(name string) string {
	panic("Implement me")
}

// Cookie grabs file name
func (req *Request) File(name string) []byte {
	panic("Implement me")
}

type Response struct {
	W io.Writer
}

func (res *Response) Cookie(name, value string) *Response {
	panic("Implement me")
	return res
}

func (res *Response) Status(status int) *Response {
	panic("Implement me")
	return res
}

func (res *Response) View(data ...interface{}) View {
	panic("Implement me")
	return nil
}
