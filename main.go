package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	Conf "go-webapp-starter/conf"
	HandlerContext "go-webapp-starter/context"

	Handlers "go-webapp-starter/handlers"

	"github.com/julienschmidt/httprouter"
)

var (
	hctx *HandlerContext.Context
)

func makeHandler(fn func(http.ResponseWriter, *http.Request, httprouter.Params, *HandlerContext.Context)) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		defer func() {
			if r := recover(); r != nil {
				if r == 401 {
					w.WriteHeader(http.StatusUnauthorized)
				} else if r == 404 {
					w.WriteHeader(http.StatusNotFound)
				} else {
					log.Printf("%s: %s", r, debug.Stack())
					w.WriteHeader(http.StatusInternalServerError)
					if Conf.Debug {
						fmt.Fprintf(w, fmt.Sprint(r))
					} else {
						fmt.Fprintf(w, "Unexpecter error happened")
					}
					fmt.Println(r)
				}
			}
		}()
		fn(w, r, ps, hctx)
	}
}

func main() {
	hctx = &HandlerContext.Context{}
	hctx.Init()

	var router = httprouter.New()

	router.GET("/version", makeHandler(Handlers.VersionHandler))
	router.POST("/login", makeHandler(Handlers.LoginHandler))
	router.GET("/logout", makeHandler(Handlers.LogoutHandler))
	router.GET("/whoami", makeHandler(Handlers.WhoamiHandler))
	router.GET("/accounts", makeHandler(Handlers.ListAccountHandler))
	router.PUT("/account", makeHandler(Handlers.CreateAccountHandler))
	// router.GET("/account/:account_id", makeHandler(HandlersReadAccountHandler))
	// router.POST("/account/:account_id", makeHandler(Handlers.UpdateAccountHandler))

	log.Printf("Listening on %s", Conf.Port)
	log.Fatal(http.ListenAndServe(":"+Conf.Port, &Server{router}))
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
