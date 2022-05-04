package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

//WebServerRoute stores uri path and the appropriate handler
type WebServerRoute struct {
	path    string
	handler func(resp http.ResponseWriter, req *http.Request)
}

//WebServer stores
type WebServer struct {
	l      *log.Logger
	server *http.Server
}

//NewWebServer creates the webServer instance
func NewWebServer(bindTo string, port int, l *log.Logger, routes []WebServerRoute) (ws WebServer) {
	var mux http.ServeMux

	ws = WebServer{
		l: l,
		server: &http.Server{
			Addr:    fmt.Sprintf("%s:%d", bindTo, port),
			Handler: &mux,
		},
	}

	for _, route := range routes {
		mux.HandleFunc(route.path, route.handler)
	}
	return
}

//Run starts the webserver in background and returns error on startup errors
func (s *WebServer) Run() error {

	const serviceStartupTimeout = time.Millisecond * 500
	closeCh := make(chan error)
	go func() {
		err := s.server.ListenAndServe()
		if err != nil {
			closeCh <- err
		}
		close(closeCh)
	}()

	select {
	case err := <-closeCh:
		return fmt.Errorf("can't start web server: %v", err)
	case <-time.After(serviceStartupTimeout):

	}
	s.l.Printf("Web server has been started at %s", s.server.Addr)
	return nil
}

//Stop() stops the service with grateful shutdown
//Stop is waiting till all http requests is finished
func (s *WebServer) Stop() {
	ctx := context.Background()
	s.l.Printf("Web server is stopped")
	_ = s.server.Shutdown(ctx)
}

//ReadBody is the helper function reading the body from http.Request to the byte buffer
func (s WebServer) ReadBody(req *http.Request) ([]byte, error) {
	b := bytes.NewBuffer(make([]byte, 0))

	_, err := b.ReadFrom(io.LimitReader(req.Body, req.ContentLength))

	if err != nil {
		s.l.Printf("Can't read request body %v\n", err)
		return nil, errors.New("body read error")
	}
	_ = req.Body.Close()

	return b.Bytes(), nil
}

//ReadBodyAsJSON is the helper function reading the body from http.Request and unmarshalling to JSON
func (s WebServer) ReadBodyAsJSON(req *http.Request, j interface{}) (err error) {

	defer func() {
		if err != nil {
			s.l.Printf("got wrong request: %v", err)
		}
	}()

	if req.Header.Get("Content-type") != "application/json" {
		err = errors.New("wrong content-type (not a json)")
		return
	}

	b, err := s.ReadBody(req)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, j)
	if err != nil {
		err = fmt.Errorf("can't parse json from request body %v", err)
	}
	return
}

//JsonResponse is a helper function writing d struct as JSON to http.ResponseWriter
func (s *WebServer) JsonResponse(d interface{}, res http.ResponseWriter) (err error) {
	defer func() {
		if err != nil {
			s.l.Printf("can't write JSON response: %v", err)
		}
	}()

	b, err := json.Marshal(d)
	if err != nil {
		err = fmt.Errorf("can't prepare JSON response data: %v", err)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	_, err = res.Write(b)

	return
}
