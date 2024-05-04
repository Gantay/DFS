package main

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	key := "mybestpicture"
	pathKey := CASPathTransformFunc(key)

	expectedOriginalKey := "be17b32c2870b1c0c73b59949db6a3be7814dd23"
	expectedPathName := "be17b/32c28/70b1c/0c73b/59949/db6a3/be781/4dd23"
	if pathKey.PathName != expectedPathName {
		t.Errorf("Have %s want %s", pathKey.PathName, expectedPathName)
	}
	if pathKey.Filename != expectedPathName {
		t.Errorf("Have %s want %s", pathKey.Filename, expectedOriginalKey)
	}
}

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	s := NewStore(opts)
	key := "momsspecials"

	data := []byte("some jpg bytes")
	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

	r, err := s.Read(key)
	if err != nil {
		t.Error(err)
	}

	b, _ := ioutil.ReadAll(r)

	fmt.Println(string(b))

	if string(b) != string(data) {
		t.Errorf("want this: %s got this: %s", data, b)
	}

}
