# webCrawler
Simple multithreaded web crawler written in golang.

The service receives list of urls (JSON array) from a POST request, loads all urls content in concurrent threads, and returns the pages metadata as a combined JSON response 

Returned response is a array of records containing:
 - Url
 - Response status of this url
 - Content length
 - Content type
 - Additional page parsing data (in this example it just counts pages html tags)

Service accepts two env configuration parameters 
 - WEBCRAWLER_HTTPPORT (default 8001)
 - WEBCRAWLER_WORKERS (default 10)

POST request example:
```js
[
  "http://www.example1.com/",
  "http://www.example2.com/",
  // ...
]
```
Response example:
```js
[
  {
    "url": "http://www.example1.com/",
    "meta": {
      "status": 200,
      "content-type": "text\/html",
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
