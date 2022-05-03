package main

type PageInfo struct {
	Url  string      `json: "url"`
	Meta PageMeta    `json: "meta"`
	Data interface{} `json: "data"`
}

type PageMeta struct {
	Status 			int		 `json: "status"`
	ContentType		*string  `json: "content-type"`
	ContentLength	*int	 `json: "content-length"`
}

type JsonRequest []string

type JsonResponse []PageInfo