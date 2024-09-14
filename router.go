package main

import (
	"log"
	"net/http"
)

type APIServer struct {
	addr string
}

func NewAPIServer(addr string) *APIServer {
	return &APIServer{
		addr: addr,
	}
}

func (a *APIServer) Start() error {

	router := http.NewServeMux()

	router.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	router.Handle("/api/v1/", http.StripPrefix("/api/v1", MdwRequireAuthMiddleware(apiV1())))
	router.Handle("/api/v2/", http.StripPrefix("/api/v2", MdwChainRequestResponseLogAuth(apiV2())))

	server := &http.Server{
		Addr:    a.addr,
		Handler: router,
	}

	log.Printf("Starting server on %s", a.addr)
	return server.ListenAndServe()
}

func apiV1() http.Handler {

	v1 := http.NewServeMux()

	v1.HandleFunc("GET /user/name", handleApiV1GetUser)
	v1.HandleFunc("POST /user/name", handleApiV1PostUser)
	v1.Handle("GET /user/id/{userID}", MdwRequireSuperUserMiddleware(http.HandlerFunc(handleApiV1GetUserID)))

	return v1
}

func apiV2() http.Handler {

	v2 := http.NewServeMux()

	v2.HandleFunc("GET /user/name", handleApiV2GetUser)
	v2.HandleFunc("GET /version", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("API v2"))
	})

	return v2
}
