# WebCrawler
Simple multithreaded web crawler written in golang.

The service gets list of urls (JSON array) as a POST request, loads all urls content in concurrent threads, and returns the pages metadata as a combined JSON response 

Returned response is a array of records containing:
 - Url
 - Response status of this url
 - Content length
 - Content type
 - Additional page parsing data (in this example it just counts pages html tags)

Service accepts two env configuration parameters 
 - WEBCRAWLER_HTTPPORT (default 8001)
 - WEBCRAWLER_WORKERS (default 10)


