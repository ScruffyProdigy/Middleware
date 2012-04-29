/*
	session wraps gorilla sessions within a Rack Middleware framework
*/
package sessioner

import (
	"code.google.com/p/gorilla/sessions"
	"net/http"
)

var store = sessions.NewCookieStore([]byte("Go Game Lobby!"))

/*
	The Session is the interface exposed to the rest of the program
*/
type Session interface {
	Set(k, v interface{})            // Set will set a session variable
	Get(k interface{}) interface{}   //Get will obtain the result of a session variable
	Clear(k interface{}) interface{} //Clear will obstain the result of a session variable, and then delete it from the session value
}

type session struct {
	sess *sessions.Session
	r    *http.Request
}

func (this *session) Set(k, v interface{}) {
	this.sess.Values[k] = v
}

func (this *session) Get(k interface{}) interface{} {
	return this.sess.Values[k]
}

func (this *session) Clear(k interface{}) (result interface{}) {
	result = this.sess.Values[k]
	delete(this.sess.Values, k)
	return
}

func (this *session) save(w http.ResponseWriter) {
	err := this.sess.Save(this.r, w)
	if err != nil {
		panic(err)
	}
}

func get(r *http.Request) *session {
	sess := new(session)
	sess.sess, _ = store.Get(r, "session")
	sess.r = r
	return sess
}
