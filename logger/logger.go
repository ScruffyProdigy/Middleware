package logger

import (
	"github.com/HairyMezican/TheRack/rack"
	"net/http"
	"log"
	"io"
	"os"
)

type Logger struct {
	log	*log.Logger
}

func (this Logger) Run(r *http.Request, vars rack.Vars, next rack.Next) (status int, header http.Header, message []byte) {
	vars["Logger"] = this.log
	return next()
}

func Set(out io.Writer,prefix string, flag int) Logger{
	return Logger{log.New(out, prefix, flag)}
}

func Get(vars rack.Vars) interface{} {
	return vars["Logger"]
}

var StandardLogger = Set(os.Stdout,"",log.LstdFlags)