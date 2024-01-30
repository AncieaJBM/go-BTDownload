package torrent

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"os"
)

type TorrentTask struct {
	PeerId   [20]byte
	PeerList []PeerInfo
	InfoSHA  [SHALEN]byte
	FileName string
	FileLen  int
	PieceLen int
	PieceSHA [][SHALEN]byte
}

type pieceTask struct {
	index  int
	sha    [SHALEN]byte
	length int
}

type pieceResult struct {
	index int
	data  []byte
}

type taskState struct {
	index      int
	conn       *PeerConn
	requested  int
	downloaded int
	backlog    int
	data       []byte
}

const BLOCKSIZE = 16384
const MAXBACKLOG = 5

func DownLoad(task *TorrentTask) error {

	fmt.Println("start downloading " + task.FileName)
	taskQueue := make(chan *pieceTask, len(task.PieceSHA))
	resultQueue := make(chan *pieceResult)
	//find每个piece的start end
	for index, sha := range task.PieceSHA {
		begin, end := task.getPieceBounds(index)
		taskQueue <- &pieceTask{index, sha, (end - begin)}
	}

	for _, peer := range task.PeerList {
		go task.peerRoutine(peer, taskQueue, resultQueue)
	}

	buf := make([]byte, task.FileLen)
	count := 0
	for count < len(task.PieceSHA) {
		res := <-resultQueue
		begin, end := task.getPieceBounds(res.index)
		copy(buf[begin:end], res.data)
		count++
		percent := float64(count) / float64(len(task.PieceSHA)) * 100
		fmt.Printf("downloading, progress : (%0.2f%%)\n", percent)
	}

	close(taskQueue)
	close(resultQueue)
	file, err := os.Create(task.FileName)
	if err != nil {
		fmt.Println("fail to create file: " + task.FileName)
		return err
	}
	_, err = file.Write(buf)
	if err != nil {
		fmt.Println("fail to write data")
		return err
	}
	return nil
}

func (t *TorrentTask) getPieceBounds(index int) (bengin, end int) {
	bengin = index * t.PieceLen
	end = bengin + t.PieceLen
	if end > t.FileLen {
		end = t.FileLen
	}
	return
}

func (t *TorrentTask) peerRoutine(peer PeerInfo, taskQueue chan *pieceTask, resultQueue chan *pieceResult) {
	conn, err := NewConn(peer, t.InfoSHA, t.PeerId)
	if err != nil {
		fmt.Println("fail to connect peer : " + peer.Ip.String())
		return
	}
	defer conn.Close()

	fmt.Println("complete handshake with peer : " + peer.Ip.String())
	// send Interested msg
	conn.WriteMsg(&PeerMsg{MsgInterested, nil})

	for task := range taskQueue {
		if !conn.Field.HasPiece(task.index) {
			// peer hasnt this piece
			taskQueue <- task
			continue
		}
		res, err := downloadPiece(conn, task)
		if err != nil {
			taskQueue <- task
			fmt.Println("fail to download piece" + err.Error())
			return
		}
		if !checkPiece(task, res) {
			taskQueue <- task
			continue
		}
		resultQueue <- res
	}

}

func checkPiece(task *pieceTask, res *pieceResult) bool {
	bt := sha1.Sum(res.data)
	if !bytes.Equal(bt[:], task.sha[:]) {
		fmt.Printf("check integrity failed, index :%v\n", res.index)
		return false
	}
	return true
}

func downloadPiece(conn *PeerConn, task *pieceTask) (*pieceResult, error) {
	state := &taskState{
		index: task.index,
		conn:  conn,
		data:  make([]byte, task.length),
	}

	for state.downloaded < task.length {
		if !conn.Choked {
			for state.requested < task.length && state.backlog < MAXBACKLOG {
				length := BLOCKSIZE
				if task.length-state.requested < length {
					length = task.length - state.requested
				}
				msg := NewRequestMsg(state.index, state.requested, length)
				_, err := state.conn.WriteMsg(msg)
				if err != nil {
					return nil, err
				}
				state.backlog++
				state.requested += length
			}
			err := state.handleMsg()
			if err != nil {
				return nil, err
			}
		}
	}
	return &pieceResult{state.index, state.data}, nil
}

func (state *taskState) handleMsg() error {
	msg, err := state.conn.ReadMsg()
	if err != nil {
		return err
	}
	// handle keep-alive
	if msg == nil {
		return nil
	}
	switch msg.Id {
	case MsgChoke:
		state.conn.Choked = true
	case MsgUnchoke:
		state.conn.Choked = false
	case MsgHave:
		index, err := GetHaveIndex(msg)
		if err != nil {
			return err
		}
		state.conn.Field.SetPiece(index)
	case MsgPiece:
		n, err := CopyPieceData(state.index, state.data, msg)
		if err != nil {
			return err
		}
		state.downloaded += n
		state.backlog--
	}
	return nil
}
