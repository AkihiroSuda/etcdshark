package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/google/gopacket"
)

// FIXME: eliminate global variables
var http11ReqM = make(map[string]*http.Request)
var http11ReqMMu sync.Mutex

func prettifyFlow(net, transport gopacket.Flow) string {
	return fmt.Sprintf("%s:%s->%s:%s",
		net.Src(), transport.Src(),
		net.Dst(), transport.Dst())
}

func retainHTTP11Req(net, transport gopacket.Flow, req *http.Request) {
	http11ReqMMu.Lock()
	defer http11ReqMMu.Unlock()
	k := prettifyFlow(net, transport)
	http11ReqM[k] = req
}

func takeHTTP11Req(net, transport gopacket.Flow, resp *http.Response) *http.Request {
	http11ReqMMu.Lock()
	defer http11ReqMMu.Unlock()
	k := prettifyFlow(net.Reverse(), transport.Reverse())
	req, ok := http11ReqM[k]
	if ok {
		delete(http11ReqM, k)
		return req
	}
	return nil
}

func prettifyHTTP11Req(net, transport gopacket.Flow, req *http.Request) string {
	s := fmt.Sprintf("%s %s (%s) [header:%s]",
		req.Method, req.URL,
		prettifyFlow(net, transport),
		req.Header)
	return s
}

func prettifyHTTP11Resp(net, transport gopacket.Flow, resp *http.Response) string {
	s := fmt.Sprintf("%s (%s) [header:%v]",
		resp.Status,
		prettifyFlow(net, transport),
		resp.Header)
	return s
}

func handleHTTP11Req(net, transport gopacket.Flow, req *http.Request, w io.Writer) error {
	s := fmt.Sprintf("==>%s\n", prettifyHTTP11Req(net, transport, req))
	w.Write([]byte(s))
	defer req.Body.Close()
	retainHTTP11Req(net, transport, req)
	return nil
}

func handleHTTP11Resp(net, transport gopacket.Flow, resp *http.Response, w io.Writer) error {
	req := takeHTTP11Req(net, transport, resp)
	reqS := ""
	if req != nil {
		reqS = prettifyHTTP11Req(net.Reverse(), transport.Reverse(), req)
	}
	s := fmt.Sprintf("<==%s [req:%s]\n",
		prettifyHTTP11Resp(net, transport, resp), reqS)
	w.Write([]byte(s))
	defer resp.Body.Close()

	handled, err := handleHTTP11JSONResp(net, transport, resp, req, w)
	if handled {
		return err
	}
	handled, err = handleHTTP11RaftResp(net, transport, resp, req, w)
	if handled {
		return err
	}
	log.Printf("hexdumping unknown HTTP/1.1 stream %s to stderr",
		prettifyFlow(net, transport))
	br := bufio.NewReader(resp.Body)
	return hexdumpAll(br, os.Stderr)
}
