package torrent

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"net"
	"os"
	"testing"
)

func TestPeer(t *testing.T) {
	var peer PeerInfo
	peer.Ip = net.ParseIP("5.135.160.83")
	peer.Port = uint16(51413)

	file, _ := os.Open("../testfile/debian-12.4.0-arm64-netinst.iso.torrent")
	tf, _ := ParseFile(bufio.NewReader(file))

	var peerId [IDLEN]byte
	_, _ = rand.Read(peerId[:])

	conn, err := NewConn(peer, tf.InfoSHA, peerId)
	if err != nil {
		t.Error("new peer err : " + err.Error())
	}
	fmt.Println(conn)
}
