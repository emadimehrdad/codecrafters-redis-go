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

func (r *Resp) readLine() ([]byte, int, error) {
	line, err := r.reader.ReadString('\n')
	if err != nil {
		return nil, 0, err
	}
	if len(line) < 2 || line[len(line)-2] != '\r' {
		return nil, 0, fmt.Errorf("invalid line ending")
	}
	return []byte(line[:len(line)-2]), len(line), nil
}

func (r *Resp) readInteger() (x int, n int, err error) {
	line, n, err := r.readLine()
	if err != nil {
		fmt.Println(err)
		return 0, 0, err
	}
	i64, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		fmt.Println(err)
		return 0, n, err
	}
	return int(i64), n, nil
}

func (r *Resp) readArray() (val Value, err error) {
	val = Value{typ: "array"}

	arrayLength, _, err := r.readInteger()
	if err != nil {
		fmt.Println(err)
		return val, err
	}

	val.array = make([]Value, 0, arrayLength)
	for i := 0; i < arrayLength; i++ {
		v, err := r.Read()
		if err != nil {
			fmt.Println(err)
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
		fmt.Println(err)
		return val, err
	}

	if bulkLength == -1 {
		val.bulk = ""
		return val, nil
	}

	bulk := make([]byte, bulkLength)
	_, err = io.ReadFull(r.reader, bulk)
	if err != nil {
		return val, err
	}
	val.bulk = string(bulk)

	_, _, _ = r.readLine() // Read the trailing CRLF

	return val, nil
}

func (r *Resp) Read() (Value, error) {
	_type, err := r.reader.ReadByte()

	if err != nil {
		fmt.Println(err)
		return Value{}, err
	}

	switch _type {
	case ARRAY:
		return r.readArray()
	case BULK:
		return r.readBulk()
	case STRING:
		line, _, err := r.readLine()
		return Value{typ: "string", str: string(line)}, err
	case ERROR:
		line, _, err := r.readLine()
		return Value{typ: "error", str: string(line)}, err
	case INTEGER:
		num, _, err := r.readInteger()
		return Value{typ: "integer", num: num}, err
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
