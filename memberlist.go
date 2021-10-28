package p2p

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

type Config struct {
	Name           string        // Node Name
	BindAddr       string        // Binding address
	UDPPort        int           // UDP port
	TCPPort        int           // TCP port
	Fanout         int           // Number of nodes to publish to per round
	IndirectChecks int           // Number of indirect checks to use
	Retransmits    int           // Retransmit multiplier for message
	PushPullFreq   float32       // How often we do a Push/Pull update
	RTT            time.Duration // 99% percentile of round-trip-time
	Interval       time.Duration // Interval length
}

type Memberlist struct {
	config *Config

	tcpListener *net.TCPListener
	udpListener *net.UDPConn

	notifyLock sync.RWMutex

	notifyJoin  []chan<- net.Addr
	notifyLeave []chan<- net.Addr
	notifyFail  []chan<- net.Addr

	sequenceNum uint32
	incarnation uint32

	nodeLock sync.RWMutex
	nodes    []*NodeState
	nodeMap  map[string]*NodeState

	tickerLock sync.Mutex
	ticker     *time.Ticker
	stopTick   chan struct{}

	ackLock     sync.Mutex
	ackHandlers map[uint32]*ackHandler
}

func DefaultConfig() *Config {
	hostname, _ := os.Hostname()
	return &Config{
		Name:           hostname,
		BindAddr:       "0.0.0.0",
		UDPPort:        7946,
		TCPPort:        7946,
		Fanout:         3,
		IndirectChecks: 3,
		Retransmits:    4,
		PushPullFreq:   0.05,
		RTT:            20 * time.Millisecond,
		Interval:       1 * time.Second,
	}

}

func newMemberList(conf *Config) (*Memberlist, error) {
	tcpAddr := fmt.Sprintf("%s:%d", conf.BindAddr, conf.TCPPort)
	tcpLn, err := net.Listen("tcp", tcpAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to start TCP listener. Err: %s", err)
	}

	udpAddr := fmt.Sprintf("%s:%d", conf.BindAddr, conf.UDPPort)
	udpLn, err := net.ListenPacket("udp", udpAddr)

	if err != nil {
		tcpLn.Close()
		return nil, fmt.Errorf("fail to start UDP listener. Err: %s", err)
	}

	m := &Memberlist{
		config:      conf,
		tcpListener: tcpLn.(*net.TCPListener),
		udpListener: udpLn.(*net.UDPConn),
		nodeMap:     make(map[string]*NodeState),
		stopTick:    make(chan struct{}),
		ackHandlers: make(map[uint32]*ackHandler),
	}
	go m.tcpListen()
	go m.udpListen()
	return m, nil
}

func Create(conf *Config) (*Memberlist, error) {
	m, err := newMemberList(conf)
	if err != nil {
		return nil, err
	}
	m.schedule()
	return m, nil
}

func Join(conf *Config) (*Memberlist, error) {
	m, err := newMemberList(conf)
	if err != nil {
		return nil, err
	}
	m.schedule()
	return m, nil
}

func (m *Memberlist) NotifyJoin(ch chan<- net.Addr) {
	m.notifyLock.Lock()
	defer m.notifyLock.Unlock()
	if channelIndex(m.notifyJoin, ch) >= 0 {
		return
	}
	m.notifyJoin = append(m.notifyJoin, ch)
}

func (m *Memberlist) NotifyLeave(ch chan<- net.Addr) {
	m.notifyLock.Lock()
	defer m.notifyLock.Unlock()
	if channelIndex(m.notifyLeave, ch) >= 0 {
		return
	}
	m.notifyLeave = append(m.notifyLeave, ch)
}

func (m *Memberlist) NotifyFail(ch chan<- net.Addr) {
	m.notifyLock.Lock()
	defer m.notifyLock.Unlock()

	if channelIndex(m.notifyFail, ch) >= 0 {
		return
	}
	m.notifyFail = append(m.notifyFail, ch)
}
