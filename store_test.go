package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	key := "mybestpicture"
	pathname := CASPathTransformFunc(key)
	fmt.Println(pathname)
}

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: DefaultPathTransformFunc,
	}
	s := NewStore(opts)

	data := bytes.NewReader([]byte("some jpg pics"))
	if err := s.writeStream("myspecialpicture", data); err != nil {
		t.Error(err)
	}
}
