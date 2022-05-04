package main

import (
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
)

const (
	CrawlerStateNew = iota
	CrawlerStateRun
	CrawlerStateShutdown
)

type workersPool struct {
	maxWorkers int
	workers    int
	sync.Mutex
}

type PageParser func(r io.Reader) (data interface{}, err error)

//CountingReader is a reader Proxy counting up bytes
type CountingReader struct {
	reader io.Reader
	count  int
}

func (cr *CountingReader) Read(b []byte) (int, error) {
	n, err := cr.reader.Read(b)
	cr.count += n
	return n, err
}

type TaskId int

//CrawlerTask stores the task info
type CrawlerTask struct {
	TaskId     TaskId
	UrlRemains int
	Result     []PageInfo
	resultChan chan []PageInfo
	sync.Mutex
}

//todo add the task timeout
//Crawler defines parameters and the statement for the crawler service
type Crawler struct {
	l *log.Logger
	//urlQueue define the queue of urls to process
	urlQueue UrlsQueue
	//tasks define the active tasks list
	tasks       map[TaskId]*CrawlerTask
	state       int
	getWork     chan bool
	wg          sync.WaitGroup
	workers     workersPool
	taskCounter TaskId
	pageParser  PageParser
	sync.Mutex
}

//NewCrawler creates a new multithreaded http crawler
func NewCrawler(l *log.Logger, maxWorkers int, pageParser PageParser) Crawler {
	return Crawler{
		l:       l,
		state:   CrawlerStateNew,
		getWork: make(chan bool),
		workers: workersPool{
			maxWorkers: maxWorkers,
		},
		tasks:      make(map[TaskId]*CrawlerTask),
		pageParser: pageParser,
	}
}

//run starts the main loop
func (c *Crawler) run() {
	c.l.Printf("Crawler is started, current queue size: %v", c.urlQueue.Len())

	for {
		<-c.getWork
		c.workers.Lock()
		c.Lock()
		state := c.state
		c.Unlock()
		if state == CrawlerStateShutdown && c.workers.workers == 0 {
			c.l.Printf("Crawler is stopped")
			c.wg.Done()
			c.workers.Unlock()
			return
		}

		for c.workers.maxWorkers-c.workers.workers > 0 {
			taskId, url, err := c.urlQueue.Get()
			if err != nil {
				break
			}
			c.workers.workers++
			go c.crawlerWorker(taskId.(*TaskId), url)
		}
		c.workers.Unlock()
	}
}

//crawlerWorker starts crawler worker, should be run in goroutine
func (c *Crawler) crawlerWorker(taskId *TaskId, url string) {

	c.l.Printf("[%v] [stared] %v", *taskId, url)
	pageInfo, err := c.processUrl(url, c.pageParser)

	if err != nil {
		c.l.Printf("Can't process the page: %v: %v", url, err)
		pageInfo = PageInfo{
			Url:  url,
			Meta: PageMeta{Status: 0},
			Data: nil,
		}
	}
	c.l.Printf("[%v] [finished] [%v] %v", *taskId, pageInfo.Meta.Status, url)
	c.workers.Lock()
	c.workers.workers--
	c.workers.Unlock()

	task := c.tasks[*taskId]
	task.Lock()
	task.Result = append(task.Result, pageInfo)
	task.UrlRemains--
	urlRemains := task.UrlRemains
	task.Unlock()
	if urlRemains == 0 {
		c.Lock()
		delete(c.tasks, *taskId)
		c.Unlock()
		task.resultChan <- task.Result
	}

	c.getWork <- true
}

//processUrl download and parse the document by url
//returns downloaded page info (PageInfo) or error
func (c *Crawler) processUrl(url string, parser func(r io.Reader) (interface{}, error)) (result PageInfo, err error) {

	//todo add get timeout
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	result.Url = url
	result.Meta.Status = resp.StatusCode
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		_ = resp.Body.Close()
		return
	}
	ct := resp.Header.Get("Content-Type")
	if ct != "" {
		ct = strings.Split(ct, ";")[0]
		result.Meta.ContentType = &ct
	}

	//count real Content Length
	reader := CountingReader{
		reader: resp.Body,
	}

	d, err := parser(&reader)
	result.Meta.ContentLength = &reader.count
	result.Data = d
	return
}

//Add adds urls to crawler processing
//pageInfo array will be returned to the resultChan channel
func (c *Crawler) Add(urls []string, resultChan chan []PageInfo) (err error) {
	if c.state == CrawlerStateShutdown {
		return errors.New("crawler is shutting down")
	}
	c.Lock()
	c.taskCounter++
	taskId := c.taskCounter

	c.tasks[taskId] = &CrawlerTask{
		TaskId:     taskId,
		UrlRemains: len(urls),
		Result:     []PageInfo{},
		resultChan: resultChan,
	}
	c.Unlock()

	c.urlQueue.Put(&taskId, urls)
	c.getWork <- true
	return
}

//Run starts crawler for tasks listening and processing
func (c *Crawler) Run() {
	if c.state != CrawlerStateNew {
		return
	}
	c.state = CrawlerStateRun
	c.wg.Add(1)
	go c.run()
}

//Stop stops crawler with graceful shutdown
//call Wait() to be sure when all tasks is finished
func (c *Crawler) Stop() {
	if c.state != CrawlerStateRun {
		return
	}
	c.l.Printf("Crawler shutdown is in progress...")
	c.Lock()
	c.state = CrawlerStateShutdown
	c.Unlock()
	c.getWork <- true
}

//Wait waits till crawler graceful shutdown
func (c *Crawler) Wait() {
	c.wg.Wait()
}
