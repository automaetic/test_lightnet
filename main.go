package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func main() {
	routes := Route{
		Sum: "/calculator.sum",
		Div: "/calculator.div",
		Mul: "/calculator.mul",
		Sub: "/calculator.sub",
	}
	go func() { initServer(routes) }()
	initProxy(routes)
}

func initServer(routes Route) {
	port := 8081
	http.HandleFunc(routes.Sum, func(w http.ResponseWriter, req *http.Request) {
		params := getParams(req.Body)
		output := params.A + params.B
		fmt.Fprintf(w, strconv.Itoa(output))
	})

	http.HandleFunc(routes.Sub, func(w http.ResponseWriter, req *http.Request) {
		params := getParams(req.Body)
		output := params.A - params.B
		fmt.Fprintf(w, strconv.Itoa(output))
	})

	http.HandleFunc(routes.Mul, func(w http.ResponseWriter, req *http.Request) {
		params := getParams(req.Body)
		output := params.A * params.B
		fmt.Fprintf(w, strconv.Itoa(output))
	})

	http.HandleFunc(routes.Div, func(w http.ResponseWriter, req *http.Request) {
		params := getParams(req.Body)
		output := params.A / params.B
		fmt.Fprintf(w, strconv.Itoa(output))
	})
	log.Println(fmt.Printf("Listing for requests at http://localhost:%d\n", port))
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func initProxy(routes Route) {
	port := 8080
	proxyScheme := "http"
	proxyHost := "localhost:8081"
	mux := http.NewServeMux()

	forward := func(w http.ResponseWriter, req *http.Request) {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		req.Body = ioutil.NopCloser(bytes.NewReader(body))
		url := fmt.Sprintf("%s://%s%s", proxyScheme, proxyHost, req.RequestURI)
		println("Forward Req to url: ", url)
		proxyReq, err := http.NewRequest(req.Method, url, bytes.NewReader(body))
		proxyReq.Header = make(http.Header)
		for h, val := range req.Header {
			proxyReq.Header[h] = val
		}

		httpClient := http.Client{}
		resp, err := httpClient.Do(proxyReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		w.Write(bodyBytes)
		defer resp.Body.Close()
	}

	mux.HandleFunc(routes.Sum, forward)
	mux.HandleFunc(routes.Div, forward)
	mux.HandleFunc(routes.Mul, forward)
	mux.HandleFunc(routes.Sub, forward)
	log.Println(fmt.Printf("Proxy Listing @ at http://localhost:%d\n", port))
	http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
}

func getParams(body io.ReadCloser) Params {
	var params Params
	err := json.NewDecoder(body).Decode(&params)
	if err != nil {
		panic(err.Error())
	}
	return params
}

type Params struct {
	A int `json:"a"`
	B int `json:"b"`
}
type Route struct {
	Sum string
	Div string
	Mul string
	Sub string
}
