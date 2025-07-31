package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

type Value struct {
	typ   string
	str   string
	num   int
	bulk  string
	array []Value
}

type Resp struct {
	reader *bufio.Reader
}

func NewResp(rd io.Reader) *Resp {
	return &Resp{bufio.NewReader(rd)}
}

func (r *Resp) readLine() (line []byte, n int, err error) {
	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			return nil, 0, err
		}
		n += 1
		line = append(line, b)
		if len(line) >= 2 && line[len(line)-2] == '\r' {
			break
		}
	}
	return line[:len(line)-2], n, nil
}

func (r *Resp) readInteger() (x int, n int, err error) {
	line, n, err := r.readLine()
	if err != nil {
		return 0, 0, err
	}
	i64, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, n, err
	}
	return int(i64), n, nil
}

func (r *Resp) readArray() (val Value, err error) {
	val = Value{typ: "array"}

	arrayLength, _, err := r.readInteger()
	if err != nil {
		return val, err
	}

	val.array = make([]Value, 0, arrayLength)
	for i := 0; i < arrayLength; i++ {
		v, err := r.Read()
		if err != nil {
			return val, err
		}
		val.array = append(val.array, v)
	}

	return val, nil
}

func (r *Resp) readBulk() (val Value, err error) {
	val = Value{typ: "bulk"}

	bulkLength, _, err := r.readInteger()
	if err != nil {
		return val, err
	}

	bulk := make([]byte, bulkLength)
	_, _ = r.reader.Read(bulk)
	val.bulk = string(bulk)

	_, _, _ = r.readLine() // Read the trailing CRLF
	return val, nil
}

func (r *Resp) Read() (Value, error) {
	_type, err := r.reader.ReadByte()

	if err != nil {
		return Value{}, err
	}

	switch _type {
	case ARRAY:
		return r.readArray()
	case BULK:
		return r.readBulk()
	default:
		fmt.Printf("Unknown type: %v", string(_type))
		return Value{}, nil
	}
}

type Type byte

func ParseInput(data []byte) {
	rd := Type(data[0])
	switch rd {
	case INTEGER:
	case STRING:
	case BULK:
	case ARRAY:
	case ERROR:
	default:
		fmt.Println("Unknown response type: ", rd)
	}
}
