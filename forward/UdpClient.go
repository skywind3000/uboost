// =====================================================================
//
// UdpClient.go -
//
// Created by skywind on 2024/11/17
// Last Modified: 2024/11/17 05:47:05
//
// =====================================================================
package forward

import (
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

// udp session
type UdpClient struct {
	receiver func(client *UdpClient, data []byte) error
	closer   func(client *UdpClient)
	conn     *net.UDPConn
	closing  atomic.Bool
	dstAddr  *net.UDPAddr
	srcAddr  *net.UDPAddr
	logger   *log.Logger
	side     ForwardSide
	key      string
	mask     []byte
	cache    []byte
	timeout  int
	lock     sync.Mutex
	wg       sync.WaitGroup
}

func NewUdpClient() *UdpClient {
	self := &UdpClient{
		side:     ForwardSideServer,
		receiver: nil,
		closer:   nil,
		conn:     nil,
		dstAddr:  nil,
		srcAddr:  nil,
		mask:     nil,
		cache:    nil,
		key:      "",
		timeout:  300,
		logger:   nil,
		lock:     sync.Mutex{},
		wg:       sync.WaitGroup{},
		closing:  atomic.Bool{},
	}
	self.closing.Store(false)
	return self
}

func (self *UdpClient) SetCallback(receiver func(client *UdpClient, data []byte) error) {
	self.receiver = receiver
}

func (self *UdpClient) SetCloser(closer func(client *UdpClient)) {
	self.closer = closer
}

func (self *UdpClient) Open(srcAddr *net.UDPAddr, dstAddr *net.UDPAddr) error {
	self.lock.Lock()
	defer self.lock.Unlock()
	self.shutdown()
	err := error(nil)
	self.conn, err = net.DialUDP("udp", nil, dstAddr)
	if err != nil {
		self.conn = nil
		return err
	}
	self.dstAddr = AddressClone(dstAddr)
	self.srcAddr = AddressClone(srcAddr)
	self.closing.Store(false)
	self.cache = make([]byte, 65536)
	self.wg = sync.WaitGroup{}
	self.wg.Add(1)
	go self.recvLoop()
	return nil
}

func (self *UdpClient) shutdown() {
	self.closing.Store(true)
	if self.conn != nil {
		self.conn.Close()
		self.conn = nil
		self.wg.Wait()
	}
	self.cache = nil
}

func (self *UdpClient) Close() {
	self.lock.Lock()
	self.shutdown()
	self.lock.Unlock()
}

func (self *UdpClient) recvLoop() {
	buf := make([]byte, 65536)
	for !self.closing.Load() {
		duration := time.Second * time.Duration(self.timeout)
		self.conn.SetReadDeadline(time.Now().Add(duration))
		n, err := self.conn.Read(buf)
		if err != nil {
			break
		}
		data := buf[:n]
		if self.receiver != nil {
			if len(self.mask) > 0 {
				EncryptRC4(data, data, self.mask)
			}
			self.receiver(self, data)
		}
	}
	if self.closer != nil {
		self.closer(self)
	}
	self.wg.Done()
}

func (self *UdpClient) SendTo(data []byte) error {
	if self.conn == nil {
		return nil
	}
	var dst []byte = self.cache[:len(data)]
	EncryptRC4(dst, data, self.mask)
	now := time.Now()
	duration := time.Second * time.Duration(self.timeout)
	self.conn.SetWriteDeadline(now.Add(time.Millisecond * 100))
	self.conn.SetReadDeadline(now.Add(duration))
	_, err := self.conn.Write(dst)
	if err != nil {
		if self.logger != nil {
			self.logger.Printf("sendto error: %s\n", err)
		}
	}
	return err
}
