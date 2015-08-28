package utee

import (
	"crypto/md5"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io"
	"math/rand"
	"net/http"
	"net/smtp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

var (
	//passwd md5
	PlainMd5 = Md5Str("")
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func Md5Str(salt string) func(string) string {
	return func(s string) string {
		h := md5.New()
		io.WriteString(h, s)
		io.WriteString(h, salt)
		return hex.EncodeToString(h.Sum(nil))
	}
}

func Chk(err error) {
	if err != nil {
		panic(err)
	}
}

//truncate string
func Truncate(s string, n int) string {
	length := utf8.RuneCountInString(s)
	if length <= n || n < 0 {
		return ""
	}

	l := []rune{}
	for _, r := range s {
		l = append(l, r)
	}

	l = l[:(length - n)]
	return string(l)
}

func Tick(t ...time.Time) int64 {
	if len(t) == 0 {
		return time.Now().UnixNano() / 1e6
	} else {
		return t[0].UnixNano() / 1e6
	}
}

func TickSec() int64 {
	return time.Now().Unix()
}

func TickHour() int64 {
	return time.Now().Unix() / 3600 * 3600
}

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func ContainsInf(s []interface{}, e interface{}) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func ContainsObjectId(s []bson.ObjectId, e string) bool {
	for _, a := range s {
		if a.Hex() == e {
			return true
		}
	}
	return false
}

func ToTimeMillis(fmt string, timeStr string) (int64, error) {
	data, err := time.Parse(fmt, timeStr)
	if err != nil {
		return 0, err
	}
	return data.UnixNano() / 1e6, nil
}

func Md5(b []byte) []byte {
	h := md5.New()
	h.Write(b)
	return h.Sum(nil)
}

// parse form into J
func F2j(r *http.Request) J {
	r.ParseForm()
	j := J{}
	for k, v := range r.Form {
		if len(v) == 1 {
			if len(v[0]) > 0 {
				j[k] = v[0]
			}
		} else {
			j[k] = v
		}
	}
	return j
}

func Fint64(s interface{}, dft ...int64) int64 {
	var i int64
	var err error
	if s == nil {
		if len(dft) > 0 {
			return dft[0]
		} else {
			s = "0"
		}
	}
	i, err = strconv.ParseInt(s.(string), 10, 64)
	if err != nil && len(dft) > 0 {
		i = dft[0]
	}
	return i
}

func Fint(s interface{}, dft ...int64) int {
	if len(dft) > 0 {
		return int(Fint64(s, dft[0]))
	} else {
		return int(Fint64(s))
	}
}

func DeleteMap(m map[string]interface{}, ks ...string) {
	for _, v := range ks {
		delete(m, v)
	}
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

func NotFound(err error) bool {
	return err == mgo.ErrNotFound
}

func Shuffle(src []string) []string {
	dest := make([]string, len(src))
	perm := rand.Perm(len(src))
	for i, v := range perm {
		dest[v] = src[i]
	}
	return dest
}

func SendMail(user, password, host, to, subject, body, mailtype string) error {
	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	var content_type string
	if mailtype == "html" {
		content_type = "Content-Type: text/" + mailtype + "; charset=UTF-8"
	} else {
		content_type = "Content-Type: text/plain" + "; charset=UTF-8"
	}

	msg := []byte("To: " + to + "\r\nFrom: " + user + "<" + user + ">\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)
	send_to := strings.Split(to, ";")
	err := smtp.SendMail(host, auth, user, send_to, msg)
	return err
}

func ParseUrl(s string) (string, int, error) {
	a := strings.Split(s, ":")
	if len(a) != 2 {
		return "", 0, fmt.Errorf("bad url %s", s)
	}
	port, err := strconv.Atoi(a[1])
	return a[0], port, err
}

func Unique(data []interface{}) []interface{} {
	m := map[interface{}]interface{}{}

	for _, d := range data {
		m[d] = "0"
	}

	l := []interface{}{}
	for key := range m {
		l = append(l, key)
	}
	return l
}

func UniqueStr(data []string) []string {
	m := map[string]string{}

	for _, d := range data {
		m[d] = "0"
	}

	l := []string{}
	for key := range m {
		l = append(l, key)
	}
	return l
}

//split a into several parts, no more than n
func SplitSlice(a []string, n int) [][]string {
	if len(a) < n || n == 1 {
		return [][]string{a}
	}

	result := make([][]string, n)
	for i, s := range a {
		idx := i % n
		result[idx] = append(result[idx], s)
	}
	return result
}
