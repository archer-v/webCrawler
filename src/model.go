package main

type PageInfo struct {
	Url  string      `json:"url"`
	Meta PageMeta    `json:"meta"`
	Data interface{} `json:"data,omitempty"`
}

type PageMeta struct {
	Status        int     `json:"status"`
	ContentType   *string `json:"content-type,omitempty"`
	ContentLength *int    `json:"content-length,omitempty"`
}

type JsonRequest []string
