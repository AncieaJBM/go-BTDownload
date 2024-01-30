#### BT下载器流程

![782a19928f1222a33d8b7cdd5f3fab41](https://github.com/AncieaJBM/go-BTDownload/assets/157377021/003df1ab-5791-4333-b584-d4460c02309e)


1、MsgPiece消息用于发送一个块（piece）的数据给对端。消息格式如下：

```
+---------------+
| Length Prefix | 4 bytes, big-endian integer
+---------------+
| Message ID    | 1 byte
+---------------+
| Piece Index   | 4 bytes, big-endian integer
+---------------+
| Block Offset  | 4 bytes, big-endian integer
+---------------+
| Block Data    | variable-length string
+---------------+
```

其中，Length Prefix表示整个消息的长度（不包括自己长度前缀这个字段），Message ID表示消息类型，对于MsgPiece消息来说就是0x07，Piece Index表示块所属的分片（piece）的索引，Block Offset表示块在分片中的偏移量（即分片内的块编号），Block Data表示块的实际数据内容。

2、MsgHave消息表示发送端已经拥有了某个块（piece）的数据。消息的大小固定为9个字节，前四个字节为0x00000005，表示消息长度为5个字节，紧随其后的一个字节为0x04表示该消息的类型是MsgHave，后面跟着4个字节的整数表示所拥有的块的索引。

  因此可以看作是一个固定长度为9字节的二进制数据。以下是MsgHave消息的完整结构：

```
  +---------------+
  | Length Prefix | 4 bytes, big-endian integer
  +---------------+
  | Message ID    | 1 byte
  +---------------+
  | Piece Index   | 4 bytes, big-endian integer
  +---------------+
```

