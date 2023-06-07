package utils

import (
	"bytes"
	"encoding/binary"
	"github.com/google/uuid"
	"github.com/username/schoolapp/models"
)

func Int64ArrayToBytes(arr []int64) []byte {
	buf := new(bytes.Buffer)
	for _, value := range arr {
		err := binary.Write(buf, binary.LittleEndian, value)
		if err != nil {
			panic(err) // handle the error according to your needs
		}
	}
	return buf.Bytes()
}

func BytesToInt64Array(data []byte) []int64 {
	buf := bytes.NewBuffer(data)
	arrayLength := len(data) / 8 // Each int64 occupies 8 bytes
	int64Array := make([]int64, arrayLength)
	for i := 0; i < arrayLength; i++ {
		err := binary.Read(buf, binary.LittleEndian, &int64Array[i])
		if err != nil {
			panic(err) // handle the error according to your needs
		}
	}
	return int64Array
}

func Int64ArrayToDbIdArray(arr []int64) []models.DbId {
	dbIdArray := make([]models.DbId, len(arr))
	for i, value := range arr {
		dbIdArray[i] = models.DbId(value)
	}
	return dbIdArray
}

func UuidArrayToBytes(arr []uuid.UUID) []byte {
	byteArray := make([]byte, 0)
	for _, u := range arr {
		byteArray = append(byteArray, u[:]...)
	}
	return byteArray
}

func BytesToUUIDArray(data []byte) []uuid.UUID {
	arrayLength := len(data) / 16 // Each UUID occupies 16 bytes
	uuidArray := make([]uuid.UUID, arrayLength)
	for i := 0; i < arrayLength; i++ {
		copy(uuidArray[i][:], data[i*16:(i+1)*16])
	}
	return uuidArray
}

func DbIdArrayToInt64Array(arr []models.DbId) []int64 {
	int64Array := make([]int64, len(arr))
	for i, value := range arr {
		int64Array[i] = int64(value)
	}
	return int64Array
}
