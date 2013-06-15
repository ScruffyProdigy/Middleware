package parser

import (
	"bytes"
	"fmt"
	"github.com/ScruffyProdigy/TheRack/httper"
	"github.com/ScruffyProdigy/TheRack/rack"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
)

func PostTo(loc string, vals url.Values) {
	resp, err := http.PostForm(loc, vals)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(string(body))
}

func writeMultipartFile(buf io.Writer, filename string) (string, error) {
	w := multipart.NewWriter(buf)
	defer w.Close()

	fw, err := w.CreateFormFile("file", "file")
	if err != nil {
		return "", err
	}

	fd, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer fd.Close()

	_, err = io.Copy(fw, fd)
	if err != nil {
		return "", err
	}

	return w.FormDataContentType(), nil
}

func SendFileTo(loc string, filename string) {
	buf := new(bytes.Buffer)
	filetype, err := writeMultipartFile(buf, filename)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	res, err := http.Post(loc, filetype, buf)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(string(body))
}

var HttpWare rack.Func = func(vars map[string]interface{}, next func()) {
	h := httper.V(vars)
	p := V(vars)

	name, named := p.FormGetValue("Name")
	if named {
		h.SetMessageString("Welcome, " + name)
	} else {
		h.SetMessageString("You are nameless")
	}
}

var FileWare rack.Func = func(vars map[string]interface{}, next func()) {
	h := httper.V(vars)
	p := V(vars)

	file, err := p.GetFile("file")
	if err != nil {
		h.SetMessageString("No File Found - " + err.Error())
		return
	}

	body, err := ioutil.ReadAll(file)
	if err != nil {
		h.SetMessageString("Can't Load File - " + err.Error())
	}

	h.SetMessage(body)
}

func init() {
	rackup := rack.New()
	rackup.Add(Form)
	rackup.Add(HttpWare)

	conn := httper.HttpConnection(":4007")
	go conn.Go(rackup)
}

func Example_None() {
	PostTo("http://localhost:4007", url.Values{})
	//output: You are nameless
}

func Example_Basic() {
	PostTo("http://localhost:4007", url.Values{"Name": {"Bob"}})
	//output: Welcome, Bob
}

func Example_Overloaded() {
	PostTo("http://localhost:4007", url.Values{"Name": {"Jim", "Bob"}})
	//output: Welcome, Jim
}

func Example_Empty() {
	PostTo("http://localhost:4007", url.Values{"Name": {}})
	//output: You are nameless
}

func Example_Incorrect() {
	PostTo("http://localhost:4007", url.Values{"irrelevant": {}})
}

func Example_Skipped() {
	rackup := rack.New()
	rackup.Add(HttpWare)

	conn := httper.HttpConnection(":3001")
	go conn.Go(rackup)

	PostTo("http://localhost:3001", url.Values{"Name": {"Jim"}})
	//output: Welcome, Jim
}

func Example_Multipart() {
	rackup := rack.New()
	rackup.Add(Multipart{256})
	rackup.Add(FileWare)

	conn := httper.HttpConnection(":3002")
	go conn.Go(rackup)

	SendFileTo("http://localhost:3002", "./test_files/helloworld.txt")
	//output: Hello World
}

var setSize rack.Func = func(vars map[string]interface{}, next func()) {
	vars["size"] = 256
	next()
}

func Example_VarMultipart() {
	rackup := rack.New()
	rackup.Add(setSize)
	rackup.Add(VarMultipart{"size"})
	rackup.Add(FileWare)

	conn := httper.HttpConnection(":3003")
	go conn.Go(rackup)

	SendFileTo("http://localhost:3003", "./test_files/helloworld.txt")
	//output: Hello World
}
