package redirecthttps

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/acoshift/header"
	"github.com/acoshift/middleware"
	"github.com/stretchr/testify/assert"
)

// TestNew asserts no fundamental issues with default configuration
func TestNew(t *testing.T) {
	mw := New(Config{})

	ts := httptest.NewServer(mw(getTestHandler()))
	defer ts.Close()

	request, err := http.NewRequest("HEAD", ts.URL, strings.NewReader(""))
	assert.NoError(t, err)

	client := http.Client{CheckRedirect: checkRedirect}

	res, err := client.Do(request)
	assert.Error(t, err)
	assert.Equal(t, http.StatusMovedPermanently, res.StatusCode)

	e := "https" + strings.Replace(ts.URL, "http", "", 4)
	assert.True(t, strings.Contains(err.Error(), e))

}

// TestNew_HTTP asserts that http requests are redirected to https
func TestNew_HTTP(t *testing.T) {
	mw := New(Config{})

	ts := httptest.NewServer(mw(getTestHandler()))
	defer ts.Close()

	request, err := http.NewRequest("HEAD", ts.URL, strings.NewReader(""))
	assert.NoError(t, err)

	client := http.Client{CheckRedirect: checkRedirect}

	res, err := client.Do(request)
	assert.Error(t, err)
	assert.Equal(t, http.StatusMovedPermanently, res.StatusCode)

	e := "https" + strings.Replace(ts.URL, "http", "", 4)
	assert.True(t, strings.Contains(err.Error(), e))
}

// TestNew_HTTPS asserts that https requests are served without redirect
func TestNew_HTTPS(t *testing.T) {
	mw := New(Config{})

	ts := httptest.NewTLSServer(mw(getTestHandler()))
	defer ts.Close()

	client, err := getTLSClient(ts)
	assert.NoError(t, err)
	client.CheckRedirect = checkRedirect

	res, err := client.Head(ts.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

// TestNew_OnlyProxyMode_HTTP asserts that requests are redirected when the client connects to proxy with http
func TestNew_OnlyProxyMode_HTTP(t *testing.T) {
	mw := New(Config{
		Mode: OnlyProxy,
	})

	ts := httptest.NewServer(mw(getTestHandler()))
	defer ts.Close()

	request, err := http.NewRequest("HEAD", ts.URL, strings.NewReader(""))
	assert.NoError(t, err)

	request.Header.Add(header.XForwardedProto, "http")
	client := http.Client{CheckRedirect: checkRedirect}

	res, err := client.Do(request)
	assert.Error(t, err)
	assert.Equal(t, http.StatusMovedPermanently, res.StatusCode)
}

// TestNew_OnlyProxyMode asserts that requests are not redirected when the client connects to proxy with http
func TestNew_OnlyProxyMode_HTTPS(t *testing.T) {
	mw := New(Config{
		Mode: OnlyProxy,
	})

	ts := httptest.NewServer(mw(getTestHandler()))
	defer ts.Close()

	request, err := http.NewRequest("HEAD", ts.URL, strings.NewReader(""))
	assert.NoError(t, err)

	request.Header.Add(header.XForwardedProto, "https")
	client := http.Client{CheckRedirect: checkRedirect}

	res, err := client.Do(request)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

// TestNew_OnlyConnectionMode asserts that requests made without TLS are redirected
func TestNew_OnlyConnectionMode_HTTP(t *testing.T) {
	mw := New(Config{Mode: OnlyConnectionState})

	ts := httptest.NewServer(mw(getTestHandler()))
	defer ts.Close()

	request, err := http.NewRequest("HEAD", ts.URL, strings.NewReader(""))
	assert.NoError(t, err)

	client := http.Client{CheckRedirect: checkRedirect}

	res, err := client.Do(request)
	assert.Error(t, err)
	assert.Equal(t, http.StatusMovedPermanently, res.StatusCode)

	e := "https" + strings.Replace(ts.URL, "http", "", 4)
	assert.True(t, strings.Contains(err.Error(), e))
}

// TestNew_OnlyConnectionMode asserts that requests made with TLS are not redirected
func TestNew_OnlyConnectionMode_HTTPS(t *testing.T) {
	mw := New(Config{Mode: OnlyConnectionState})

	ts := httptest.NewTLSServer(mw(getTestHandler()))
	defer ts.Close()

	client, err := getTLSClient(ts)
	assert.NoError(t, err)
	client.CheckRedirect = checkRedirect

	res, err := client.Head(ts.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

// TestNew_AllMode asserts that requests made without TLS or through a proxy are redirected
func TestNew_AllMode(t *testing.T) {
	mw := New(Config{Mode: All})

	ts := httptest.NewServer(mw(getTestHandler()))
	defer ts.Close()

	request, err := http.NewRequest("HEAD", ts.URL, strings.NewReader(""))
	assert.NoError(t, err)

	request.Header.Add(header.XForwardedProto, "http")
	client := http.Client{CheckRedirect: checkRedirect}

	res, err := client.Do(request)
	assert.Error(t, err)
	assert.Equal(t, http.StatusMovedPermanently, res.StatusCode)

	e := "https" + strings.Replace(ts.URL, "http", "", 4)
	assert.True(t, strings.Contains(err.Error(), e))
}

// TestNew_SkipHTTP asserts that http requests are not redirected
func TestNew_SkipHTTP(t *testing.T) {
	mw := New(Config{Skipper: middleware.SkipHTTP})

	ts := httptest.NewServer(mw(getTestHandler()))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

// TestNew_SkipHTTPS asserts that https requests are not redirected
func TestNew_SkipHTTPS(t *testing.T) {
	mw := New(Config{Skipper: middleware.SkipHTTPS})

	ts := httptest.NewTLSServer(mw(getTestHandler()))
	defer ts.Close()

	client, err := getTLSClient(ts)
	assert.NoError(t, err)
	client.CheckRedirect = checkRedirect

	res, err := client.Head(ts.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

// getTLSClient configures a client for a TLS test and returns it.
//
// In go version 1.9, this is not necessary and the API has been updated to do this boilerplate.
//
// Reference: https://github.com/golang/go/issues/18411
func getTLSClient(ts *httptest.Server) (http.Client, error) {
	cert, err := x509.ParseCertificate(ts.TLS.Certificates[0].Certificate[0])
	if err != nil {
		return http.Client{}, err
	}

	certpool := x509.NewCertPool()
	certpool.AddCert(cert)

	return http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: certpool,
			},
		},
	}, nil
}

func getTestHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "test")
	})
}

func checkRedirect(req *http.Request, via []*http.Request) error {
	return errors.New(fmt.Sprintf("redirected to %s", req.URL))
}
