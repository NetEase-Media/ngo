// Copyright Ngo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

// Encrypt 加密
func Encrypt(content, key string) (string, error) {
	if len(content) == 0 || len(key) == 0 {
		return "", errors.New("content or key must be not empty")
	}
	kb := []byte(key)
	cb := []byte(content)
	//创建加密算法实例
	block, err := aes.NewCipher(kb)
	if err != nil {
		return "", err
	}
	//获取块的大小
	blockSize := block.BlockSize()
	//对数据进行填充，让数据长度满足需求
	cb = PKCS7Padding(cb, blockSize)
	//采用AES加密方法中CBC加密模式
	crypted := make([]byte, len(cb))
	for bs, be := 0, blockSize; bs < len(cb); bs, be = bs+blockSize, be+blockSize {
		block.Encrypt(crypted[bs:be], cb[bs:be])
	}
	return base64.StdEncoding.EncodeToString(crypted), nil
}

// Decrypt 解密
func Decrypt(content, key string) (string, error) {
	if len(content) == 0 || len(key) == 0 {
		return "", errors.New("content or key must be not empty")
	}
	//解密base64字符串
	cb, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return "", err
	}
	kb := []byte(key)
	//创建加密算法实例
	block, err := aes.NewCipher(kb)
	if err != nil {
		return "", err
	}
	origData := make([]byte, len(cb))
	blockSize := block.BlockSize()
	for bs, be := 0, blockSize; bs < len(cb); bs, be = bs+blockSize, be+blockSize {
		block.Decrypt(origData[bs:be], cb[bs:be])
	}
	//去除填充字符串
	origData, err = PKCS7UnPadding(origData)
	if err != nil {
		return "", err
	}
	return string(origData), err
}

// PKCS7Padding PKCS7 填充模式
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	//Repeat()函数的功能是把切片[]byte{byte(padding)}复制padding个，然后合并成新的字节切片返回
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// PKCS7UnPadding 填充的反向操作，删除填充字符串
func PKCS7UnPadding(origData []byte) ([]byte, error) {
	//获取数据长度
	length := len(origData)
	if length == 0 {
		return nil, errors.New("加密字符串错误！")
	}
	//获取填充字符串长度
	unpadding := int(origData[length-1])

	if unpadding > length {
		return nil, errors.New("加密字符串错误！")
	}

	//截取切片，删除填充字节，并且返回明文
	return origData[:(length - unpadding)], nil

}

// AES/ECB/PKCS7模式加密--签名加密方式
func ECBEncrypt(data []byte, key string) (string, error) {
	if data == nil || len(data) == 0 || len(key) == 0 {
		return "", errors.New("data or key must be not empty")
	}
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	ecb := NewECBEncryptEr(block)
	// 加PKCS7填充
	content := PKCS7Padding(data, block.BlockSize())
	encryptData := make([]byte, len(content))
	// 生成加密数据
	ecb.CryptBlocks(encryptData, content)
	return base64.StdEncoding.EncodeToString(encryptData), nil
}

// AES/ECB/PKCS7模式解密--签名解密方式
func ECBDecrypt(data []byte, key string) (string, error) {
	if data == nil || len(data) == 0 || len(key) == 0 {
		return "", errors.New("data or key must be not empty")
	}
	decryptData, _ := base64.StdEncoding.DecodeString(string(data))
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	ecb := NewECBDecryptEr(block)
	retData := make([]byte, len(decryptData))
	ecb.CryptBlocks(retData, decryptData)
	// 解PKCS7填充
	retData, err = PKCS7UnPadding(retData)
	if err != nil {
		return "", err
	}
	return string(retData), nil
}

type ecb struct {
	b         cipher.Block
	blockSize int
}

func newECB(b cipher.Block) *ecb {
	return &ecb{
		b:         b,
		blockSize: b.BlockSize(),
	}
}

type ecbEncryptEr ecb

func NewECBEncryptEr(b cipher.Block) cipher.BlockMode {
	return (*ecbEncryptEr)(newECB(b))
}

func (x *ecbEncryptEr) BlockSize() int { return x.blockSize }

func (x *ecbEncryptEr) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Encrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

// ecb解密方法
type ecbDecryptEr ecb

func NewECBDecryptEr(b cipher.Block) cipher.BlockMode {
	return (*ecbDecryptEr)(newECB(b))
}

func (x *ecbDecryptEr) BlockSize() int { return x.blockSize }

func (x *ecbDecryptEr) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Decrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}
