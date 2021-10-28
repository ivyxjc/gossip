package gossip

import (
	"net"
	"time"
)

const (
	StateAlive = iota
	StateSuspect
	StateDead
)

type Node struct {
	Name string
	Addr net.IP
}

type NodeState struct {
	Node
	Incarnation int       // Last known incarnation number
	State       int       // Current State
	StateChange time.Time // Time last state change happened
}

type ackHandler struct {
	handler func()
	timer   *time.Timer
}

func (m *Memberlist) schedule() {
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

func (m *Memberlist) deSchedule() {
	m.tickerLock.Lock()
	m.ticker.Stop()
	m.ticker = nil
	m.tickerLock.Unlock()
	m.stopTick <- struct{}{}
}

func (m *Memberlist) tick() {

}

func (m *Memberlist) setAckHandler(seqNo uint32, handler func(), timeout time.Duration) {
	ah := &ackHandler{handler: handler, timer: nil}
	m.ackLock.Lock()
	m.ackHandlers[seqNo] = ah
	m.ackLock.Unlock()

	ah.timer = time.AfterFunc(timeout, func() {
		m.ackLock.Lock()
		delete(m.ackHandlers, seqNo)
		m.ackLock.Unlock()
	})
}

func (m *Memberlist) invokeAckHandler(seqNo uint32) {
	m.ackLock.Lock()
	ah, ok := m.ackHandlers[seqNo]
	delete(m.ackHandlers, seqNo)
	m.ackLock.Unlock()
	if !ok {
		return
	}
	ah.timer.Stop()
	ah.handler()
}

func (m *Memberlist) aliveNode(a *alive) {
}

func (m *Memberlist) suspectNode(s *suspect) {
}

func (m *Memberlist) deadNode(d *dead) {
}

func (m *Memberlist) mergeState(remote []pushNodeState) {

}
