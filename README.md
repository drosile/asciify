# Overview

This is a Go web service that will turn an image into ASCII art.

# Installation

`go get github.com/drosile/asciify`

# Operation

`go run $GOPATH/src/github.com/drosile/asciify`

This will run a web server on local port 9999. This server provides two endpoints:
- `/` will return ASCII art as plain text
- `/json` will return a JSON payload containing the ASCII art under the attribute `image_string`

To use either endpoint, either POST an image file with multipart/form-data using `image_file`, or do a GET with the URL query parameter `image_url=[your_image_url]`. A default image will be used if neither of these options is used.

The default width of the generate ASCII art is 100 characters. To change that, pass a `width` query parameter as an integer.

# Known Issues

If you pass something that is not an image to this service, it will try to treat it like an image, but it will likely result in a failure.

# Dependencies

This project makes use of [https://github.com/stdupp/goasciiart](goasciiart) code, but since it is a program, not a package, I have included it in this repo. This code does the following:
- scale an image to a given width (this may also be understood as the width in ASCII characters)--this includes changing the proportions to match monospace characters
- read the grayscale value of each pixel in the re-scaled image
- print an ASCII character for each pixel, with higher-grayscale (whiter) pixels being represented by "denser" ASCII characters

This code, in turn, depends on [https://github.com/nfnt/resize](resize). Alternative image libraries may be used if necessary for the resizing.
