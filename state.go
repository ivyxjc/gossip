package p2p

import (
	"net"
	"time"
)

type Node struct {
	Name string
	Addr net.Addr
}

type NodeState struct {
	Node
	Incarnation int       // Last known incarnation number
	State       int       // Current State
	StateChange time.Time // Time last state change happened
}

func (m *MemberList) schedule() {
	m.tickerLock.Lock()
	m.ticker = time.NewTicker(m.config.Interval)
	C := m.ticker.C
	m.tickerLock.Unlock()
	go func() {
		for {
			select {
			case <-C:
				m.tick()
			case <-m.stopTick:
				return
			}
		}
	}()
}

func (m *MemberList) deSchedule() {
	m.tickerLock.Lock()
	m.ticker.Stop()
	m.ticker = nil
	m.tickerLock.Unlock()
	m.stopTick <- struct{}{}
}

func (m *MemberList) tick() {

}
