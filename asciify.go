package main

import (
  "./goasciiart"

  "bytes"
  "fmt"
  "html"
  "image"
  "io/ioutil"
  "net/http"
  "strconv"
)

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

  // TODO: add image upload support
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

func main() {
  const (
    http_port = ":9999"
  )
  http.HandleFunc("/", asciify_handler)
  http.ListenAndServe(http_port, nil)
}
