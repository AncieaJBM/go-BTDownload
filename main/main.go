package main

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"go-BTDownload/torrent"
	"os"
)

func main() {
	//server := "www.baidu.com:80"
	//_, err := net.DialTimeout("tcp", server, 5*time.Second)
	//if err != nil {
	//	fmt.Print("nimade")
	//}
	//parse torrent file
	file, err := os.Open("../testfile/debian-12.4.0-arm64-netinst.iso.torrent")
	if err != nil {
		fmt.Println("open file error")
		return
	}
	defer file.Close()
	tf, err := torrent.ParseFile(bufio.NewReader(file))
	if err != nil {
		fmt.Println("parse file error")
		return
	}
	// random peerId
	var peerId [torrent.IDLEN]byte
	_, _ = rand.Read(peerId[:])
	//connect tracker & find peers
	peers := torrent.FindPeers(tf, peerId)
	if len(peers) == 0 {
		fmt.Println("can not find peers")
		return
	}
	// build torrent task
	task := &torrent.TorrentTask{
		PeerId:   peerId,
		PeerList: peers,
		InfoSHA:  tf.InfoSHA,
		FileName: tf.FileName,
		FileLen:  tf.FileLen,
		PieceLen: tf.PieceLen,
		PieceSHA: tf.PieceSHA,
	}
	//download from peers & make file
	torrent.DownLoad(task)
}
