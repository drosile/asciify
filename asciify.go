package main

import (
  "fmt"
  "net/http"
)

func asciify_handler(w http.ResponseWriter, r *http.Request) {
  const (
    test_image = "http://i.imgur.com/ewD2qvQ.jpg"
  )
  resp, err := http.Get(test_image)
  if err != nil {
    fmt.Fprintf(w, "request failed")
    return
  }
  fmt.Fprintf(w, "request succeeded %s", resp)
}

func main() {
  const (
    http_port = ":9999"
  )
  http.HandleFunc("/", asciify_handler)
  http.ListenAndServe(http_port, nil)
}
