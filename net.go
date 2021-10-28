package p2p

import (
	"log"
	"net"
)

func (m *MemberList) tcpListen() {
	for {
		conn, err := m.tcpListener.AcceptTCP()
		if err != nil {
			if neterr, ok := err.(net.Error); ok && !neterr.Temporary() {
				break
			}
			log.Printf("[ERR] Error acception TCP connection: %s", err)
		}
		go m.handleConnection(conn)
	}

}

func (m *MemberList) udpListen() {
}

func (m *MemberList) handleConnection(conn *net.TCPConn) {
}
