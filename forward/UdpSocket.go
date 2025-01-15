// =====================================================================
//
// # UdpSocket.go - UdpSocket implementation
//
// Last Modified: 2024/07/15 10:25:24
//
// =====================================================================
package forward

import (
	"log"
	"net"
	"sync"
	"sync/atomic"
)

const (
	UDP_RECV_COUNT int = 1
)

// basic udp socket implementation
type UdpSocket struct {
	receiver  func(data []byte, addr *net.UDPAddr) error
	metric    UdpMetric
	conn      *net.UDPConn
	closing   atomic.Bool
	wg        sync.WaitGroup
	count     int
	lock      sync.Mutex
	logPacket *log.Logger
	logError  *log.Logger
}

func NewUdpSocket() *UdpSocket {
	self := &UdpSocket{
		receiver:  nil,
		conn:      nil,
		count:     4,
		logPacket: nil,
		logError:  nil,
		lock:      sync.Mutex{},
	}
	self.closing.Store(false)
	self.metric.Clear()
	return self
}

func (self *UdpSocket) SetCallback(receiver func(data []byte, addr *net.UDPAddr) error) {
	self.receiver = receiver
}

func (self *UdpSocket) Open(addr *net.UDPAddr, count int) error {
	self.lock.Lock()
	defer self.lock.Unlock()
	self.shutdown()
	self.metric.Clear()
	err := error(nil)
	self.conn, err = net.ListenUDP("udp", addr)
	if err != nil {
		self.conn = nil
		return err
	}
	self.count = count
	if self.count < 1 {
		self.count = 1
	}
	self.closing.Store(false)
	self.wg.Add(self.count)
	for i := 0; i < self.count; i++ {
		go self.backgroundReceiving()
	}
	return nil
}

func (self *UdpSocket) shutdown() error {
	self.closing.Store(true)
	if self.conn != nil {
		self.conn.Close()
	}
	self.wg.Wait()
	if self.conn != nil {
		self.conn.Close()
		self.conn = nil
	}
	return nil
}

func (self *UdpSocket) Close() {
	self.lock.Lock()
	defer self.lock.Unlock()
	self.shutdown()
}

func (self *UdpSocket) backgroundReceiving() {
	defer self.wg.Done()
	if self.logPacket != nil {
		self.logPacket.Printf("backgroundReceiving start")
	}
	data := make([]byte, 1024*96)
	for !self.closing.Load() {
		n, addr, err := self.conn.ReadFromUDP(data)
		if err != nil {
			if self.logError != nil {
				self.logError.Printf("RecvFrom error: %s", err)
			}
			break
		} else {
			if self.logPacket != nil {
				self.logPacket.Printf("ReadFromUDP: size=%d addr=%s", n, addr.String())
			}
		}
		if n >= 0 {
			if self.receiver != nil {
				err := self.receiver(data[:n], addr)
				if err != nil {
					self.metric.IncPacketDropped()
				} else {
					self.metric.IncPacketReceived()
				}
			} else {
				self.metric.IncPacketDropped()
			}
		}
	}
	if self.logPacket != nil {
		self.logPacket.Printf("backgroundReceiving exit")
	}
}

func (self *UdpSocket) SendTo(data []byte, addr *net.UDPAddr) int {
	if self.conn != nil {
		n, err := self.conn.WriteToUDP(data, addr)
		if err != nil {
			if self.logError != nil {
				self.logError.Printf("SendTo: %s", err)
			}
			return -1
		} else {
			if self.logPacket != nil {
				self.logPacket.Printf("SendTo: size=%d addr=%s", len(data), addr.String())
			}
		}
		if n > 0 {
			self.metric.IncPacketSent()
		}
	}
	return 0
}

func (self *UdpSocket) SendBatch(data [][]byte, addr []*net.UDPAddr) int {
	for i := 0; i < len(data); i++ {
		self.SendTo(data[i], addr[i])
	}
	return 0
}

func (self *UdpSocket) GetMetric() *UdpMetric {
	return &self.metric
}

func (self *UdpSocket) SetLogPacket(logger *log.Logger) {
	self.logPacket = logger
}

func (self *UdpSocket) SetLogError(logger *log.Logger) {
	self.logError = logger
}

func (self *UdpSocket) SetOption(option int, value interface{}) error {
	var err error = nil
	switch option {
	case UDP_RECV_COUNT:
		self.count = value.(int)
	}
	return err
}

func (self *UdpSocket) IsClosing() bool {
	return self.closing.Load()
}
