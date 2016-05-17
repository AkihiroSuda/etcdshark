// original: https://github.com/google/gopacket/blob/master/examples/httpassembly/main.go

package main

import (
	"flag"
	"log"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/tcpassembly"
)

func main() {
	// parse flags
	iface := flag.String("i", "eth0", "Interface to get packets from")
	fname := flag.String("r", "", "Filename to read from, overrides -i")
	snaplen := flag.Int("s", 1600, "SnapLen for pcap packet capture")
	filter := flag.String("f", "tcp", "BPF filter for pcap")
	colored := flag.Bool("c", false, "colored output")
	flag.Parse()

	var ph *pcap.Handle
	var err error

	// set up pcap
	if *fname != "" {
		log.Printf("Reading from pcap dump %q", *fname)
		ph, err = pcap.OpenOffline(*fname)
	} else {
		log.Printf("Starting capture on interface %q", *iface)
		ph, err = pcap.OpenLive(*iface, int32(*snaplen), true, pcap.BlockForever)
	}
	if err != nil {
		log.Fatal(err)
	}

	if err := ph.SetBPFFilter(*filter); err != nil {
		log.Fatal(err)
	}

	packetSource := gopacket.NewPacketSource(ph, ph.LinkType())

	// set up assembler for our stream
	streamFactory, err := newStreamFactory(*colored)
	if err != nil {
		log.Fatal(err)
	}
	streamPool := tcpassembly.NewStreamPool(streamFactory)
	assembler := tcpassembly.NewAssembler(streamPool)
	if err != nil {
		log.Fatal(err)
	}

	// main loop
	for {
		select {
		case packet := <-packetSource.Packets():
			// A nil packet indicates the end of a pcap file.
			if packet == nil {
				// FIXME: flush the printer here
				return
			}
			tcp := tcpLayer(packet)
			if tcp == nil {
				continue
			}
			assembler.AssembleWithTimestamp(packet.NetworkLayer().NetworkFlow(), tcp, packet.Metadata().Timestamp)
		case <-time.Tick(time.Minute):
			// Every minute, flush connections that haven't seen activity in the past 2 minutes.
			assembler.FlushOlderThan(time.Now().Add(time.Minute * -2))
		}
	}
}

func tcpLayer(packet gopacket.Packet) *layers.TCP {
	if packet.NetworkLayer() == nil || packet.TransportLayer() == nil || packet.TransportLayer().LayerType() != layers.LayerTypeTCP {
		return nil
	}
	tcp := packet.TransportLayer().(*layers.TCP)
	return tcp
}
