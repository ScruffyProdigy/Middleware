package sessioner

import (
	"fmt"
	"github.com/ScruffyProdigy/TheRack/httper"
	"github.com/ScruffyProdigy/TheRack/rack"
	"io/ioutil"
	"net/http"
)

func GetWithCookies(loc string, cookies []*http.Cookie) []*http.Cookie {
	req, err := http.NewRequest("get", loc, nil)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	fmt.Print(string(body))

	return res.Cookies()
}

var HelloWorldWare rack.Func = func(vars map[string]interface{}, next func()) {
	s := (V)(vars)
	times, ok := s.Get("times").(int)
	if !ok {
		times = 0
	}

	times++

	s.Set("times", times)

	(httper.V)(vars).SetMessageString(fmt.Sprint(times))
}

func Example_Session() {
	rackup := rack.New()
	rackup.Add(Middleware)
	rackup.Add(HelloWorldWare)

	conn := httper.HttpConnection(":3000")
	go conn.Go(rackup)

	var cookies []*http.Cookie
	for i := 0; i < 4; i++ {
		cookies = GetWithCookies("http://localhost:3000", cookies)
	}
	//output: 1234
}
