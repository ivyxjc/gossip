package gossip

import "testing"

func GetMemberlist(t *testing.T) *Memberlist {
	c := DefaultConfig()

	var m *Memberlist
	var err error
	for i := 0; i < 100; i++ {
		m, err = Create(c)
		if err == nil {
			return m
		}
		c.TCPPort++
		c.UDPPort++
	}
	t.Fatalf("failed to start: %v", err)
	return nil
}
