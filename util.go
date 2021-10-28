package p2p

import (
	"bytes"
	"encoding/gob"
	"net"
)

// channelIndex returns the index of a channel in a list or
// -1 if not found
func channelIndex(list []chan<- net.Addr, ch chan<- net.Addr) int {
	for i := 0; i < len(list); i++ {
		if list[i] == ch {
			return i
		}
	}
	return -1
}

// channelDelete takes an index and removes the element at that index
func channelDelete(list []chan<- net.Addr, idx int) []chan<- net.Addr {
	n := len(list)
	list[idx], list[n-1] = list[n-1], nil
	return list[0 : n-1]
}

func decode(buf []byte, out interface{}) error {
	r := bytes.NewBuffer(buf)
	dec := gob.NewDecoder(r)
	return dec.Decode(out)
}

func encode(in interface{}) (*bytes.Buffer, error) {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(in)
	return buf, err
}
