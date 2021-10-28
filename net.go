package p2p

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
)

const (
	pingMsg = iota
	indirectPingMsg
	ackRespMsg
	suspectMsg
	aliveMsg
	deadMsg
	pushPullMsg
)

const (
	udpBufSize = 65536
	udpSendBuf = 1500
)

// ping request sent directly to node
type ping struct {
	SeqNo int
}

type indirectPingReq struct {
	SeqNo  int
	Target string
}

// ack response is sent for a ping
type ackResp struct {
	SeqNo int
}

// suspect is broadcast when we suspect a node is dead
type suspect struct {
	Incarnation int
	Node        string
}

// alive is broadcast when we  know a node is alive
type alive struct {
	Incarnation int
	Node        string
}

// dead is broadcast when we confirm a node is dead
type dead struct {
	Incarnation int
	Node        string
}

func (m *Memberlist) tcpListen() {
	for {
		conn, err := m.tcpListener.AcceptTCP()
		if err != nil {
			if neterr, ok := err.(net.Error); ok && !neterr.Temporary() {
				break
			}
			log.Printf("[ERR] Error acception TCP connection: %s", err)
			continue
		}
		go m.handleConnection(conn)
	}

}

func (m *Memberlist) udpListen() {
	mainBuf := make([]byte, udpBufSize)
	var n int
	var msgType uint32
	var addr net.Addr
	var err error

	for {
		buf := mainBuf[0:udpBufSize]
		n, addr, err = m.udpListener.ReadFrom(buf)
		if err != nil {
			if neterr, ok := err.(net.Error); ok && !neterr.Temporary() {
				break
			}
			log.Printf("[ERR] Error acception UDP packet: %s", err)
			continue
		}
		buf = buf[0:n]
		if len(buf) < 4 {
			log.Printf("[ERR] UDP packet too short (%d bytes). From %s", len(buf), addr)
		}

		msgType = binary.BigEndian.Uint32(buf[0:4])
		buf = buf[4:]
		switch msgType {
		case pingMsg:
			m.handlePing(buf, addr)
		case indirectPingMsg:
			m.handleIndirectPing(buf, addr)
		case ackRespMsg:
			m.handleAck(buf, addr)
		case suspectMsg:
			m.handleSuspect(buf, addr)
		case aliveMsg:
			m.handleAlive(buf, addr)
		case deadMsg:
			m.handleDead(buf, addr)
		default:
			log.Printf("[ERR] UDP msg type (%d) not supported. From: %s", msgType, addr)
			continue
		}
	}
}

func (m *Memberlist) sendMsg(to net.Addr, msg *bytes.Buffer) error {
	_, err := m.udpListener.WriteTo(msg.Bytes(), to)
	return err
}

func (m *Memberlist) handleConnection(conn *net.TCPConn) {
}

func (m *Memberlist) handlePing(buf []byte, from net.Addr) {
	var p ping
	if err := decode(buf, &p); err != nil {
		log.Printf("[ERR]fail to parse ping request: %+v", err)
	}
	ack := ackResp{SeqNo: p.SeqNo}
	out, err := encode(ackRespMsg, ack)
	if err != nil {
		log.Printf("[ERR] Fail to encode ack response: %+v", err)
	}
	if err := m.sendMsg(from, out); err != nil {
		log.Printf("[ERR] Failed to send ack: %s", err)
	}
}

func (m *Memberlist) handleIndirectPing(buf []byte, from net.Addr) {

}

func (m *Memberlist) handleAck(buf []byte, from net.Addr) {

}

func (m *Memberlist) handleSuspect(buf []byte, from net.Addr) {

}

func (m *Memberlist) handleAlive(buf []byte, from net.Addr) {

}

func (m *Memberlist) handleDead(buf []byte, from net.Addr) {

}
