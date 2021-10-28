package p2p

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"testing"
)

func TestChannelIndex(t *testing.T) {
	ch1 := make(chan net.Addr)
	ch2 := make(chan net.Addr)
	ch3 := make(chan net.Addr)
	list := []chan<- net.Addr{ch1, ch2, ch3}

	if channelIndex(list, ch1) != 0 {
		t.Fatalf("bad index")
	}
	if channelIndex(list, ch2) != 1 {
		t.Fatalf("bad index")
	}
	if channelIndex(list, ch3) != 2 {
		t.Fatalf("bad index")
	}

	ch4 := make(chan net.Addr)
	if channelIndex(list, ch4) != -1 {
		t.Fatalf("bad index")
	}
}

func TestChannelIndex_Empty(t *testing.T) {
	ch := make(chan net.Addr)
	if channelIndex(nil, ch) != -1 {
		t.Fatalf("bad index")
	}
}

func TestChannelDelete(t *testing.T) {
	ch1 := make(chan net.Addr)
	ch2 := make(chan net.Addr)
	ch3 := make(chan net.Addr)
	list := []chan<- net.Addr{ch1, ch2, ch3}

	// Delete ch2
	list = channelDelete(list, 1)

	if len(list) != 2 {
		t.Fatalf("bad len")
	}
	if channelIndex(list, ch1) != 0 {
		t.Fatalf("bad index")
	}
	if channelIndex(list, ch3) != 1 {
		t.Fatalf("bad index")
	}
}

func TestEncodeDecode(t *testing.T) {
	msg := &ping{SeqNo: 100}
	buf, err := encode(pingMsg, msg)
	if err != nil {
		t.Fatalf("unexpected err: %s", err)
	}
	var out ping
	if err := decode(buf.Bytes()[4:], &out); err != nil {
		t.Fatalf("unexpected err: %s", err)
	}
	if msg.SeqNo != out.SeqNo {
		t.Fatalf("bad sequence no")
	}
}

type s struct {
	data map[string]interface{}
}

func TestGob(t *testing.T) {
	var s1 = s{
		data: make(map[string]interface{}, 8),
	}
	s1.data["count"] = 1
	s1.data["count2"] = 2

	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(s1.data)
	if err != nil {
		fmt.Printf("gob encode failed, err: %+v", err)
	}
	b := buf.Bytes()
	fmt.Println(b)

	var s2 = s{
		data: make(map[string]interface{}, 8),
	}
	dec := gob.NewDecoder(bytes.NewBuffer(b))
	err = dec.Decode(&s2.data)
	if err != nil {
		fmt.Println("gob decode failed, err", err)
		return
	}
	fmt.Println(s2.data)

}
