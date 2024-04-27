package model

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

type Coord struct {
	Lat, Long float64
}

func (c *Coord) Scan(src interface{}) error {
	switch src.(type) {
	case []byte:
		var b = src.([]byte)
		if len(b) != 25 {
			return errors.New(fmt.Sprintf("Expected []bytes with length 25, got %d", len(b)))
		}
		var longitude float64
		var latitude float64
		buf := bytes.NewReader(b[9:17])
		err := binary.Read(buf, binary.LittleEndian, &longitude)
		if err != nil {
			return err
		}
		buf = bytes.NewReader(b[17:25])
		err = binary.Read(buf, binary.LittleEndian, &latitude)
		if err != nil {
			return err
		}
		*c = Coord{longitude, latitude}
	default:
		return errors.New(fmt.Sprintf("Expected []byte for Location type, got  %T", src))
	}
	return nil
}
