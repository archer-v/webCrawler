package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var gitTag, gitCommit, gitBranch, build, version string
var webServer WebServer
var crawler Crawler

var httpLogger, crawlerLogger *log.Logger

//ProcessWebRequest is called by webserver and processes the web request
func ProcessWebRequest(res http.ResponseWriter, req *http.Request) {

	if req.Method != http.MethodPost {
		http.Error(res, "only POST requests is allowed", http.StatusBadRequest)
		return
	}

	var urls JsonRequest
	err := webServer.ReadBodyAsJSON(req, &urls)

	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	resultChan := make(chan []PageInfo)

	if len(urls) == 0 {
		http.Error(res, "empty request", http.StatusBadRequest)
		return
	}

	httpLogger.Printf("Got a new request for %v urls", len(urls))

	err = crawler.Add(urls, resultChan)
	if err != nil {
		//crawler is in shutting down mode
		http.Error(res, err.Error(), http.StatusForbidden)
		return
	}

	data, ok := <-resultChan
	if !ok {
		e := "the channel is closed unexpectedly by crawler"
		crawlerLogger.Printf(e)
		http.Error(res, e, http.StatusInternalServerError)
		return
	}

	err = webServer.JsonResponse(data, res)

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func main() {

	if build == "" {
		build = "DEV"
		version = fmt.Sprintf("version: DEV, build: %v", build)
	} else {
		version = fmt.Sprintf("version: %v-%v-%v, build: %v", gitTag, gitBranch, gitCommit, build)
	}

	fmt.Printf("WebCrawler %v\n", version)

	c, err := InitConfig()

	if err != nil {
		panic(err)
	}

	confPrint, _ := json.MarshalIndent(c, "", "\t")
	fmt.Println("Startup configuration: ")
	fmt.Println(string(confPrint))

	httpLogger = log.New(os.Stdout, "HTTP ", log.LstdFlags)
	crawlerLogger = log.New(os.Stdout, "CRAWLER ", log.LstdFlags)

	webServer = NewWebServer(
		"",
		c.HTTPPort,
		httpLogger,
		[]WebServerRoute{{"/", ProcessWebRequest}})

	crawler = NewCrawler(
		crawlerLogger,
		c.Workers,
		func(r io.Reader) (data interface{}, err error) {
			return PageParseTagsCounter(r)
		})

	err = webServer.Run()
	if err != nil {
		fmt.Printf("Can't start the web server: %v", err.Error())
		return
	}
	crawler.Run()

	//wait for signal and perform the graceful shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-sig

	crawler.Stop()
	crawler.Wait()
	webServer.Stop()

	fmt.Printf("WebCrawler exit\n")
}
