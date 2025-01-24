package websocket

import "time"

type FrameType uint8

const (
	FrameData      FrameType = 0x0
	FramePing      FrameType = 0x1
	FrameAck       FrameType = 0x2
	FrameNoAck     FrameType = 0x3
	FrameTranspond FrameType = 0x6
	FrameErr       FrameType = 0x9
)

type Message struct {
	Id        string `json:"id"`
	FrameType `json:"frameType"`
	AckSeq    int         `json:"ackSeq"`
	ackTime   time.Time   `json:"ackTime"`
	errCount  int         `json:"errCount"`
	Method    string      `json:"method"`
	FromId    string      `json:"fromId"`
	Data      interface{} `json:"data"`
	// 转发目标
	TranspondUid string `json:"transpondUid"`
}

func NewMsg(fromId string, data interface{}) Message {
	return Message{
		FrameType: FrameData,
		FromId:    fromId,
		Data:      data,
	}
}

func NewErrMsg(err error) *Message {
	return &Message{
		FrameType: FrameErr,
		Data:      err.Error(),
	}
}
