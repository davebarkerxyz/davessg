# davessg - Dave's Static Site Generator

Because none of the hundreds of other static site generators fit my needs just perfectly.

Ultra simple, minimal, and probably missing most of the features *you* need.

Dave's Static Site Generator is written in Go and uses the [Goldmark](https://github.com/yuin/goldmark) library for Markdown parsing.


## License

davessg is MIT licensed.


## Building

Clone the repo and run `go build davessg.go`. Nice and easy.


## Usage

`davessg -help` is pretty self explanatory.

```
% ./davessg -help
Usage of ./davessg:
  -base-url string
        Base URL (default "/")
  -bind-addr string
        Listen address for web server (use with -serve) (default "localhost:8009")
  -build-dir string
        Build directory (created if necessary) (default "build/")
  -force
        Overwrite existing build files, if necessary
  -serve
        Start web server
  -source-dir string
        Source content dir (default "content/")
  -verbose
        Verbose output
```

- Drop your content files in *content/*.
- Output goes in *build/* by default.
- Markdown (.md) and HTML (.htm and .html) files have the template in *templates/index.html* applied, with pretty-paths (e.g. *content/mypage.md* builds to *build/mypage/index.html* so it can be accessed at *https://mysite/mypage*.
- Page bundles (extra files like images bundled with pages are supported) - *content/mypage/header.png* will be copied to *build/mypage/header.png*.
- Files in *templates/static/* are copied as-is to *build/static/*.
- If a file in the build dir is newer than the content file, it's ignored (unless the `-force` arg is passed).

