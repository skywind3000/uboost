// =====================================================================
//
// service.go -
//
// Created by skywind on 2024/11/18
// Last Modified: 2024/11/18 19:49:04
//
// =====================================================================
package service

import (
	"log"
	"os"
	"os/signal"

	"github.com/skywind3000/uboost/forward"
)

type ServiceConfig struct {
	Side    forward.ForwardSide
	SrcAddr string
	DstAddr string
	Mask    string
	Fec     int
	Mark    uint32
}

func StartService(config ServiceConfig) int {
	service := forward.NewUdpForward(config.Side)
	service.SetFec(config.Fec)
	saddr := forward.AddressResolve(config.SrcAddr)
	daddr := forward.AddressResolve(config.DstAddr)

	logger := log.Default()
	logger.Printf("config: %v\n", config)

	if saddr == nil {
		logger.Printf("config: invalid src address: \"%s\"\n", config.SrcAddr)
		return 1
	}

	if daddr == nil {
		logger.Printf("config: invalid dst address: \"%s\"\n", config.DstAddr)
		return 1
	}

	service.SetLogger(logger)
	service.SetMark(config.Mark)

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt, os.Kill)

	err := service.Open(saddr, daddr, config.Mask)
	if err != nil {
		return 2
	}

	<-sigch
	service.Close()

	return 0
}
