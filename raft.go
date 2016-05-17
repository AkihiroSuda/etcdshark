package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/coreos/etcd/raft/raftpb"
	"github.com/google/gopacket"
)

func handleHTTP11RaftResp(net, transport gopacket.Flow,
	resp *http.Response, req *http.Request, w io.Writer) (bool, error) {
	if resp.StatusCode/100 != 2 {
		return false, nil
	}
	if strings.HasPrefix(req.URL.Path, "/raft/stream/message") {
		return handleHTTP11RaftMessageResp(net, transport,
			resp, req, w)
	}
	if strings.HasPrefix(req.URL.Path, "/raft/stream/msgapp") {
		return handleHTTP11RaftMsgAppV2Resp(net, transport,
			resp, req, w)
	}
	return false, nil
}

func handleHTTP11RaftMessageResp(net, transport gopacket.Flow,
	resp *http.Response, req *http.Request, w io.Writer) (bool, error) {
	var err error
	var m raftpb.Message
	var l uint64
	for {
		if err = binary.Read(resp.Body, binary.BigEndian, &l); err != nil && err != io.EOF {
			return true, err
		}
		if err == io.EOF {
			return true, nil
		}
		buf := make([]byte, int(l))
		if _, err := io.ReadFull(resp.Body, buf); err != nil {
			return true, err
		}
		err = m.Unmarshal(buf)
		if err != nil {
			return true, err
		}
		s := fmt.Sprintf("%#v\n", m)
		_, err = w.Write([]byte(s))
		if err != nil {
			return true, err
		}
	}

}

func handleHTTP11RaftMsgAppV2Resp(net, transport gopacket.Flow,
	resp *http.Response, req *http.Request, w io.Writer) (bool, error) {
	br := bufio.NewReader(resp.Body)
	// not yet supported, just hexdump it
	log.Printf("hexdumping unknown Raft stream %s to stderr",
		prettifyFlow(net, transport))
	return true, hexdumpAll(br, os.Stderr)
}
