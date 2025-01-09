// =====================================================================
//
// UdpForward.go -
//
// Created by skywind on 2024/11/18
// Last Modified: 2024/11/18 17:19:16
//
// =====================================================================
package forward

import (
	"log"
	"net"
	"sync"
	"sync/atomic"
)

// udp forward
type UdpForward struct {
	side      ForwardSide
	clientMap sync.Map
	lock      sync.Mutex
	count     int
	mask      []byte
	udpServer *UdpSocket
	closing   atomic.Bool
	srcAddr   *net.UDPAddr
	dstAddr   *net.UDPAddr
	logger    *log.Logger
	wg        sync.WaitGroup
	mark      uint32
}

func NewUdpForward(side ForwardSide) *UdpForward {
	self := &UdpForward{
		side:      side,
		count:     0,
		udpServer: nil,
		srcAddr:   nil,
		dstAddr:   nil,
		mask:      nil,
		logger:    nil,
		mark:      0,
	}
	self.udpServer = NewUdpSocket()
	self.udpServer.SetCallback(self.onServerPacket)
	self.closing.Store(false)
	return self
}

func (self *UdpForward) SetLogger(logger *log.Logger) {
	self.logger = logger
	self.udpServer.SetLogError(logger)
	self.udpServer.SetLogPacket(logger)
	self.udpServer.SetLogPacket(nil)
}

func (self *UdpForward) SetMark(mark uint32) {
	self.mark = mark
}

func (self *UdpForward) Open(srcAddr *net.UDPAddr, dstAddr *net.UDPAddr, mask string) error {
	self.lock.Lock()
	defer self.lock.Unlock()
	self.shutdown()
	self.mask = []byte(mask)
	self.srcAddr = AddressClone(srcAddr)
	self.dstAddr = AddressClone(dstAddr)
	self.closing.Store(false)
	err := self.udpServer.Open(self.srcAddr, 0)
	if err == nil {
		if self.logger != nil {
			self.logger.Printf("Start forwarding %s -> %s\n", srcAddr.String(), dstAddr.String())
		}
	} else {
		if self.logger != nil {
			self.logger.Printf("Failed to open %s\n", srcAddr.String())
		}
		return err
	}
	return err
}

func (self *UdpForward) shutdown() {
	self.closing.Store(true)
	self.udpServer.Close()
	self.clientMap.Range(func(key, value interface{}) bool {
		client := value.(*UdpClient)
		client.Close()
		return true
	})
	self.clientMap = sync.Map{}
	self.wg.Wait()
}

func (self *UdpForward) Close() {
	self.lock.Lock()
	self.shutdown()
	self.lock.Unlock()
	if self.logger != nil {
		self.logger.Printf("Forwarding closed\n")
	}
}

func (self *UdpForward) newClient(key string, srcAddr *net.UDPAddr) *UdpClient {
	client := NewUdpClient()
	client.SetCallback(self.onClientPacket)
	client.SetCloser(self.onClientClose)
	client.side = self.side
	client.key = key
	client.mask = self.mask
	client.logger = self.logger
	self.wg.Add(1)
	err := client.Open(srcAddr, self.dstAddr)
	if err != nil {
		if self.logger != nil {
			self.logger.Printf("connection open failed: %s: %s\n", key, err)
		}
		client.Close()
		client = nil
		self.wg.Done()
	}
	return client
}

func (self *UdpForward) getClient(srcAddr *net.UDPAddr) *UdpClient {
	key := srcAddr.String()
	var client *UdpClient = nil
	cc, ok := self.clientMap.Load(key)
	if ok {
		client = cc.(*UdpClient)
	} else {
		self.lock.Lock()
		cc, ok = self.clientMap.Load(key)
		if ok {
			client = cc.(*UdpClient)
		} else {
			client = self.newClient(key, srcAddr)
			if client != nil {
				self.clientMap.Store(key, client)
			}
		}
		self.lock.Unlock()
		if self.logger != nil {
			if client != nil {
				self.logger.Printf("connection new: %s\n", key)
			} else {
			}
		}
	}
	return client
}

func (self *UdpForward) onClientClose(client *UdpClient) {
	key := client.key
	// don't block the Close() function
	go func() {
		if self.closing.Load() == false {
			self.lock.Lock()
			_, ok := self.clientMap.Load(key)
			if ok {
				self.clientMap.Delete(key)
			}
			client.Close()
			client.closer = nil
			client.receiver = nil
			self.lock.Unlock()
		}
		if self.logger != nil {
			self.logger.Printf("connection closed: %s\n", key)
		}
		self.wg.Done()
	}()
}

func (self *UdpForward) onClientPacket(client *UdpClient, data []byte) error {
	if self.closing.Load() == false {
		self.udpServer.SendTo(data, client.srcAddr)
	}
	return nil
}

func (self *UdpForward) onServerPacket(data []byte, addr *net.UDPAddr) error {
	if self.logger != nil {
		// self.logger.Printf("server packet: size=%d addr=%s\n", len(data), addr.String())
	}
	if self.closing.Load() == false {
		client := self.getClient(addr)
		client.SendTo(data)
	} else {
		if self.logger != nil {
			self.logger.Printf("server is closing\n")
		}
	}
	return nil
}
