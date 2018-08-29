package multihash

import (
	"bytes"
	"io"
	"testing"
)

type evilReader struct {
	buffer []byte
}

func (er *evilReader) Read(buf []byte) (int, error) {
	n := copy(buf, er.buffer)
	er.buffer = er.buffer[n:]
	var err error
	if len(er.buffer) == 0 {
		err = io.EOF
	}
	return n, err
}

func TestEvilReader(t *testing.T) {
	emptyHash, err := Sum(nil, ID, 0)
	if err != nil {
		t.Fatal(err)
	}
	r := NewReader(&evilReader{emptyHash.Bytes()})
	h, err := r.ReadMultihash()
	if err != nil {
		t.Fatal(err)
	}
	if h != emptyHash {
		t.Fatal(err)
	}
	h, err = r.ReadMultihash()
	if h != Nil || err != io.EOF {
		t.Fatal("expected end of file")
	}
}

func TestReader(t *testing.T) {

	var buf bytes.Buffer

	for _, tc := range testCases {
		m, err := tc.Multihash()
		if err != nil {
			t.Fatal(err)
		}

		buf.Write(m.Bytes())
	}

	r := NewReader(&buf)

	for _, tc := range testCases {
		h, err := tc.Multihash()
		if err != nil {
			t.Fatal(err)
		}

		h2, err := r.ReadMultihash()
		if err != nil {
			t.Error(err)
			continue
		}

		if h != h2 {
			t.Error("h and h2 should be equal")
		}
	}
}

func TestWriter(t *testing.T) {

	var buf bytes.Buffer
	w := NewWriter(&buf)

	for _, tc := range testCases {
		m, err := tc.Multihash()
		if err != nil {
			t.Error(err)
			continue
		}

		if err := w.WriteMultihash(m); err != nil {
			t.Error(err)
			continue
		}

		buf2 := make([]byte, len(m.Bytes()))
		if _, err := io.ReadFull(&buf, buf2); err != nil {
			t.Error(err)
			continue
		}

		if m.Binary() != string(buf2) {
			t.Error("m and buf2 should be equal")
		}
	}
}
