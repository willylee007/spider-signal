package main

// SignalMsg 信令消息
type SignalMsg struct {
	Type    string `json:"type"`
	Content `json:"content"`
}

// Content 信令内容
type Content struct {
	RoomName string `json:"roomName"`
	Msg      string `json:"msg"`
}
