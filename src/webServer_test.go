package main

import (
	"fmt"
	"github.com/phayes/freeport"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"
)

func TestWebServer(t *testing.T) {

	port, err := freeport.GetFreePort()
	if err != nil {
		t.Fatal("Can't get a free tcp port")
	}
	pattern := "SUCCESS"

	WebServer := NewWebServer(
		"",
		port,
		log.New(os.Stdout, "HTTP ", log.LstdFlags),
		[]WebServerRoute{
		{"/", func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(http.StatusOK)
			_, err := res.Write([]byte(pattern))
			if err != nil {
				t.Fatalf("can't prepare response")
			}
		}},
	})

	t.Run("Webserver start", func(t *testing.T) {
		err := WebServer.Run()
		if err != nil {
			t.Fatalf("can't start webserver: %v", err)
		}
	})

	defer func() {
		WebServer.Stop()
	}()

	t.Run("Webserver serve the request", func(t *testing.T) {
		if err != nil {
			t.Fatalf("can't start webserver: %v", err)
		}

		resp, err := http.Get(fmt.Sprintf("http://localhost:%v/", port))
		if err != nil {
			t.Fatalf("can't get webserver response: %v", err)
		}

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("can't read webserver response: %v", err)
		}

		if string(b) != pattern {
			t.Fatalf("wrong answer")
		}
	})
}