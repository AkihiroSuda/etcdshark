package main

import (
	"bufio"
	"io"
	"log"
	"os"

	"github.com/google/gopacket"
	"github.com/google/gopacket/tcpassembly"
	"github.com/google/gopacket/tcpassembly/tcpreader"
)

type streamFactory struct {
	colored bool
}

func (sf *streamFactory) New(net, transport gopacket.Flow) tcpassembly.Stream {
	s := &stream{
		net:       net,
		transport: transport,
		rs:        tcpreader.NewReaderStream(),
		colored:   sf.colored,
	}
	go s.run()
	return &s.rs
}

type stream struct {
	net, transport gopacket.Flow
	rs             tcpreader.ReaderStream
	colored        bool
}

func (s *stream) run() {
	br := bufio.NewReader(&s.rs)
	var w io.Writer = os.Stdout
	if s.colored {
		w = newColoredWriter(s.net, s.transport, os.Stdout)
	}

	err := handleMessage(s.net, s.transport, br, w)
	if err != nil && err != io.EOF {
		log.Printf("Error: %s", err)
	}
}

func newStreamFactory(colored bool) (tcpassembly.StreamFactory, error) {
	return &streamFactory{colored}, nil
}
