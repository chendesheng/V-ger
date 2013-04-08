package bencoding

import (
	"errors"
	// "io"
	// "fmt"
	"strconv"
	// "log"
	// "io/ioutil"
)

type list []interface{}
type dict map[string]interface{}

func Parse(bytes []byte) (res interface{}, err error) {
	defer (func() {
		if r := recover(); r != nil {
			res = nil
			err = r.(error)
		}
	})()
	_, res = parseObject(bytes, 0)
	err = nil
	return
}

func parseObject(bytes []byte, pos int) (int, interface{}) {
	b := bytes[pos]
	switch b {
	case 'l':
		return parseList(bytes, pos)
	case 'd':
		return parseDict(bytes, pos)
	case 'i':
		return parseInteger(bytes, pos)
	default:
		if b >= '0' && b <= '9' {
			return parseString(bytes, pos)
		} else {
			panic(errors.New("Unexcepted char"))
		}
	}

	return pos, nil
}

func parseList(bytes []byte, pos int) (int, []interface{}) {
	if bytes[pos] != 'l' {
		panic(errors.New("Unexcepted char"))
	}
	pos++
	res := make([]interface{}, 0)
	for {
		var obj interface{}
		pos, obj = parseObject(bytes, pos)
		res = append(res, obj)

		if bytes[pos] == 'e' {
			pos++
			break
		}
	}

	return pos, res
}

func parseDict(bytes []byte, pos int) (int, map[string]interface{}) {
	if bytes[pos] != 'd' {
		panic(errors.New("Unexcepted char"))
	}
	pos++

	res := make(map[string]interface{})
	for {
		var key string
		var val interface{}
		pos, key = parseString(bytes, pos)
		pos, val = parseObject(bytes, pos)
		res[key] = val

		if bytes[pos] == 'e' {
			pos++
			break
		}
	}
	return pos, res
}

func parseInteger(bytes []byte, pos int) (int, int64) {
	if bytes[pos] != 'i' {
		panic(errors.New("Unexcepted char"))
	}
	pos++

	s := pos
	for ; pos < len(bytes); pos++ {
		b := bytes[pos]
		if b >= '0' && b <= '9' {
			continue
		}
		if b == 'e' {
			break
		} else {
			panic(errors.New("Unexcepted char"))
		}
	}

	n, err := strconv.ParseInt(string(bytes[s:pos]), 10, 64)
	if err != nil {
		panic(err)
	}

	return pos + 1, n
}
func parseString(bytes []byte, pos int) (int, string) {
	s := pos
	for ; pos < len(bytes); pos++ {
		b := bytes[pos]
		if b >= '0' && b <= '9' {
			continue
		}
		if b == ':' {
			break
		} else {
			panic(errors.New("Unexcepted char"))
		}
	}

	n, err := strconv.Atoi(string(bytes[s:pos]))
	if err != nil {
		panic(err)
	}
	pos++
	if pos+n <= len(bytes) {
		return pos + n, string(bytes[pos : pos+n])
	} else {
		panic("Unexcepted EOF")
	}

	return pos, ""
}
