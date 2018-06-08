package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/alinz/go-callbag"
)

func downloadPage(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func main() {
	callbag.Pipe(
		// provide list of URLs
		callbag.FromValues(
			"https://golang.org",
			"https://golang.org/pkg",
		),
		// download url and return either an error or
		// content
		callbag.Map(func(val interface{}) interface{} {
			content, err := downloadPage(val.(string))
			if err != nil {
				return err
			}

			return content
		}),
		// filter those things that return an error
		callbag.Filter(func(val interface{}) bool {
			switch val.(type) {
			case []byte:
				return true
			default:
				return false
			}
		}),
		// print the size of downloaded content
		callbag.ForEach(func(val interface{}) {
			fmt.Println("content size: ", len(val.([]byte)))
		}),
	)
}
