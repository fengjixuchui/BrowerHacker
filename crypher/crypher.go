package crypher

import (
	"syscall"
	"unsafe"
)
// 8.0之前的加密
type DATA_BLOB struct {
	cbData uint32
	pbData *byte
}
// 初始化一下DATA_BLOB结构体
func InitBlob(data []byte) *DATA_BLOB{
	if len(data) == 0{
		return &DATA_BLOB{}
	}
	return &DATA_BLOB{
		pbData: &data[0],
		cbData: uint32(len(data)),
	}
}
// 转字节数组
func (b *DATA_BLOB) ToByteArray() []byte {
	d := make([]byte, b.cbData)
	copy(d, (*[1 << 30]byte)(unsafe.Pointer(b.pbData))[:])
	return d
}
// 调用win32 api进行解密
func WinDecypt(data []byte) ([]byte, error){
	dllcrypt32 := syscall.NewLazyDLL("Crypt32.dll")
	dllkernel32 := syscall.NewLazyDLL("Kernel32.dll")
	procDecryptData := dllcrypt32.NewProc("CryptUnprotectData")
	procLocalFree := dllkernel32.NewProc("LocalFree")

	var outblob DATA_BLOB
	r, _, err := procDecryptData.Call(uintptr(unsafe.Pointer(InitBlob(data))), 0, 0, 0, 0, 0, uintptr(unsafe.Pointer(&outblob)))
	if r == 0 {
		return nil, err
	}
	defer procLocalFree.Call(uintptr(unsafe.Pointer(outblob.pbData)))
	return outblob.ToByteArray(), nil
}


// 8.0 之后的加密
