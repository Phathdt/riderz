package common

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"io"
)

const (
	wkbXDR byte = 0
	wkbNDR byte = 1
)

type Geometry interface {
	sql.Scanner
	driver.Valuer
	GetType() uint32
	Write(*bytes.Buffer) error
}

// read
func decode(value interface{}) (io.Reader, error) {
	ewkb, err := hex.DecodeString(value.(string))
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(ewkb), nil
}

func readEWKB(reader io.Reader, g Geometry) error {
	var byteOrder binary.ByteOrder
	var wkbByteOrder byte
	var wkbType uint32

	// Read as Little Endian to attempt to determine byte order
	if err := binary.Read(reader, binary.LittleEndian, &wkbByteOrder); err != nil {
		return err
	}

	// Decide byte order
	switch wkbByteOrder {
	case wkbXDR:
		byteOrder = binary.BigEndian
	case wkbNDR:
		byteOrder = binary.LittleEndian
	default:
		return errors.New("unsupported byte order")
	}

	// Determine the geometery type
	if err := binary.Read(reader, byteOrder, &wkbType); err != nil {
		return err
	}

	// Decode into our struct
	return binary.Read(reader, byteOrder, g)
}

// write
func writeEWKB(g Geometry) (*bytes.Buffer, error) {
	buffer := bytes.NewBuffer(nil)

	// Set our endianness
	if err := binary.Write(buffer, binary.LittleEndian, wkbNDR); err != nil {
		return nil, err
	}

	if err := binary.Write(buffer, binary.LittleEndian, g.GetType()); err != nil {
		return nil, err
	}

	if err := g.Write(buffer); err != nil {
		return nil, err
	}

	return buffer, nil
}

type PointS struct {
	SRID int32
	X, Y float64
}

func (p *PointS) Scan(value interface{}) error {
	reader, err := decode(value)
	if err != nil {
		return err
	}

	if err = readEWKB(reader, p); err != nil {
		return err
	}

	return nil
}

func (p PointS) Value() (driver.Value, error) {
	buffer, err := writeEWKB(&p)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (p PointS) Write(buffer *bytes.Buffer) error {
	err := binary.Write(buffer, binary.LittleEndian, &p)
	return err
}

func (p PointS) GetType() uint32 {
	return 0x20000001
}
