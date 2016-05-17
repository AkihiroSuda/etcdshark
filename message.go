package main

import (
	"bufio"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/google/gopacket"
)

type messageType uint

const (
	unsupportedMessageType messageType = iota
	http11Request
	http11Response
	http2Request  // not supported yet
	http2Response // not supported yet
)

func peekMessageType(br *bufio.Reader) (messageType, error) {
	head4, err := br.Peek(4)
	if err != nil {
		return unsupportedMessageType, err
	}
	if string(head4) == "HTTP" {
		return http11Response, nil
	}
	// FIXME: check HTTP/1.1 methods?
	return http11Request, nil
}

// handle parses messages in br and write human-readable information to w.
func handleMessage(net, transport gopacket.Flow, br *bufio.Reader, w io.Writer) error {
	for {
		messageType, err := peekMessageType(br)
		if err != nil {
			return err
		}
		switch messageType {
		case http11Request:
			req, err := http.ReadRequest(br)
			if err != nil {
				return err
			}
			if err = handleHTTP11Req(net, transport, req, w); err != nil {
				return err
			}
		case http11Response:
			resp, err := http.ReadResponse(br, nil)
			if err != nil {
				return err
			}
			if err = handleHTTP11Resp(net, transport, resp, w); err != nil {
				return err
			}
		default:
			log.Printf("hexdumping unknown stream %s to stderr",
				prettifyFlow(net, transport))
			// we use stderr rather than w for unknown stream
			return hexdumpAll(br, os.Stderr)
		}
	}
}
