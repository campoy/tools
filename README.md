Tools
=====

Below you'll find a list of tools that I've developed for fun or because they turn out to be useful.

They are mostly all written in Go and distributed under the license specified in the LICENSE file.

## flags

flags provides a set of custom defined flags that you can easily use with the flag package from the standard library.

Learn more about it on this [justforfunc](https://www.youtube.com/watch?v=4D506W1AjeM) video.

[docs](https://godoc.org/github.com/campoy/tools/flags)

## httplog

httplog provides an implementation of http.RoundTripper that logs every single request and response using a given logging function.

[docs](http://godoc.org/github.com/campoy/tools/httplog)

## tree

tree is a very simple implementation of the tree unix command.
This implementation doesn't provide any options as flags.

### Disclaimer

This is not an official Google product (experimental or otherwise), it is just
code that happens to be owned by Google.
