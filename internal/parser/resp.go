// Package parser implements the RESP protocol parser.
package parser

import (
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"unicode/utf8"
)

var CRLF = "\r\n"

var ErrIncomplete = errors.New("incomplete RESP value")

type Parser struct {
	buf    []byte
	logger *slog.Logger
}

func NewParser(logger *slog.Logger) *Parser {
	return &Parser{logger: logger}
}

func (p *Parser) Append(data []byte) {
	p.buf = append(p.buf, data...)
}

func (p *Parser) ParseNext() (Value, error) {
	value, consumed, err := p.parseOne(p.buf)
	if err != nil {
		return Value{}, err
	}
	p.buf = p.buf[consumed:]

	p.logger.Debug("recieved command", "value", value.Print(), "client buffer length", len(p.buf))
	return value, nil
}

func (p *Parser) parseOne(data []byte) (value Value, consumed int, err error) {
	if len(data) == 0 {
		return Value{}, 0, ErrIncomplete
	}

	dataPrefix, size := utf8.DecodeRune(data)
	dataAfterPrefix := data[size:]

	switch ValueType(dataPrefix) {
	case ArrayType:
		return p.parseArray(dataAfterPrefix)
	case BulkType:
		return p.parseBulk(dataAfterPrefix)
	default:
		return Value{}, 0, fmt.Errorf("unknown data type: %c", dataPrefix)
	}
}

func (p *Parser) parseArray(data []byte) (Value, int, error) {
	length, afterHeader, err := parseLength(data)
	if err != nil {
		return Value{}, 0, err
	}

	if length < 0 {
		return Value{}, 0, fmt.Errorf("invalid array length: %d", length)
	}

	value := Value{
		Type:  ArrayType,
		Array: make([]Value, 0, length),
	}

	consumed := len(data) - len(afterHeader) + 1

	for range length {
		element, n, err := p.parseOne(afterHeader)
		if err != nil {
			return Value{}, 0, fmt.Errorf("parse array element: %w", err)
		}
		value.Array = append(value.Array, element)
		afterHeader = afterHeader[n:]
		consumed += n
	}

	return value, consumed, nil
}

func (p *Parser) parseBulk(data []byte) (Value, int, error) {
	length, payload, err := parseLength(data)
	if err != nil {
		return Value{}, 0, err
	}

	if length < 0 {
		return Value{}, 0, fmt.Errorf("invalid bulk string length: %d", length)
	}

	if len(payload) < length+len(CRLF) {
		return Value{}, 0, ErrIncomplete
	}

	value := Value{
		Type: BulkType,
		Bulk: string(payload[:length]),
	}

	// prefix + header + length + CRLF
	consumed := len(data) - len(payload) + length + len(CRLF) + 1

	return value, consumed, nil
}

// parseLength reads a decimal length terminated by CRLF.
func parseLength(data []byte) (int, []byte, error) {
	before, after, ok := bytes.Cut(data, []byte(CRLF))
	if !ok {
		return 0, nil, ErrIncomplete
	}

	length, err := strconv.Atoi(string(before))
	if err != nil {
		return 0, nil, fmt.Errorf("invalid RESP length %q: %w", before, err)
	}

	return length, after, nil
}

func RESPValue(v Value) []byte {
	switch v.Type {
	case SimpleStingType:
		return []byte(string(v.Type) + v.SimpleString + CRLF)
	case ErrorType:
		return []byte(string(v.Type) + v.Error + CRLF)
	case NullType:
		return []byte(string(v.Type) + CRLF)
	case BulkType:
		return []byte(
			string(v.Type) +
				strconv.Itoa(len(v.Bulk)) +
				CRLF +
				v.Bulk +
				CRLF,
		)
	case ArrayType:
		var buf bytes.Buffer
		buf.WriteString(string(v.Type))
		buf.WriteString(strconv.Itoa(len(v.Array)))
		buf.WriteString(CRLF)

		for _, element := range v.Array {
			encoded := RESPValue(element)
			buf.Write(encoded)
		}
		return buf.Bytes()
	default:
		return []byte(
			string(v.Type) +
				strconv.Itoa(len(v.Array)) +
				CRLF,
		)
	}
}
