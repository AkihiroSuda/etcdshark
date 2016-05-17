package main

import (
	"bufio"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/google/gopacket"
)

func handleHTTP11JSONResp(net, transport gopacket.Flow,
	resp *http.Response, req *http.Request, w io.Writer) (bool, error) {
	// NOTE: Response for /raft/probing contains
	// {"Content-Type": "text/plain"}, but its actual body is JSON
	if resp.Header.Get("Content-Type") != "application/json" &&
		(req != nil && req.URL.Path != "/raft/probing") {
		return false, nil
	}
	br := bufio.NewReader(resp.Body)
	s, err := ioutil.ReadAll(br)
	if err != nil {
		return true, err
	}
	_, err = w.Write(append(s, []byte("\n")...))
	return true, err
}
