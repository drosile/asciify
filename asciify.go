package main

import (
  "fmt"
  "html"
  "io/ioutil"
  "net/http"
)

func read_image_data(r *http.Request) (image_data []byte, err error) {
  const (
    test_image_url = "http://i.imgur.com/ewD2qvQ.jpg"
  )
  r.ParseForm()
  image_url := html.UnescapeString(r.FormValue("image_url"))

  // TODO: add image upload support
  if image_data == nil {
    if image_url == "" {
      image_url = test_image_url
    }

    resp, image_err := http.Get(image_url)
    // FIXME: there must be a more idiomatic way to handle errors
    if image_err != nil {
      err = image_err
      return
    }
    i, read_err := ioutil.ReadAll(resp.Body)
    if read_err != nil {
      err = read_err
      return
    }
    image_data = i
  }
  return
}

func asciify_handler(w http.ResponseWriter, r *http.Request) {
  image_data, err := read_image_data(r)
  if err != nil {
    fmt.Fprintf(w, "request failed")
    return
  }

  fmt.Fprintf(w, "%s", image_data)
}

func main() {
  const (
    http_port = ":9999"
  )
  http.HandleFunc("/", asciify_handler)
  http.ListenAndServe(http_port, nil)
}
