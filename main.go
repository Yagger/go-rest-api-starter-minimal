package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime/debug"

	conf "hashsnap-event-server/conf"
	context "hashsnap-event-server/context"
	utils "hashsnap-event-server/utils"

	"github.com/julienschmidt/httprouter"
)

var (
	ctx *context.Context
)

func makeHandler(fn func(http.ResponseWriter, *http.Request, httprouter.Params, *context.Context)) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("%s: %s", r, debug.Stack())
				if r == 401 {
					w.WriteHeader(http.StatusUnauthorized)
				} else if r == 404 {
					w.WriteHeader(http.StatusNotFound)
				} else {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(w, fmt.Sprint(r))
					fmt.Println(r)
				}
			}
		}()
		fn(w, r, ps, ctx)
	}
}

func main() {
	ctx = &context.Context{}
	ctx.Init()

	var router = httprouter.New()
	router.GET("/version", makeHandler(versionHandler))
	log.Printf("Listening on %s", conf.Port)
	log.Fatal(http.ListenAndServe(":"+conf.Port, &Server{router}))
}

// Server struct
type Server struct {
	r *httprouter.Router
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("Access %s %s from %s | %s | %s | %s | %s | %s | %+v", r.Method, r.URL.Path, r.Header.Get("Origin"), r.Host, r.Proto, r.RemoteAddr, r.RequestURI, r.UserAgent(), r.Header)
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Headers", "X-ACCESS_TOKEN, Access-Control-Allow-Origin, Authorization, Origin, x-requested-with, Content-Type, Content-Range, Content-Disposition, Content-Description")
		w.Header().Set("Access-Control-Allow-Methods", "POST, PUT, GET, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Max-Age", "3600")
	}
	s.r.ServeHTTP(w, r)
}

func versionHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx *context.Context) {
	data, err := ioutil.ReadFile("version")
	utils.CheckErr(err, "versionHandler", "Unable to read version file")
	fmt.Fprintf(w, "Version: %s", data)
}
