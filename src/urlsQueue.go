package main

import (
	"errors"
	"sync"
)

type task struct {
	url string
	id  interface{}
}

// UrlsQueue define the queue of urls to process
type UrlsQueue struct {
	sync.Mutex
	q  	[]task
}

//Put saves urls and the related taskId to the queue
func (q *UrlsQueue) Put(taskId interface{}, urls []string) {
	q.Lock()
	for _, url := range urls {
		q.q = append(q.q, task{url, taskId})
	}
	q.Unlock()
	return
}

//Get returns the next url from the queue and the taskId
func (q *UrlsQueue) Get() (taskId interface{}, url string, err error){
	q.Lock()
	if len(q.q) == 0 {
		err = errors.New("empty")
	} else {
		url = q.q[0].url
		taskId = q.q[0].id
		q.q = q.q[1:]
	}
	q.Unlock()
	return
}

//Len returns the size of the queue
func (q *UrlsQueue) Len() (l int) {
	q.Lock()
	l = len(q.q)
	q.Unlock()
	return
}