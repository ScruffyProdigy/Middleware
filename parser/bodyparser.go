/*
Parser helps parse the http request for values or files sent in with it
*/
package parser

import (
	"errors"
	"fmt"
	"github.com/ScruffyProdigy/TheRack/httper"
	"github.com/ScruffyProdigy/TheRack/rack"
	"mime/multipart"
	"reflect"
)

/*
Form will make sure all form variables are parsed properly before 
*/
var Form rack.Func = func(vars map[string]interface{}, next func()) {
	V(vars).Parse()
	next()
}

/*
Multipart will load any files that were sent in.  Files larger than MaxBytes might not be loaded properly
*/
type Multipart struct {
	MaxBytes int64
}

func (this Multipart) Run(vars map[string]interface{}, next func()) {
	V(vars).ParseMultipartOfSize(this.MaxBytes)
	next()
}

/*
VarMultipart is like Multipart, in that it will load any files that were sent in.  
It uses a variable in the vars to determine what the maximum file size should be.
This is useful if the max size isn't constant but you can determine it some other way
*/
type VarMultipart struct {
	MaxBytesVar string
}

func (this VarMultipart) Run(vars map[string]interface{}, next func()) {
	v, ok := vars[this.MaxBytesVar]
	if !ok {
		next()
		return
	}

	val := reflect.ValueOf(v)
	kind := val.Kind()
	if kind < reflect.Int || kind > reflect.Int64 {
		next()
		return
	}

	maxSize := val.Int()
	V(vars).ParseMultipartOfSize(maxSize)
	next()
}

/*
V allows you to recast vars, and get access to a set of functions
*/
type V map[string]interface{}

/*
Parse() will parse the request for form values
*/
func (vars V) Parse() {
	r := (httper.V)(vars).GetRequest()
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
}

/*
ParseMultipartOfSize() will parse the request to load all files up to a maximum size
*/
func (vars V) ParseMultipartOfSize(maxSize int64) {
	r := (httper.V)(vars).GetRequest()
	err := r.ParseMultipartForm(maxSize)
	if err != nil {
		panic(err)
	}
}

/*
FormGetSlice() will get all strings at a key.
Most of the time there will only be one string per key, but FormGetSlice becomes useful when there are many in one key
*/
func (vars V) FormGetSlice(key string) ([]string, bool) {
	r := (httper.V)(vars).GetRequest()
	f := r.Form

	if f == nil {
		//form data not ready, parse to get data and retry
		vars.Parse()
		f = r.Form

		if f == nil {
			//still not ready, something is wrong
			return nil, false
		}
	}

	val, ok := f[key]
	return val, ok
}

/*
FormValue() will get the first string at a key.
It will return "" if there are no strings at that key
*/
func (vars V) FormValue(key string) string {
	val, _ := vars.FormGetValue(key)
	return val
}

/*
FormGetValue() will get the first string at a key.
It will also tell you whether or not there was a string at that value, so you can distinguish between an empty string and no value
*/
func (vars V) FormGetValue(key string) (string, bool) {
	slice, ok := vars.FormGetSlice(key)
	if !ok {
		return "", false
	}

	if len(slice) == 0 {
		return "", false
	}

	return slice[0], true
}

/*
FileGetAt will return the ith file stored at a key (using a 0-based offset)
*/
func (vars V) GetFileAt(key string, i int) (multipart.File, error) {
	r := (httper.V)(vars).GetRequest()

	if i < 0 {
		return nil, errors.New("No negative files")
	}

	m := r.MultipartForm
	if m == nil {
		return nil, errors.New("multipart not parsed")
	}

	slice, ok := m.File[key]
	if !ok {
		for k, _ := range m.File {
			return nil, errors.New("No file at key:" + key + " (try key:\"" + k + "\")")
		}
	}
	if slice == nil {
		return nil, errors.New("No files at key:" + key)
	}

	if len(slice) <= i {
		return nil, errors.New(fmt.Sprint("No File #", i, " at key:", key))
	}

	return slice[i].Open()
}

/*
FileGet will get the first file stored at a key
*/
func (vars V) GetFile(key string) (multipart.File, error) {
	return vars.GetFileAt(key, 0)
}
