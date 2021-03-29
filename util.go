package httpRequest

import (
	"bytes"
	"encoding/binary"
)

func IntByte(v interface{}) []byte {
	b := bytes.NewBuffer([]byte{})
	switch v.(type) {
	case int:
		binary.Write(b, binary.BigEndian, int64(v.(int)))
	case int8:
		binary.Write(b, binary.BigEndian, v.(int8))
	case int16:
		binary.Write(b, binary.BigEndian, v.(int16))
	case int32:
		binary.Write(b, binary.BigEndian, v.(int32))
	case int64:
		binary.Write(b, binary.BigEndian, v.(int64))
	case uint:
		binary.Write(b, binary.BigEndian, uint64(v.(uint)))
	case uint8:
		binary.Write(b, binary.BigEndian, v.(uint8))
	case uint16:
		binary.Write(b, binary.BigEndian, v.(uint16))
	case uint32:
		binary.Write(b, binary.BigEndian, v.(uint32))
	case uint64:
		binary.Write(b, binary.BigEndian, v.(uint64))
	}
	return b.Bytes()
}
