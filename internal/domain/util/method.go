package contextutil

// Method é o tipo do method
type Method string

// consts que representam os methods
const (
	MethodGET    Method = "GET"
	MethodPOST   Method = "POST"
	MethodPUT    Method = "PUT"
	MethodPATCH  Method = "PATCH"
	MethodDELETE Method = "DELETE"
)
