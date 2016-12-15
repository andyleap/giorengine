package giorengine

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

const (
	resetPosition = 32768
)

type Uint16Reader interface {
	ReadUint16() (uint16, error)
}

type uint16Reader struct {
	data  io.Reader
	order binary.ByteOrder
}

func (r *uint16Reader) ReadUint16() (uint16, error) {
	var word uint16
	e := binary.Read(r.data, r.order, &word)
	return word, e
}

func NewReaderUint16BE(r io.Reader) (io.ReadCloser, error) {
	return NewReader(&uint16Reader{r, binary.BigEndian})
}

func NewReader(r Uint16Reader) (io.ReadCloser, error) {
	d := &Reader{
		position:  resetPosition,
		dict:      map[int]string{0: "0", 1: "1", 2: "2"},
		dictSize:  4,
		result:    &bytes.Buffer{},
		numBits:   3,
		enlargeIn: 4,
		data:      r,
	}

	if e := d.process(); e != nil {
		return nil, e
	}

	// Get rid of original
	d.data = nil

	return d, nil
}

type Reader struct {
	currentVal                                         uint16
	position                                           uint16
	data                                               Uint16Reader
	dict                                               map[int]string
	lastEntry                                          string
	result                                             *bytes.Buffer
	numBits, errorCount, dictSize, enlargeIn, readSize int
}

func (r *Reader) readBit() uint16 {
	res := r.currentVal & r.position
	r.position = r.position >> 1
	if r.position == 0 {
		r.position = resetPosition

		if val, e := r.data.ReadUint16(); e == io.EOF || e == nil {
			r.currentVal = val
		} else {
			fmt.Println(e)
		}
	}

	if res > 0 {
		return 1
	}
	return 0
}

func (r *Reader) readBits(numBits int) (res uint16) {
	maxpower := math.Pow(2, float64(numBits))
	power := 1
	for float64(power) != maxpower {
		res = res | (uint16(power) * r.readBit())
		power = power << 1
	}
	return
}

func (r *Reader) Close() error {
	if r.result == nil {
		return fmt.Errorf("Reader closed.")
	}
	r.result = nil
	return nil
}

func (r *Reader) process() error {

	if val, e := r.data.ReadUint16(); e == io.EOF || e == nil {
		r.currentVal = val
	} else {
		return e
	}

	var c uint16
	switch r.readBits(2) {
	case 0:
		c = r.readBits(8)
	case 1:
		c = r.readBits(16)
	case 2:
		return nil
	}

	r.dict[3] = fmt.Sprintf("%c", c)
	r.result.WriteRune(rune(c))
	r.lastEntry = r.dict[3]

	for {

		bitsRead := r.readBits(r.numBits)
		var dictIndex = r.dictSize

		switch bitsRead {
		case 0:
			if r.errorCount++; r.errorCount > 10000 {
				return fmt.Errorf("Error count exceded 10000")
			}
			c = r.readBits(8)
			r.dict[r.dictSize] = fmt.Sprintf("%c", c)
			r.dictSize++
			r.enlargeIn--

		case 1:
			c = r.readBits(16)
			r.dict[r.dictSize] = fmt.Sprintf("%c", c)
			r.dictSize++
			r.enlargeIn--

		case 2:
			return nil

		default:
			dictIndex = int(bitsRead)
		}

		if r.enlargeIn == 0 {
			r.enlargeIn = int(math.Pow(2, float64(r.numBits)))
			r.numBits++
		}

		v, exists := r.dict[dictIndex]
		if !exists {
			if dictIndex == r.dictSize {
				v = fmt.Sprintf("%s%c", r.lastEntry, []rune(r.lastEntry)[0])
			} else {
				return nil
			}
		}
		r.result.WriteString(v)
		r.dict[r.dictSize] = fmt.Sprintf("%s%c", r.lastEntry, []rune(v)[0])
		r.dictSize++
		r.lastEntry = v

		if r.enlargeIn--; r.enlargeIn == 0 {
			r.enlargeIn = int(math.Pow(2, float64(r.numBits)))
			r.numBits++
		}
	}

	return nil
}

func (r *Reader) Read(b []byte) (int, error) {
	if r.result == nil {
		return 0, fmt.Errorf("Reader closed.")
	}
	return r.result.Read(b)
}

func (r *Reader) String() string {
	return r.result.String()
}
