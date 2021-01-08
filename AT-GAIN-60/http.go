package atgain60

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
)

var _ http.RoundTripper = &transport{}

// transport is an http.RoundTripper that is meant for use with this atlona amp. DO NOT USE FOR GENERAL HTTP REQUESTS. It handles the random HTTP/0.9 and HTTP/1.0 responses that the amp responds with. HTTP/1.0 requests are parsed very simply, likely incorrectly for a lot of HTTP requests, and does not set all fields on the http.Response. HTTP/0.9 'spec' was found here https://www.w3.org/Protocols/HTTP/AsImplemented.html.
type transport struct {
	dialer net.Dialer
}

// RoundTrip implements the http.RoundTripper interface
func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	port := req.URL.Port()
	if port == "" {
		port = "80"
	}

	conn, err := t.dialer.DialContext(req.Context(), "tcp", req.URL.Hostname()+":"+port)
	if err != nil {
		return nil, fmt.Errorf("unable to dial: %w", err)
	}
	defer conn.Close()

	deadline, _ := req.Context().Deadline()

	if err := conn.SetDeadline(deadline); err != nil {
		return nil, fmt.Errorf("unable to set connection deadline: %w", err)
	}

	// write the http request line
	// we're gonna send HTTP/0.9, even though we don't need to, just in case
	// it encourages the amp to send HTTP/0.9 responses.
	line := []byte(fmt.Sprintf("%s %s HTTP/0.9\r\n", req.Method, req.URL.Path))

	n, err := conn.Write(line)
	switch {
	case err != nil:
		return nil, fmt.Errorf("unable to write request line: %w", err)
	case n != len(line):
		return nil, fmt.Errorf("unable to write request line: wrote %v/%v bytes", n, len(line))
	}

	// write the final newline
	line = []byte("\r\n")

	n, err = conn.Write(line)
	switch {
	case err != nil:
		return nil, fmt.Errorf("unable to write header: %w", err)
	case n != len(line):
		return nil, fmt.Errorf("unable to write header: wrote %v/%v bytes", n, len(line))
	}

	// read the whole response
	body, err := ioutil.ReadAll(conn)
	if err != nil {
		return nil, fmt.Errorf("unable to read response: %w", err)
	}

	resp := &http.Response{}

	// if first line contains HTTP/, then parse the response as a HTTP/1.0+ response
	r := bufio.NewReader(bytes.NewBuffer(body))
	first, _, err := r.ReadLine()
	if err != nil || !bytes.Contains(first, []byte("HTTP/")) {
		// response was HTTP/0.9
		resp.Proto = "HTTP/0.9"
		resp.ProtoMinor = 9
		resp.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		return resp, nil
	}

	split := strings.SplitN(string(first), " ", 3)
	if len(split) != 3 {
		return nil, fmt.Errorf("unable to parse response line: %q", string(first))
	}

	resp.Proto = split[0]
	resp.Status = split[1] + " " + split[2]
	resp.StatusCode, _ = strconv.Atoi(split[1])

	resp.Header = make(http.Header)

	// parse the headers
	for {
		line, _, err := r.ReadLine()
		if err != nil || len(line) == 0 {
			break
		}

		split := strings.SplitN(string(line), ":", 2)
		if len(split) != 2 {
			return nil, fmt.Errorf("unable to parse header: %q", string(line))
		}

		resp.Header.Add(split[0], split[1])
	}

	resp.Body = ioutil.NopCloser(r)
	return resp, nil
}
