// =====================================================================
//
// udpmask.go -
//
// Created by skywind on 2024/11/18
// Last Modified: 2024/11/18 19:45:07
//
// =====================================================================
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/skywind3000/uboost/forward"
	"github.com/skywind3000/uboost/service"
)

func main() {
	src := flag.String("src", "", "local address, eg: 0.0.0.0:8080")
	dst := flag.String("dst", "", "destination address, eg: 8.8.8.8:443")
	mask := flag.String("mask", "", "encryption/decryption key")
	mark := flag.Uint("mark", 0, "fwmark value")
	side := flag.String("side", "", "forward side: client/server")
	fec := flag.Int("fec", 0, "fec redundancy")
	flag.Parse()

	flag.Usage = func() {
		flagSet := flag.CommandLine
		fmt.Printf("Usage of %s:\n", os.Args[0])
		order := []string{"side", "src", "dst", "mask", "fec", "mark"}
		for _, name := range order {
			flag := flagSet.Lookup(name)
			tname := flag.Value
			ttext := fmt.Sprintf("%T", tname)
			ttext = ttext[6 : len(ttext)-5]
			fmt.Printf("  -%s %s\n", flag.Name, ttext)
			fmt.Printf("        %s\n", flag.Usage)
		}
	}

	if src == nil || dst == nil || side == nil {
		flag.Usage()
		return
	}
	if *src == "" || *dst == "" {
		flag.Usage()
		return
	}

	config := service.ServiceConfig{
		SrcAddr: *src,
		DstAddr: *dst,
		Mask:    *mask,
		Mark:    uint32(*mark),
		Fec:     *fec,
	}
	if *side == "server" {
		config.Side = forward.ForwardSideServer
	} else if *side == "client" {
		config.Side = forward.ForwardSideClient
	} else {
		flag.Usage()
		return
	}
	if config.Fec < 0 || config.Fec > 10 {
		flag.Usage()
		return
	}
	service.StartService(config)
}
