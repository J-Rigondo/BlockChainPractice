package utils

import (
	"bytes"
	"encoding/gob"
	"log"
)

func HandleError(err error) {
	if err != nil {
		log.Panic(err)
	}
}

/*
*
byte 배열로 변환해주는 함수
interface{} 모든 형태 다 받음
*/
func ToBytes(i interface{}) []byte {
	var Buffer bytes.Buffer

	encoder := gob.NewEncoder(&Buffer)
	err := encoder.Encode(i)
	HandleError(err)

	return Buffer.Bytes()
}

func FromBytes(i interface{}, data []byte) {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(i)

	HandleError(err)
}
