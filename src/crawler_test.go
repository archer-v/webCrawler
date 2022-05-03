package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestCrawler(t *testing.T) {

	pageParser := func (r io.Reader) (data interface{}, err error) {
		return PageParseTagsCounter(r)
	}

	crawler := NewCrawler(log.New(os.Stdout, "SPIDER ", log.LstdFlags), 10, pageParser)
	crawler.Run()

	resultChan := make(chan []PageInfo)
	testCases1 := []string{"https://amazon.com/", "https://google.com", "https://google.com/badpage", "http://www.microsoft.com/", "http://www.microsoft.com/404"}
	testCases2 := []string{"https://facebook.com/", "https://yahooo.com/", "https://naggggg.com"}

	t.Run("Add task 1 to the crawler", func(t *testing.T) {
		err := crawler.Add(testCases1, resultChan)
		if err != nil {
			t.Fatalf("error on adding url to the crawler: %v", err)
		}
	})

	t.Run("Add task 2 to the crawler", func(t *testing.T) {
		err := crawler.Add(testCases2, resultChan)
		if err != nil {
			t.Fatalf("error on adding url to the crawler: %v", err)
		}
	})

	t.Run("Waiting for tasks finished", func(t *testing.T) {
		countDown := 2
		for countDown > 0 {
			pageInfo, ok := <- resultChan
			if !ok {
				t.Fatalf("the channel is closed unexpectedly by crawler")
			}
			countDown--
			for _, p := range pageInfo {
				fmt.Printf("%v %v\n", p.Url, p.Meta.Status)
			}
		}
	})

	crawler.Stop()
	crawler.Wait()
}

func TestProcessUrl(t *testing.T) {

	type testCase struct {
		url 				string
		httpCode 			int
		nilContentLength 	bool
		ContentType     	string
	}

	testCases := []testCase{
		{
			"https://google.com/badpage",
			404,
			true,
			"",
		},
		{
			"https://google.com",
			200,
			false,
			"text/html",
		},
	}

	c := Crawler{}
	for _, testCase := range testCases {
		result, err := c.processUrl(testCase.url, func(r io.Reader) (interface{}, error) {
			_, err := ioutil.ReadAll(r)
			return nil, err
		})
		if err != nil {
			t.Fatalf("Expected error %v", err)
		}
		if result.Url != testCase.url {
			t.Fatalf("Wrong url field in result")
		}
		if result.Meta.Status != testCase.httpCode {
			t.Fatalf("Wrong response status %v for url: %v", result.Meta.Status, testCase.url)
		}
		if (result.Meta.ContentLength == nil) != testCase.nilContentLength  {
			t.Fatalf("Wrong Content-Length for url: %v", testCase.url)
		}

		if result.Meta.ContentType == nil {
			if testCase.ContentType != "" {
				t.Fatalf("Content-Type shouldn't be null for url: %v", testCase.url)
			}
		} else {
			if *result.Meta.ContentType != testCase.ContentType {
				t.Fatalf("Wrong Content-Type (got %v, expected %v) for url: %v", *result.Meta.ContentType, testCase.ContentType, testCase.url)
			}
		}
	}
}