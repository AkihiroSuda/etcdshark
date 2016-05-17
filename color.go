package main

import (
	"io"
	"math/rand"
	"time"

	"github.com/google/gopacket"
	"github.com/hpcloud/golor"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type coloredWriter struct {
	w  io.Writer
	fg int
	bg int
}

func newColoredWriter(net, transport gopacket.Flow, w io.Writer) *coloredWriter {

	// set bg and fg color
	// https://en.wikipedia.org/wiki/File:Xterm_256color_chart.svg
	//
	// TODO: should we use hash(net,transport) rather than rand?
	bg := rand.Intn(16)
	fg := 15 // white
	if bg == 7 || bg == 10 || bg == 11 ||
		bg == 14 || bg == 15 {
		fg = 0 // black
	}

	cw := &coloredWriter{
		w:  w,
		fg: fg,
		bg: bg,
	}
	return cw
}

func (cw *coloredWriter) Write(p []byte) (n int, err error) {
	if cw.fg >= 0 || cw.bg >= 0 {
		colored := golor.Colorize(string(p), cw.fg, cw.bg)
		written, err := cw.w.Write([]byte(colored))
		// FIXME: if written < len([]byte(colored))
		return written, err
	}
	return cw.w.Write(p)
}
