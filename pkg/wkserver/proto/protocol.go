package proto

import wkproto "github.com/WuKongIM/WuKongIMGoProto"

type MsgType uint8 // 消息类型
const (
	Unknown          MsgType = iota
	MsgTypeConnect           // connect
	MsgTypeConnack           // connack
	MsgTypeRequest           // request
	MsgTypeResp              // response
	MsgTypeHeartbeat         // heartbeat
	MsgTypeMessage           // message
)

const (
	MsgTypeLength    = 1
	MsgContentLength = 4
)

func (m MsgType) Uint8() uint8 {
	return uint8(m)
}

type Protocol interface {
	Decode(data []byte) ([]byte, MsgType, int, error)
	Encode(data []byte, msgType uint8) ([]byte, error)
}

type DefaultProto struct {
}

func New() *DefaultProto {

	return &DefaultProto{}
}

func (d *DefaultProto) Decode(data []byte) ([]byte, MsgType, int, error) {
	if len(data) <= MsgContentLength {
		return nil, 0, 0, nil
	}
	decoder := wkproto.NewDecoder(data)
	msgType, err := decoder.Uint8()
	if err != nil {
		return nil, 0, 0, err
	}
	if msgType == MsgTypeHeartbeat.Uint8() {
		return nil, MsgTypeHeartbeat, MsgTypeLength, nil
	}
	contentLen, err := decoder.Uint32()
	if err != nil {
		return nil, 0, 0, err
	}
	if contentLen > uint32(len(data)-MsgTypeLength-MsgContentLength) {
		return nil, 0, 0, nil
	}
	contentBytes, err := decoder.Bytes(int(contentLen))
	if err != nil {
		return nil, 0, 0, err
	}
	return contentBytes, MsgType(msgType), len(contentBytes) + MsgTypeLength + MsgContentLength, nil
}

func (d *DefaultProto) Encode(data []byte, msgType uint8) ([]byte, error) {
	encoder := wkproto.NewEncoder()
	defer encoder.End()
	encoder.WriteUint8(msgType)
	encoder.WriteUint32(uint32(len(data)))
	encoder.WriteBytes(data)
	return encoder.Bytes(), nil
}
