package aescbc

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"log"
	"sync"
	"vshop/lib/config"
)

type Crypter struct {
	key    []byte
	iv     []byte
	cipher cipher.Block
}

var theCrypter *Crypter
var once sync.Once

func initCrypter() {

	cfg, err := config.Configuration()
	if err != nil {
		log.Fatalf("Error reading configuration file: %s", err)
	}

	// initialize crypter singleton
	key := cfg.Security.WorkingKey
	keyBytes := md5.Sum([]byte(key))

	iv, err := hex.DecodeString(cfg.Security.IVStr)
	if err != nil {
		log.Fatalf("Error reading initializing vector: %s", err)
	}

	c, err := aes.NewCipher(keyBytes[:])
	if err != nil {
		log.Fatalf("Error creating new AES Cipher: %s", err)
	}

	theCrypter = &Crypter{keyBytes[:], iv, c}
}

func NewCrypter() *Crypter {

	once.Do(initCrypter)
	return theCrypter
}

func (e *Crypter) Encrypt(buf []byte) ([]byte, error) {

	padBuf := e.pad(buf)

	ciphertext := make([]byte, len(padBuf))
	mode := cipher.NewCBCEncrypter(e.cipher, e.iv)
	mode.CryptBlocks(ciphertext, padBuf)

	return ciphertext, nil
}

func (e *Crypter) Decrypt(buf []byte) ([]byte, error) {

	padBuf := make([]byte, len(buf))

	mode := cipher.NewCBCDecrypter(e.cipher, e.iv)
	mode.CryptBlocks(padBuf, buf)
	unpadBuf, err := e.unpad(padBuf)
	if err != nil {
		return nil, err
	}

	return unpadBuf, nil
}

func (e *Crypter) pad(buf []byte) []byte {

	blockSize := e.cipher.BlockSize()

	padByte := blockSize - len(buf)%blockSize
	padding := bytes.Repeat([]byte{uint8(padByte)}, padByte)

	return append(buf, padding...)
}

func (e *Crypter) unpad(buf []byte) ([]byte, error) {

	bufLen := len(buf)
	if bufLen == 0 {
		return nil, errors.New("invalid padding size")
	}

	padByte := buf[bufLen-1]
	if padByte == 0 {
		return nil, errors.New("invalid last byte of padding")
	}

	blockSize := e.cipher.BlockSize()
	padLen := int(padByte)
	if padLen > bufLen || padLen > blockSize {
		return nil, errors.New("invalid padding size")
	}

	for _, v := range buf[bufLen-padLen : bufLen-1] {
		if v != padByte {
			return nil, errors.New("invalid padding")
		}
	}

	return buf[:bufLen-padLen], nil
}
