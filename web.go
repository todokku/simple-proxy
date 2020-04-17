package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var authToken string

func main() {
	token, err := ioutil.ReadFile(os.Getenv("TOKEN_PATH"))
	if err != nil {
		os.Stderr.WriteString("Failed to load token from disk\n")
		os.Exit(1)
	}
	authToken = strings.TrimSpace(string(token))

	http.HandleFunc("/", proxy)
	http.ListenAndServe(":8080", nil)
}

func proxy(w http.ResponseWriter, r *http.Request) {
	tokens := r.Header["Authorization"]
	if len(tokens) < 1 {
		http.Error(w, "authorization is missing", http.StatusUnauthorized)
		return
	}
	tokenHeaderParts := strings.Split(tokens[0], " ")
	if len(tokenHeaderParts) < 2 {
		http.Error(w, "authorization is invalid", http.StatusUnauthorized)
		return
	}
	if tokenHeaderParts[1] != authToken {
		http.Error(w, "token is not authorized", http.StatusUnauthorized)
		return
	}

	urls := r.URL.Query()["url"]

	if len(urls) < 1 {
		http.Error(w, "missing url param", http.StatusBadRequest)
		return
	}

	url, err := url.Parse(urls[0])
	if err != nil {
		http.Error(w, "invalid url param", http.StatusBadRequest)
		return
	}

	newRequest := &http.Request{
		URL:    url,
		Header: r.Header,
	}

	client := &http.Client{}

	resp, err := client.Do(newRequest)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get url %v", err), http.StatusBadGateway)
		return
	}
	fmt.Printf("%v", resp)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))

	for k, v := range resp.Header {
		w.Header().Add(k, v[0])
	}
	_, err = w.Write(body)
	if err != nil {
		http.Error(w, "failed to write response", http.StatusInternalServerError)
		return
	}
}
