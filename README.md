# webCrawler
Simple boilerplate multithreaded web crawler written in golang.

The service receives list of urls (JSON array) from a POST request, loads all urls content in concurrent threads, and returns the pages metadata as a combined JSON response 

Returned response is a array of records containing:
 - Url
 - Response status of this url
 - Content length
 - Content type
 - Additional page parsing data (in this example it just counts html tags)

#### Options

Service accepts two env configuration parameters 
 - WEBCRAWLER_HTTPPORT (default 8001)
 - WEBCRAWLER_WORKERS (default 10)


#### Request example

```js
[
  "http://www.example1.com/",
  "http://www.example2.com/",
  // ...
]
```

#### Response example

```js
[
  {
    "url": "http://www.example1.com/",
    "meta": {
      "status": 200,
      "content-type": "text/html",
      "content-length": 605
    },
    "data": {
            "elements": [
                {
                    "tag-name": "meta",
                    "count": 4
                },
                {
                    "tag-name": "body",
                    "count": 1
                },
                {
                    "tag-name": "html",
                    "count": 1
                },
                {
                    "tag-name": "div",
                    "count": 12
                },
                //...                
            ]
        }
  },
  {
    "url": "http://www.example2.com/",
    "meta": {
      "status": 404,
    },
  },
  // ...
]
```

#### Page parser

You can implement your own page parser by implementing your own parsing handler
```
func(r io.Reader) (data interface{}, err error) {

}
```
see ``func PageParseTagsCounter(r io.Reader) (data *PageData, err error)`` as an example

#### Startup

```bash
$ git clone https://github.com/mandalorian-one/webCrawler.git
$ make build
$ docker-compose run --rm app /build/app
```
