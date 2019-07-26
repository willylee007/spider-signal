package main

// SignalMsg 信令消息
type SignalMsg struct {
	Type    string `json:"type"`
	Content `json:"content"`
}

// RelayClientSignal 服务器转发的消息
type RelayClientSignal struct {
	client *Client
	msg    []byte
}

// Content 信令内容
type Content struct {
	RoomName string `json:"roomName"`
	Msg      string `json:"msg"`
	Sdp      string `json:"sdp"`
}
