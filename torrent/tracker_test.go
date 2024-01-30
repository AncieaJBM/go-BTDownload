package torrent

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"os"
	"testing"
)

func TestTracker(t *testing.T) {
	file, _ := os.Open("../testfile/debian-12.4.0-amd64-netinst.iso.torrent")
	tf, _ := ParseFile(bufio.NewReader(file))

	var peerId [IDLEN]byte
	_, _ = rand.Read(peerId[:])

	peers := FindPeers(tf, peerId)
	for i, p := range peers {
		fmt.Printf("Peer %d, Ip: %s, Port: %d\n", i, p.Ip, p.Port)
	}
}
