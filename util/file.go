package util

import "os"
//这里应该是UInt32,函数命名需调整，后续看心情吧
func ReadUInt64(file *os.File, index int64) uint64 {
	_, _ = file.Seek(index, 0)
	bytes := make([]byte, 4)
	_, err := file.Read(bytes)
	if err != nil {
		return 0
	}
	return (uint64(bytes[0]) & 0xff) | (uint64(bytes[1]&0xff) << 8) | (uint64(bytes[2]&0xff) << 16) | (uint64(bytes[3]&0xff) << 24)
}

func WriteUInt64(file *os.File, index int64, val uint64) {
	_, _ = file.Seek(index, 0)
	_, _ = file.Write([]byte{
		byte(val), byte(val >> 8), byte(val >> 16), byte(val >> 24),
	})

}
