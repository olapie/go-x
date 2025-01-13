package xnet

import (
	"net"
	"testing"

	"go.olapie.com/x/xrand"
)

//
//func TestMulticast(t *testing.T) {
//	ifi := GetPrivateIPv4Interface()
//	if ifi == nil {
//		t.Log("No PrivateIPv4Interface")
//		t.FailNow()
//	}
//	multiIP := GetMulticastIPv4String(ifi)
//	if multiIP == "" {
//		t.Fatal("no multi ip")
//	}
//	udpAddr, err := net.ResolveUDPAddr("udp", multiIP+":9999")
//	if err != nil {
//		t.Fatal(err, multiIP)
//	}
//	t.Log(udpAddr.String())
//	conn, err := net.ListenMulticastUDP("udp", ifi, udpAddr)
//	if err != nil {
//		t.Fatal(err)
//	}
//	defer conn.Close()
//	packet := RandomBytes(10)
//	buf := make([]byte, 2000)
//	received := make(chan error)
//
//	go func() {
//		nRead, addr, err := conn.ReadFrom(buf)
//		if err == nil {
//			t.Log(nRead, addr)
//			buf = buf[:nRead]
//		}
//		received <- err
//	}()
//
//	udpConn, err := net.DialUDP("udp", nil, udpAddr)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	_, err = udpConn.Write(packet)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	time.Sleep(time.Second)
//	select {
//	case err := <-received:
//		MustNotErrorT(t, err)
//	case <-time.After(time.Second):
//		t.Fatal("failed to receive udp packet")
//	}
//
//	MustEqualT(t, packet, buf)
//	t.Log(packet)
//	t.Log(buf)
//}

func TestGetBroadcastIPv4(t *testing.T) {
	ifi := GetPrivateIPv4Interface()
	if ifi == nil {
		t.Log("No PrivateIPv4Interface")
		t.FailNow()
	}
	ip := GetBroadcastIPv4(ifi)
	if ip == nil {
		t.Logf("%v has no ipv4", ifi)
		t.FailNow()
	}
	t.Log(ip.String())

	udpAddr, err := net.ResolveUDPAddr("udp", ip.String()+":7819")
	if err != nil {
		t.Fatal(err)
	}
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	_, err = conn.Write(xrand.Bytes(10))
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetOutboundIPString(t *testing.T) {
	t.Log(GetOutboundIPString())
}

func TestGetPrivateIPv4(t *testing.T) {
	GetPrivateIPv4(GetPrivateIPv4Interface())
}
