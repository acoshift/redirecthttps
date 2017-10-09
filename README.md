<div align="center">
  <img src="logo@2x.png" alt="Logo" width='45%' />

  <h2>RedirectHTTPS: HTTPS Redirection Middleware</h2>
  <p>A minimalistic middleware that redirects all network traffic from the insecure HTTP protocol to the HTTPS transport, all written in Go<p>
</div>

<br />

<div align="center">

[![Go Report Card](https://goreportcard.com/badge/github.com/acoshift/redirecthttps?style=flat-square)](https://goreportcard.com/report/github.com/acoshift/redirecthttps)
[![build status](https://img.shields.io/travis/acoshift/redirecthttps/master.svg?style=flat-square)](https://travis-ci.org/kataras/iris)
[![github issues](https://img.shields.io/github/issues/acoshift/redirecthttps.svg?style=flat-square)](https://github.com/acoshift/redirecthttps/issues?q=is%3Aopen+is%3Aissue)
[![release](https://img.shields.io/github/release/acoshift/redirecthttps.svg?style=flat-square)](https://github.com/acoshift/redirecthttps/releases)
[![chat](https://img.shields.io/badge/community-%20chat-00BCD4.svg?style=flat-square)](https://gitter.im/acoshift)
[![license](https://img.shields.io/github/license/acoshift/redirecthttps.svg?style=flat-square)]()

</div>

<br />

### Installation

`go get github.com/acoshift/redirecthttps`

### Usage Example

This middleware can be applied in the middleware chain, as follows:

```
middlewares := middleware.Chain(
  // We can put our redirection middleware right here.
  redirecthttps.New(redirecthttps.Config{
    // The redirection mode can be specified below.
    Mode: redirecthttps.OnlyProxy
  }),
)

mux.Handle("/", middlewares(app.Handler()))
```

### Configuration: Redirection Modes

There are three redirection modes, which can be specified before instantiating the middleware.

First, the `OnlyConnectionState` Mode, which only checks the connection state from request.TLS, and perform the redirection if TLS is not present in the request.

Second, the `OnlyProxy` Mode, which only checks the X-Forwarded-Proto header from the request in order to determine if it's using plain HTTP or not. If so, perform the redirection.

Finally, the `All` Mode, which checks both the X-Forwarded-Proto Header AND the request.TLS variable.

### Contribution

If you found an issue in this library, please file file an issue at https://github.com/acoshift/redirecthttps/issues.

If you wanted to help improve this middleware, feel free to fork this project and submit a pull request through GitHub. Thanks!

### License

MIT License

Copyright (c) 2017 Thanatat Tamtan

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
