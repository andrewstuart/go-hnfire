package hnfire

import "fmt"

var hnBase Endpoint = "https://hacker-news.firebaseio.com/v0"

//An Endpoint is a convenience function for getting nested API strings.
type Endpoint string

//Child returns the current endpoint concatenated with the result of formatting
//the given format string with the passed arguments.
func (ep Endpoint) Child(format string, args ...interface{}) Endpoint {
	args = append([]interface{}{ep}, args...)
	return Endpoint(fmt.Sprintf("%s/"+format, args...))
}

//String returns the Endpoint's underlying string value
func (ep Endpoint) String() string {
	return string(ep)
}
