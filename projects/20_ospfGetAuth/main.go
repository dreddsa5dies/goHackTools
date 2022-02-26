package main

import (
	"fmt"
	"os"

	"github.com/google/gopacket/layers"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

var (
	handle *pcap.Handle
)

func main() {
	if len(os.Args) == 1 {
		fmt.Printf("Use: %v PCAPFile\n", os.Args[0])
	} else {
		// Open file instead of device
		handle, err := pcap.OpenOffline(os.Args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "PCAPFile: %v\n", err)
			os.Exit(1)
		}
		defer handle.Close()

		// Loop through packets in file
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		for packet := range packetSource.Packets() {
			// проходим по всем пакетам и ранжируем по уровню ospf
			ospf := packet.Layer(layers.LayerTypeOSPF)
			if nil != ospf {
				ospf, _ := ospf.(*layers.OSPFv2)
				switch {
				case ospf.AuType == 1:
					fmt.Printf("Simple version: %v\n", ospf.AuType)
					fmt.Printf("OSPF Pass: %v\n", ospf.Authentication)
				case ospf.AuType == 2:
					fmt.Printf("MD5 version: %v\n", ospf.AuType)
					fmt.Printf("Authentication: %d\n", ospf.Authentication)
				}
			}
		}
	}
}
