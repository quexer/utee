package utee

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

var (
	// PlainMd5 string md5 function with empty salt
	PlainMd5 = Md5Str("")
	// PlainSha1 string sha-1 function with empty salt
	PlainSha1 = Sha1Str("")
)

// Md5Str create string md5 function with salt
func Md5Str(salt string) func(string) string {
	return func(s string) string {
		h := md5.New()
		io.WriteString(h, s)
		io.WriteString(h, salt)
		return hex.EncodeToString(h.Sum(nil))
	}
}

// Sha1Str create string sha-1 function with salt
func Sha1Str(salt string) func(string) string {
	return func(s string) string {
		h := sha1.New()
		io.WriteString(h, s)
		io.WriteString(h, salt)
		return hex.EncodeToString(h.Sum(nil))
	}
}

// Sha256Str create string sha-256 function with salt
func Sha256Str(salt string) func(string) string {
	return func(s string) string {
		h := sha256.New()
		io.WriteString(h, s)
		io.WriteString(h, salt)
		return hex.EncodeToString(h.Sum(nil))
	}
}

func HmacSha256(s string, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	io.WriteString(h, s)
	return hex.EncodeToString(h.Sum(nil))
}

func Chk(err error) {
	if err != nil {
		panic(err)
	}
}

func Md5(b []byte) []byte {
	h := md5.New()
	h.Write(b)
	return h.Sum(nil)
}

func IsPemExpire(b []byte) (bool, error) {
	block, _ := pem.Decode(b)
	if block == nil {
		return false, errors.New("failed to parse certificate PEM")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return false, err
	}
	return cert.NotAfter.Before(time.Now()), nil
}

func ParseAddr(s string) (string, int, error) {
	a := strings.Split(s, ":")
	if len(a) != 2 {
		return "", 0, fmt.Errorf("bad url %s", s)
	}
	port, err := strconv.Atoi(a[1])
	return a[0], port, err
}
