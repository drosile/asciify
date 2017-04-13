package main

import (
  "./goasciiart"

  "bytes"
  "encoding/json"
  "fmt"
  "html"
  "image"
  "io/ioutil"
  "net/http"
  "strconv"
)

type AsciiAPIResponse struct {
  ImageString string `json:"image_string"`
}

func read_image_data(r *http.Request) (image_data []byte, width int, err error) {
  const (
    default_width = 100
    test_image_url = "http://i.imgur.com/ewD2qvQ.jpg"
  )
  r.ParseForm()
  image_url := html.UnescapeString(r.FormValue("image_url"))
  i, e := strconv.Atoi(r.FormValue("width"))
  width = i
  if e != nil {
    width = default_width
  }

  if width == 0 {
    width = default_width
  }

  // FIXME: these variable names are getting confusing, and I'm ignoring errors
  if r.Method == "POST" {
    img_data, _, _ := r.FormFile("image_file")
    img, _ := ioutil.ReadAll(img_data)
    image_data = img
  }

  if image_data == nil {
    if image_url == "" {
      image_url = test_image_url
    }

    // FIXME: there must be a more idiomatic way to handle errors
    resp, image_err := http.Get(image_url)
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

func convert_image_to_ascii(image_data []byte, width int) ([]byte) {
  data, _, err := image.Decode(bytes.NewReader(image_data))
  if err != nil {
    return nil
  }
  return goasciiart.Convert2Ascii(goasciiart.ScaleImage(data, width))
}

func asciify_handler(w http.ResponseWriter, r *http.Request) {
  image_data, width, err := read_image_data(r)
  if err != nil {
    fmt.Fprintf(w, "request failed")
    return
  }

  fmt.Fprintf(w, "%s", convert_image_to_ascii(image_data, width))
}

func json_asciify_handler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  image_data, width, err := read_image_data(r)
  // FIXME: Error message should be JSON
  if err != nil {
    fmt.Fprintf(w, "request failed")
    return
  }
  image_string := string(convert_image_to_ascii(image_data, width))
  p := AsciiAPIResponse{ImageString: image_string}
  json.NewEncoder(w).Encode(p)
}

func main() {
  const (
    http_port = ":9999"
  )
  http.HandleFunc("/", asciify_handler)
  http.HandleFunc("/json", json_asciify_handler)
  http.ListenAndServe(http_port, nil)
}
